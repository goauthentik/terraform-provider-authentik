package provider

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceOutpostProviderAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceOutpostProviderAttachmentCreate,
		ReadContext:   resourceOutpostProviderAttachmentRead,
		DeleteContext: resourceOutpostProviderAttachmentDelete,
		UpdateContext: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"outpost": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the outpost.",
			},
			"protocol_provider": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the provider.",
			},
		},
	}
}

func resourceOutpostProviderAttachmentCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	outpostID := d.Get("outpost").(string)
	providerID := int32(d.Get("protocol_provider").(int))

	// Get current outpost
	outpost, hr, err := c.client.OutpostsApi.OutpostsInstancesRetrieve(ctx, outpostID).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	// Check if provider is already attached
	if slices.Contains(outpost.Providers, providerID) {
		// Already attached, just set ID
		d.SetId(fmt.Sprintf("%s:%d", outpostID, providerID))
		return resourceOutpostProviderAttachmentRead(ctx, d, m)
	}

	// Add provider
	outpost.Providers = append(outpost.Providers, providerID)

	// Update outpost
	req := api.PatchedOutpostRequest{
		Providers: outpost.Providers,
	}

	_, hr, err = c.client.OutpostsApi.OutpostsInstancesPartialUpdate(ctx, outpostID).PatchedOutpostRequest(req).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(fmt.Sprintf("%s:%d", outpostID, providerID))
	return resourceOutpostProviderAttachmentRead(ctx, d, m)
}

func resourceOutpostProviderAttachmentRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	// Parse ID
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return diag.Errorf("Invalid ID format")
	}
	outpostID := parts[0]
	// providerID is parts[1] but we need to check if it exists in the outpost

	outpost, hr, err := c.client.OutpostsApi.OutpostsInstancesRetrieve(ctx, outpostID).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	providerIDInt, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	found := slices.Contains(outpost.Providers, int32(providerIDInt))

	if !found {
		d.SetId("")
		return nil
	}

	helpers.SetWrapper(d, "outpost", outpostID)
	helpers.SetWrapper(d, "protocol_provider", int(providerIDInt))

	return nil
}

func resourceOutpostProviderAttachmentDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	outpostID := d.Get("outpost").(string)
	providerID := int32(d.Get("protocol_provider").(int))

	// Get current outpost
	outpost, hr, err := c.client.OutpostsApi.OutpostsInstancesRetrieve(ctx, outpostID).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	// Remove provider
	newProviders := []int32{}
	for _, p := range outpost.Providers {
		if p != providerID {
			newProviders = append(newProviders, p)
		}
	}

	// Update outpost
	req := api.PatchedOutpostRequest{
		Providers: newProviders,
	}

	_, hr, err = c.client.OutpostsApi.OutpostsInstancesPartialUpdate(ctx, outpostID).PatchedOutpostRequest(req).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	return nil
}
