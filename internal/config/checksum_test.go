package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/72636c/stratus/internal/config"
)

func Test_CalculateChecksum(t *testing.T) {
	testCases := []struct {
		description string
		input       func() interface{}
		expected    string
	}{
		{
			description: "config.Stack",
			input: func() interface{} {
				return &config.Stack{
					Policy:   []byte("a"),
					Template: []byte("b"),
				}
			},
			expected: "5dc046d563a19c16dfae96d8b530873f7b6a758af3e68bde7aefa3d96790770c",
		},
		{
			description: "map",
			input: func() interface{} {
				return map[string]interface{}{
					"c": -1,
					"a": "hello",
					"b": true,
				}
			},
			expected: "a43ce2809f97260129f7bfd1352033ced8379cb4bff7a87a606bc7b136ba496d",
		},
		{
			description: "struct",
			input: func() interface{} {
				return struct {
					C int    `json:"c"`
					A string `json:"a"`
					B bool   `json:"b"`
				}{
					C: -1,
					A: "hello",
					B: true,
				}
			},
			expected: "481dadb27dfdb867222436adf2b1ae52715838e21fe1cae7ae9ebfbd16286690",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			assert := assert.New(t)

			for index := 0; index < 100; index++ {
				actual, err := config.CalculateChecksum(testCase.input())
				passed1 := assert.Equal(testCase.expected, actual)
				passed2 := assert.NoError(err)

				if !(passed1 && passed2) {
					break
				}
			}
		})
	}
}
