package stratus

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	_ S3 = new(s3.S3)
	_ S3 = new(S3Mock)
)

type S3 interface {
	PutObjectWithContext(
		aws.Context,
		*s3.PutObjectInput,
		...request.Option,
	) (*s3.PutObjectOutput, error)
}
