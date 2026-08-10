package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

// TestPolicyGeoIPToRequest_BooleanCanBeTurnedOff is the executable spec for the bug
// migration-plan.md flags on this resource ("the main user of GetP[bool], so it is where the
// null-semantics change is user-visible").
//
// check_history_distance and check_impossible_travel are Optional with no Default. SDKv2
// built them with GetP[bool], which is d.GetOk underneath and so returned nil for any false
// value - it could not distinguish "unset" from "explicitly false". Both request fields are
// omitempty, and the API leaves absent fields unchanged on PUT, so a user who set either
// check to true and then back to false got: request omits the field, API keeps true, Read
// returns true, plan shows the same diff forever. There was no way to turn the check off
// through Terraform.
//
// The assertion that matters is the "explicitly false" row, and specifically that the key is
// present in the marshalled body - a non-nil pointer to false is necessary but not
// sufficient, since omitempty is what actually decided the outcome.
func TestPolicyGeoIPToRequest_BooleanCanBeTurnedOff(t *testing.T) {
	ctx := context.Background()
	r := &policyGeoIPResource{}

	base := func() *policyGeoIPModel {
		return &policyGeoIPModel{
			Name:                  types.StringValue("geoip"),
			ExecutionLogging:      types.BoolValue(false),
			Asns:                  types.ListNull(types.Int32Type),
			Countries:             types.ListNull(types.StringType),
			HistoryMaxDistanceKm:  types.Int64Value(100),
			DistanceToleranceKm:   types.Int32Value(50),
			HistoryLoginCount:     types.Int32Value(5),
			ImpossibleToleranceKm: types.Int32Value(100),
		}
	}

	t.Run("explicitly false is sent so the check can be turned off", func(t *testing.T) {
		data := base()
		data.CheckHistoryDistance = types.BoolValue(false)
		data.CheckImpossibleTravel = types.BoolValue(false)

		body, diags := r.toRequest(ctx, data)
		require.False(t, diags.HasError(), diags)

		require.NotNil(t, body.CheckHistoryDistance)
		assert.False(t, *body.CheckHistoryDistance)
		require.NotNil(t, body.CheckImpossibleTravel)
		assert.False(t, *body.CheckImpossibleTravel)

		// omitempty is what broke this before, so assert on the wire format.
		b, err := json.Marshal(body)
		require.NoError(t, err)
		assert.Contains(t, string(b), `"check_history_distance"`,
			"an explicit false must reach the API or the check can never be disabled")
		assert.Contains(t, string(b), `"check_impossible_travel"`)
	})

	t.Run("null is omitted so the server default stands", func(t *testing.T) {
		data := base()
		data.CheckHistoryDistance = types.BoolNull()
		data.CheckImpossibleTravel = types.BoolNull()

		body, diags := r.toRequest(ctx, data)
		require.False(t, diags.HasError(), diags)

		assert.Nil(t, body.CheckHistoryDistance)
		assert.Nil(t, body.CheckImpossibleTravel)

		b, err := json.Marshal(body)
		require.NoError(t, err)
		assert.NotContains(t, string(b), `"check_history_distance"`)
		assert.NotContains(t, string(b), `"check_impossible_travel"`)
	})

	t.Run("lists are sent as empty arrays rather than omitted", func(t *testing.T) {
		data := base()
		body, diags := r.toRequest(ctx, data)
		require.False(t, diags.HasError(), diags)

		// discovery #6: a nil slice is gated out by !IsNil, so clearing would not clear.
		require.NotNil(t, body.Asns)
		assert.Empty(t, body.Asns)
		require.NotNil(t, body.Countries)
		assert.Empty(t, body.Countries)
	})
}

// TestPolicyGeoIPFromAPI_NullVsDefault covers the read side of the same split: the two
// no-Default booleans must stay null when never configured, while everything with a Default
// takes the API value verbatim so import does not plan a spurious update (discovery #5).
func TestPolicyGeoIPFromAPI_NullVsDefault(t *testing.T) {
	ctx := context.Background()
	r := &policyGeoIPResource{}

	// Everything null, as after ImportState, and every API value its type's zero value.
	data := &policyGeoIPModel{
		ID:        types.StringValue("uuid-1"),
		Asns:      types.ListNull(types.Int32Type),
		Countries: types.ListNull(types.StringType),
	}

	res := &api.GeoIPPolicy{
		Pk:                    "uuid-1",
		Name:                  "geoip",
		ExecutionLogging:      new(false),
		CheckHistoryDistance:  new(false),
		CheckImpossibleTravel: new(false),
		HistoryMaxDistanceKm:  new(int64(0)),
		DistanceToleranceKm:   new(int32(0)),
		HistoryLoginCount:     new(int32(0)),
		ImpossibleToleranceKm: new(int32(0)),
	}

	diags := r.fromAPI(ctx, data, res)
	require.False(t, diags.HasError(), diags)

	assert.True(t, data.CheckHistoryDistance.IsNull(), "no Default, so null + false stays null")
	assert.True(t, data.CheckImpossibleTravel.IsNull(), "no Default, so null + false stays null")

	for name, v := range map[string]attrNullable{
		"execution_logging":       data.ExecutionLogging,
		"history_max_distance_km": data.HistoryMaxDistanceKm,
		"distance_tolerance_km":   data.DistanceToleranceKm,
		"history_login_count":     data.HistoryLoginCount,
		"impossible_tolerance_km": data.ImpossibleToleranceKm,
	} {
		assert.Falsef(t, v.IsNull(), "%s has a Default, so it must never import as null", name)
	}

	assert.True(t, data.Asns.IsNull(), "null in, nothing from the API, so null out")
	assert.True(t, data.Countries.IsNull())
}

// TestPolicyPasswordToRequest_DefaultedBooleansAlwaysSent is the sibling case: these three
// were also built with GetP[bool], but they have schema Defaults, so the fix is to send the
// concrete value unconditionally rather than to use BoolPtr. check_static_rules defaults to
// true, so before this it could not be switched off either.
func TestPolicyPasswordToRequest_DefaultedBooleansAlwaysSent(t *testing.T) {
	r := &policyPasswordResource{}

	body := r.toRequest(&policyPasswordModel{
		Name:                 types.StringValue("password"),
		ExecutionLogging:     types.BoolValue(false),
		PasswordField:        types.StringValue("password"),
		CheckStaticRules:     types.BoolValue(false),
		CheckHaveIBeenPwned:  types.BoolValue(false),
		CheckZxcvbn:          types.BoolValue(false),
		ErrorMessage:         types.StringValue(""),
		SymbolCharset:        types.StringValue("!"),
		HibpAllowedCount:     types.Int32Value(1),
		ZxcvbnScoreThreshold: types.Int32Value(2),
		AmountUppercase:      types.Int32Null(),
		AmountLowercase:      types.Int32Null(),
		AmountSymbols:        types.Int32Null(),
		AmountDigits:         types.Int32Null(),
		LengthMin:            types.Int32Null(),
	})

	require.NotNil(t, body.CheckStaticRules)
	assert.False(t, *body.CheckStaticRules, "check_static_rules defaults to true, so false must be sent to disable it")
	require.NotNil(t, body.CheckHaveIBeenPwned)
	require.NotNil(t, body.CheckZxcvbn)

	b, err := json.Marshal(body)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"check_static_rules"`)

	// error_message is Required in Terraform but optional on the wire; SDKv2 dropped an
	// intentionally-empty message via GetP.
	require.NotNil(t, body.ErrorMessage)
	assert.Equal(t, "", *body.ErrorMessage)

	// The five no-Default rules are omitted when null, which is the correct H1 behaviour.
	assert.Nil(t, body.AmountUppercase)
	assert.Nil(t, body.LengthMin)
}
