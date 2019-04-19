package nconf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenericConfigLoad(t *testing.T) {
	gc := GenericConfig{
		Type: "something",
		Config: map[string]interface{}{
			"string": "string",
			"int":    123,
			"nested": struct {
				NV string
			}{
				NV: "nested-value",
			},
		},
	}

	result := struct {
		String string
		Int    int
		Nested struct {
			NV string
		}
	}{}
	require.NoError(t, gc.Load(&result))
	assert.Equal(t, result.String, "string")
	assert.Equal(t, result.Int, 123)
	assert.Equal(t, result.Nested.NV, "nested-value")
}

func TestGenericConfigLoadWithDefaults(t *testing.T) {
	gc := GenericConfig{
		Type: "something",
		Config: map[string]interface{}{
			"int": 123,
			"nested": struct {
				NV string
			}{
				NV: "nested-value",
			},
		},
	}

	result := struct {
		String string
		Int    int
		Nested struct {
			NV string
		}
	}{
		String: "default-string",
	}
	require.NoError(t, gc.Load(&result))
	assert.Equal(t, result.String, "default-string")
	assert.Equal(t, result.Int, 123)
	assert.Equal(t, result.Nested.NV, "nested-value")
}
