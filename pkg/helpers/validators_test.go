package helpers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestRelativeDuration_ValidateString(t *testing.T) {
	testCases := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"single key", "hours=1", false},
		{"multiple keys", "hours=1;minutes=2;seconds=3", false},
		{"negative value", "minutes=-5", false},
		{"case insensitive key", "HOURS=1", false},
		{"unknown key", "fortnights=1", true},
		{"missing equals", "hours", true},
	}

	v := RelativeDuration()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := validator.StringRequest{ConfigValue: types.StringValue(tc.value)}
			resp := &validator.StringResponse{}
			v.ValidateString(context.Background(), req, resp)
			assert.Equal(t, tc.wantErr, resp.Diagnostics.HasError())
		})
	}
}

func TestRelativeDuration_NullUnknownSkipped(t *testing.T) {
	v := RelativeDuration()

	resp := &validator.StringResponse{}
	v.ValidateString(context.Background(), validator.StringRequest{ConfigValue: types.StringNull()}, resp)
	assert.False(t, resp.Diagnostics.HasError())

	resp = &validator.StringResponse{}
	v.ValidateString(context.Background(), validator.StringRequest{ConfigValue: types.StringUnknown()}, resp)
	assert.False(t, resp.Diagnostics.HasError())
}

type testEnum string

const (
	testEnumFoo testEnum = "foo"
	testEnumBar testEnum = "bar"
)

func TestOneOf(t *testing.T) {
	v := OneOf([]testEnum{testEnumFoo, testEnumBar})

	resp := &validator.StringResponse{}
	v.ValidateString(context.Background(), validator.StringRequest{ConfigValue: types.StringValue("foo")}, resp)
	assert.False(t, resp.Diagnostics.HasError())

	resp = &validator.StringResponse{}
	v.ValidateString(context.Background(), validator.StringRequest{ConfigValue: types.StringValue("nope")}, resp)
	assert.True(t, resp.Diagnostics.HasError())
}
