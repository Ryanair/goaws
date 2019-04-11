package s3

import "github.com/aws/aws-sdk-go/service/s3"

type Client struct {
	*s3.S3
}

func NewClient() {

}
