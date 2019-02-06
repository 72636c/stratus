package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/72636c/stratus/internal/config"
)

func Test_Resolve(t *testing.T) {
	os.Setenv("SET_1", "serious")
	os.Setenv("SET_2", "prod")
	os.Setenv("SET_3", "env")
	os.Setenv("SET_4", "SET_1")
	os.Unsetenv("UNSET_1")

	testCases := []struct {
		description   string
		input         string
		expected      string
		expectedError string
	}{
		{
			description: "no placeholders",
			input:       "serious-resource-name",
			expected:    "serious-resource-name",
		},
		{
			description: "unhandled single brackets",
			input:       "serious-resource-{env:SET_1}",
			expected:    "serious-resource-{env:SET_1}",
		},
		{
			description: "unhandled opening bracket prefix and suffix",
			input:       "{-{",
			expected:    "{-{",
		},
		{
			description: "unhandled closing bracket prefix and suffix",
			input:       "}-}",
			expected:    "}-}",
		},
		{
			description: "environment placeholders",
			input:       "{{env:SET_1}}-resource-name-{{env:SET_2}}",
			expected:    "serious-resource-name-prod",
		},
		{
			description: "nested environment placeholders",
			input:       "{{{{env:SET_3}}:{{env:SET_4}}}}-resource-name",
			expected:    "serious-resource-name",
		},
		{
			description:   "environment variable not set",
			input:         "{{env:UNSET_1}}",
			expected:      "",
			expectedError: "not set",
		},
		{
			description:   "malformed placeholder",
			input:         "{{env_UNSET_1}}",
			expected:      "",
			expectedError: "malformed placeholder",
		},
		{
			description:   "unhandled placeholder",
			input:         "{{gcp:SET_1}}",
			expected:      "",
			expectedError: "unrecognised placeholder",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			actual, err := config.Resolve(testCase.input)
			assert.Equal(testCase.expected, actual)
			if testCase.expectedError == "" {
				assert.NoError(err)
			} else {
				require.Error(err)
				assert.Contains(err.Error(), testCase.expectedError)
			}
		})
	}
}

func Test_Unmarshal_JSON(t *testing.T) {
	os.Setenv("SET_1", "true")

	type Thing struct {
		Bool   config.Bool
		String config.String
	}

	testCases := []struct {
		description   string
		input         string
		expected      *Thing
		expectedError string
	}{
		{
			description: "literals",
			input:       `{"Bool": true, "String": "hello"}`,
			expected: &Thing{
				Bool:   config.Bool(true),
				String: "hello",
			},
		},
		{
			description: "placeholders",
			input:       `{"Bool": "{{env:SET_1}}", "String": "{{env:SET_1}}"}`,
			expected: &Thing{
				Bool:   config.Bool(true),
				String: "true",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			var thing *Thing

			err := config.Unmarshal(".json", []byte(testCase.input), &thing)
			assert.Equal(testCase.expected, thing)
			if testCase.expectedError == "" {
				assert.NoError(err)
			} else {
				require.Error(err)
				assert.Contains(err.Error(), testCase.expectedError)
			}
		})
	}
}

func Test_Unmarshal_YAML(t *testing.T) {
	os.Setenv("SET_1", "on")

	type Thing struct {
		Bool   config.Bool
		String config.String
	}

	testCases := []struct {
		description   string
		input         string
		expected      *Thing
		expectedError string
	}{
		{
			description: "literals",
			input:       "bool: yes\nstring: hello",
			expected: &Thing{
				Bool:   config.Bool(true),
				String: "hello",
			},
		},
		{
			description: "placeholders",
			input:       "bool: '{{env:SET_1}}'\nstring: '{{env:SET_1}}'",
			expected: &Thing{
				Bool:   config.Bool(true),
				String: "on",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			var thing *Thing

			err := config.Unmarshal(".yaml", []byte(testCase.input), &thing)
			assert.Equal(testCase.expected, thing)
			if testCase.expectedError == "" {
				assert.NoError(err)
			} else {
				require.Error(err)
				assert.Contains(err.Error(), testCase.expectedError)
			}
		})
	}
}
