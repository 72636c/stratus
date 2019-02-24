package config

type ArtefactPath struct {
	Section    string
	Resource   string
	Property   string
	Expression string
}

var (
	localArtefactPaths = []*ArtefactPath{
		&ArtefactPath{
			Section:  "Resources",
			Resource: "AWS::ApiGateway::RestApi",
			Property: "BodyS3Location",
		},
		&ArtefactPath{
			Section:  "Resources",
			Resource: "AWS::AppSync::GraphQLSchema",
			Property: "DefinitionS3Location",
		},
		&ArtefactPath{
			Section:  "Resources",
			Resource: "AWS::AppSync::Resolver",
			Property: "RequestMappingTemplateS3Location",
		},
		&ArtefactPath{
			Section:  "Resources",
			Resource: "AWS::AppSync::Resolver",
			Property: "ResponseMappingTemplateS3Location",
		},
		&ArtefactPath{
			Section:  "Resources",
			Resource: "AWS::CloudFormation::Stack",
			Property: "TemplateURL",
		},
		&ArtefactPath{
			Section:  "Resources",
			Resource: "AWS::ElasticBeanstalk::ApplicationVersion",
			Property: "SourceBundle",
		},
		&ArtefactPath{
			Section:  "Resources",
			Resource: "AWS::Lambda::Function",
			Property: "Code",
		},
		&ArtefactPath{
			Section:  "Resources",
			Resource: "AWS::Serverless::Api",
			Property: "DefinitionUri",
		},
		&ArtefactPath{
			Section:  "Resources",
			Resource: "AWS::Serverless::Function",
			Property: "CodeUri",
		},
		// TODO: this can be anywhere
		&ArtefactPath{
			Section:  "Transform",
			Resource: "AWS::Include",
			Property: "Location",
		},
	}
)

var (
	a = `{
		"Name": "AWS::Include",
		"Parameters": {
			"Location": "{{PATH}}"
		}
	}`
	b = `{
		"Resources": {
			"{{ANY}}": {
				"Name": "AWS::Include",
				"Parameters": {
					"Location": "{{PATH}}"
				}
			}
		}
	}`
)
