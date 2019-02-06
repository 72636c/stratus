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
			expectedError: "not set",
		},
		{
			description:   "malformed placeholder",
			input:         "{{env_UNSET_1}}",
			expectedError: "malformed placeholder",
		},
		{
			description:   "unsupported placeholder",
			input:         "{{gcp:SET_1}}",
			expectedError: "unsupported placeholder",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			actual, err := config.Resolve(testCase.input)
			if testCase.expectedError == "" {
				assert.Equal(testCase.expected, actual)
				assert.NoError(err)
			} else {
				require.Error(err)
				assert.Contains(err.Error(), testCase.expectedError)
			}
		})
	}
}

func Test_Unmarshal(t *testing.T) {
	os.Setenv("SET_1", "true")

	type Data struct {
		Bool   config.Bool
		String config.String
	}

	testCases := []struct {
		description   string
		input         string
		extension     string
		expected      *Data
		expectedError string
	}{
		{
			description: "JSON literals",
			input:       `{"Bool": true, "String": "hello"}`,
			extension:   ".json",
			expected: &Data{
				Bool:   true,
				String: "hello",
			},
		},
		{
			description: "JSON placeholders",
			input:       `{"Bool": "{{env:SET_1}}", "String": "{{env:SET_1}}"}`,
			extension:   ".json",
			expected: &Data{
				Bool:   true,
				String: "true",
			},
		},
		{
			description: "YAML literals",
			input:       "bool: yes\nstring: hello",
			extension:   ".yaml",
			expected: &Data{
				Bool:   true,
				String: "hello",
			},
		},
		{
			description: "YAML placeholders",
			input:       "bool: '{{env:SET_1}}'\nstring: '{{env:SET_1}}'",
			extension:   ".yaml",
			expected: &Data{
				Bool:   true,
				String: "true",
			},
		},
		{
			description:   "invalid JSON bool",
			input:         `{"Bool": "tr\"ue", "String": "hello"}`,
			extension:     ".json",
			expectedError: "invalid character",
		},
		{
			description:   "invalid YAML bool",
			input:         "bool: y\"e\"s\nstring: hello",
			extension:     ".yaml",
			expectedError: "cannot unmarshal",
		},
		{
			description:   "unsupported file extension",
			input:         "",
			extension:     ".xml",
			expectedError: "unsupported file extension",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			var data *Data

			err := config.Unmarshal(testCase.extension, []byte(testCase.input), &data)
			if testCase.expectedError == "" {
				assert.Equal(testCase.expected, data)
				assert.NoError(err)
			} else {
				require.Error(err)
				assert.Contains(err.Error(), testCase.expectedError)
			}
		})
	}
}
