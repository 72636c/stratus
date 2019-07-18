package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/72636c/stratus/internal/config"
)

func Test_CalculateChecksum(t *testing.T) {
	stackInput := func(stack *config.Stack) func() interface{} {
		return func() interface{} {
			return stack.Hashable()
		}
	}

	testCases := []struct {
		description string
		input       func() interface{}
		expected    string
	}{
		{
			description: "full config.Stack with checksum fields",
			input: stackInput(&config.Stack{
				Name:         "a",
				Capabilities: []string{"b1", "b2"},
				Parameters: config.StackParameters{
					{
						Key:   "c1",
						Value: "c2",
					},
					{
						Key:   "c3",
						Value: "c4",
					},
				},
				Tags: config.StackTags{
					{
						Key:   "d1",
						Value: "d2",
					},
					{
						Key:   "d3",
						Value: "d4",
					},
				},
				TerminationProtection: true,

				Policy:   []byte("e"),
				Template: []byte("f"),

				ArtefactBucket: "g",
				PolicyKey:      "h",
				TemplateKey:    "i",

				Checksum: "j",
			}),
			expected: "4d76060a71bb68af4c917a48f3a09298ba9a7cab8f8890be39cc867b4f20645f",
		},

		{
			description: "full config.Stack without checksum fields",
			input: stackInput(&config.Stack{
				Name:         "a",
				Capabilities: []string{"b1", "b2"},
				Parameters: config.StackParameters{
					{
						Key:   "c1",
						Value: "c2",
					},
					{
						Key:   "c3",
						Value: "c4",
					}},
				Tags: config.StackTags{
					{
						Key:   "d1",
						Value: "d2",
					},
					{
						Key:   "d3",
						Value: "d4",
					},
				},
				TerminationProtection: true,

				Policy:   []byte("e"),
				Template: []byte("f"),

				ArtefactBucket: "g",
			}),
			expected: "4d76060a71bb68af4c917a48f3a09298ba9a7cab8f8890be39cc867b4f20645f",
		},
		{
			description: "minimal config.Stack",
			input: stackInput(&config.Stack{
				Name:                  "a",
				Capabilities:          make([]string, 0),
				Parameters:            make(config.StackParameters, 0),
				Tags:                  make(config.StackTags, 0),
				TerminationProtection: false,

				Policy:   []byte("e"),
				Template: []byte("f"),

				ArtefactBucket: "",
			}),
			expected: "e5fcf031009c4780d9993fa6afcee53dd1079a513b9033176ffa8f61f148fc4d",
		},
		{
			description: "nil config.Stack",
			input:       stackInput(new(config.Stack)),
			expected:    "1245e875630639894025da5af9f4bc530eda6e824b84c362ac190b4cca141d0a",
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
