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
	sc.CreateContext = func(ctx context.Context, rd *schema.ResourceData, m any) diag.Diagnostics {
		span := sentry.StartSpan(ctx, "terraform.resource.create", sentry.WithTransactionName("terraform.resource"))
		span.Description = "Resource create"
		defer span.Finish()
		return so.CreateContext(ctx, rd, m)
	}
	sc.ReadContext = func(ctx context.Context, rd *schema.ResourceData, m any) diag.Diagnostics {
		span := sentry.StartSpan(ctx, "terraform.resource.read", sentry.WithTransactionName("terraform.resource"))
		span.Description = "Resource read"
		defer span.Finish()
		return so.ReadContext(ctx, rd, m)
	}
	if so.UpdateContext != nil {
		sc.UpdateContext = func(ctx context.Context, rd *schema.ResourceData, m any) diag.Diagnostics {
			span := sentry.StartSpan(ctx, "terraform.resource.update", sentry.WithTransactionName("terraform.resource"))
			span.Description = "Resource update"
			defer span.Finish()
			return so.UpdateContext(ctx, rd, m)
		}
	}
	sc.DeleteContext = func(ctx context.Context, rd *schema.ResourceData, m any) diag.Diagnostics {
		span := sentry.StartSpan(ctx, "terraform.resource.delete", sentry.WithTransactionName("terraform.resource"))
		span.Description = "Resource delete"
		defer span.Finish()
		return so.DeleteContext(ctx, rd, m)
	}
	return sc
}

func td(resource func() *schema.Resource) *schema.Resource {
	sc := resource()
	so := resource()
	sc.ReadContext = func(ctx context.Context, rd *schema.ResourceData, m any) diag.Diagnostics {
		span := sentry.StartSpan(ctx, "terraform.datasource.read", sentry.WithTransactionName("terraform.datasource"))
		span.Description = "Datasource read"
		defer span.Finish()
		return so.ReadContext(ctx, rd, m)
	}
	return sc
}
