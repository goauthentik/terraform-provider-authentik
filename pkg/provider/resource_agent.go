package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceAgent() *schema.Resource {
	return &schema.Resource{
		Description:   "Enterprise --- ",
		CreateContext: resourceAgentCreate,
		ReadContext:   resourceAgentRead,
		UpdateContext: resourceAgentUpdate,
		DeleteContext: resourceAgentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"parent": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "PK of the user that owns this agent. If not set, a new parent user is created automatically.",
			},
			"policy_behavior": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				Description:      helpers.EnumToDescription(api.AllowedPolicyBehaviorEnumEnumValues) + " Mirroring/copying requires an explicit `parent`; without one the server falls back to `none`.",
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedPolicyBehaviorEnumEnumValues),
			},
			"is_active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expiring": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"expires": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attributes": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				Description:      helpers.JSONDescription,
				DiffSuppressFunc: helpers.DiffSuppressJSON,
				ValidateDiagFunc: helpers.ValidateJSON,
			},
			// Computed
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"token_identifier": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifier of the agent's API token, so its key can be retrieved/copied later.",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The agent's API token. Only returned once, at creation.",
			},
		},
	}
}

func resourceAgentExpiresSet(d *schema.ResourceData, set func(*time.Time)) diag.Diagnostics {
	l, ok := d.Get("expires").(string)
	if !ok || l == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, l)
	if err != nil {
		return diag.FromErr(err)
	}
	set(&t)
	return nil
}

// resourceAgentSchemaToCreateModel builds the request for AgentsAgentsCreate, which only accepts
// the fields that can never be changed afterwards (parent, policy_behavior) plus expiry. It does
// not accept username/name/is_active/email/attributes - those are applied via a follow-up Update
// call in resourceAgentCreate, since AgentCreateRequest has no fields for them.
func resourceAgentSchemaToCreateModel(d *schema.ResourceData) (*api.AgentCreateRequest, diag.Diagnostics) {
	m := api.AgentCreateRequest{
		Parent:   helpers.GetIntP(d, "parent"),
		Expiring: new(d.Get("expiring").(bool)),
	}
	if pb, ok := d.GetOk("policy_behavior"); ok {
		m.PolicyBehavior = api.PolicyBehaviorEnum(pb.(string)).Ptr()
	}
	if diags := resourceAgentExpiresSet(d, m.Expires.Set); diags != nil {
		return nil, diags
	}
	return &m, nil
}

func resourceAgentSchemaToUpdateModel(d *schema.ResourceData) (*api.AgentRequest, diag.Diagnostics) {
	m := api.AgentRequest{
		Username: d.Get("username").(string),
		Name:     d.Get("name").(string),
		IsActive: new(d.Get("is_active").(bool)),
		Email:    helpers.GetP[string](d, "email"),
		Expiring: new(d.Get("expiring").(bool)),
	}
	attr, err := helpers.GetJSON[map[string]any](d, "attributes")
	if err != nil {
		return nil, err
	}
	m.Attributes = attr
	if diags := resourceAgentExpiresSet(d, m.Expires.Set); diags != nil {
		return nil, diags
	}
	return &m, nil
}

func resourceAgentCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	cm, diags := resourceAgentSchemaToCreateModel(d)
	if diags != nil {
		return diags
	}

	created, hr, err := c.client.AgentsAPI.AgentsAgentsCreate(ctx).AgentCreateRequest(*cm).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(created.Agent.Pk)))
	helpers.SetWrapper(d, "token", created.Token)

	um, diags := resourceAgentSchemaToUpdateModel(d)
	if diags != nil {
		return diags
	}
	if _, hr, err := c.client.AgentsAPI.AgentsAgentsUpdate(ctx, created.Agent.Pk).AgentRequest(*um).Execute(); err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	return resourceAgentRead(ctx, d, m)
}

func resourceAgentRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	res, hr, err := c.client.AgentsAPI.AgentsAgentsRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "username", res.Username)
	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "parent", res.Parent.Pk)
	helpers.SetWrapper(d, "policy_behavior", res.PolicyBehavior)
	helpers.SetWrapper(d, "is_active", res.IsActive)
	helpers.SetWrapper(d, "email", res.Email)
	helpers.SetWrapper(d, "expiring", res.Expiring)
	if res.Expires.IsSet() && res.Expires.Get() != nil {
		helpers.SetWrapper(d, "expires", res.Expires.Get().Format(time.RFC3339))
	}
	helpers.SetWrapper(d, "uid", res.Uid)
	helpers.SetWrapper(d, "uuid", res.Uuid)
	helpers.SetWrapper(d, "token_identifier", res.TokenIdentifier.Get())
	if diags := helpers.SetJSON(d, "attributes", res.Attributes); diags != nil {
		return diags
	}
	return diags
}

func resourceAgentUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	um, diags := resourceAgentSchemaToUpdateModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.AgentsAPI.AgentsAgentsUpdate(ctx, int32(id)).AgentRequest(*um).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceAgentRead(ctx, d, m)
}

func resourceAgentDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.AgentsAPI.AgentsAgentsDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
