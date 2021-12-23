package provider

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func tr(resource func() *schema.Resource) *schema.Resource {
	sc := resource()
	so := resource()
	sc.CreateContext = func(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
		span := sentry.StartSpan(ctx, "authentik.terraform.resource.create", sentry.TransactionName("authentik.terraform.resource.create"))
		span.Description = "Resource create"
		defer span.Finish()
		return so.CreateContext(ctx, rd, m)
	}
	sc.ReadContext = func(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
		span := sentry.StartSpan(ctx, "authentik.terraform.resource.read", sentry.TransactionName("authentik.terraform.resource.read"))
		span.Description = "Resource read"
		defer span.Finish()
		return so.ReadContext(ctx, rd, m)
	}
	sc.UpdateContext = func(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
		span := sentry.StartSpan(ctx, "authentik.terraform.resource.update", sentry.TransactionName("authentik.terraform.resource.update"))
		span.Description = "Resource update"
		defer span.Finish()
		return so.UpdateContext(ctx, rd, m)
	}
	sc.DeleteContext = func(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
		span := sentry.StartSpan(ctx, "authentik.terraform.resource.delete", sentry.TransactionName("authentik.terraform.resource.delete"))
		span.Description = "Resource delete"
		defer span.Finish()
		return so.DeleteContext(ctx, rd, m)
	}
	return sc
}

func td(resource func() *schema.Resource) *schema.Resource {
	sc := resource()
	so := resource()
	sc.ReadContext = func(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
		span := sentry.StartSpan(ctx, "authentik.terraform.datasource.read", sentry.TransactionName("authentik.terraform.datasource.read"))
		span.Description = "Datasource read"
		defer span.Finish()
		return so.ReadContext(ctx, rd, m)
	}
	return sc
}
