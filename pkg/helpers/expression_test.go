package helpers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpressionValue_StringSemanticEquals(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name  string
		prior string
		new   string
		equal bool
	}{
		{"identical", "return True", "return True", true},
		{"heredoc trailing newline vs API-stripped", "return True\n", "return True", true},
		{"both sides have trailing newlines", "return True\n", "return True\n", true},
		{"genuinely different", "return True", "return False", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prior := ExpressionValue{StringValue: types.StringValue(tc.prior)}
			new := ExpressionValue{StringValue: types.StringValue(tc.new)}

			equal, diags := prior.StringSemanticEquals(ctx, new)
			require.False(t, diags.HasError())
			assert.Equal(t, tc.equal, equal)

			// symmetric: must not matter which side is "prior" vs "new"
			equalReverse, diags := new.StringSemanticEquals(ctx, prior)
			require.False(t, diags.HasError())
			assert.Equal(t, tc.equal, equalReverse)
		})
	}
}
