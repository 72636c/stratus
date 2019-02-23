package stratus

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/mock"
)

type S3Mock struct {
	mock.Mock
}

func NewS3Mock() *S3Mock {
	return new(S3Mock)
}

func (client *S3Mock) PutObjectWithContext(
	_ aws.Context,
	input *s3.PutObjectInput,
	_ ...request.Option,
) (*s3.PutObjectOutput, error) {
	args := client.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}

func (client *S3Mock) PutObjectWithContextMatcher(
	expected *s3.PutObjectInput,
) func(*s3.PutObjectInput) bool {
	return func(actual *s3.PutObjectInput) bool {
		if expected == nil && actual == nil {
			return true
		}

		if expected == nil ||
			actual == nil ||
			!reflect.DeepEqual(expected.Bucket, actual.Bucket) ||
			!reflect.DeepEqual(expected.Key, actual.Key) {
			fmt.Printf(
				"S3Mock.PutObjectWithContextMatcher: expected '%+v', received '%+v'\n",
				expected,
				actual,
			)
			return false
		}

		err := compareReaders(expected.Body, actual.Body)
		if err != nil {
			fmt.Printf("S3Mock.PutObjectWithContextMatcher.Body: %+v\n", err)
			return false
		}

		return true
	}
}

func compareReaders(a, b io.Reader) error {
	var bufferA, bufferB bytes.Buffer

	teeA := io.TeeReader(a, &bufferA)
	teeB := io.TeeReader(b, &bufferB)

	dataA, err := ioutil.ReadAll(teeA)
	if err != nil {
		return fmt.Errorf("error reading a: %+v", err)
	}

	dataB, err := ioutil.ReadAll(teeB)
	if err != nil {
		return fmt.Errorf("error reading b: %+v", err)
	}

	equal := reflect.DeepEqual(dataA, dataB)
	if !equal {
		return fmt.Errorf("expected '%s', received '%s'", dataA, dataB)
	}

	return nil
}
