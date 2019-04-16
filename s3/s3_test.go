package s3

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/Ryanair/goaws"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

var (
	cli      *Client
	bucketID string
)

func TestClient_GeneratePutURL_ok(t *testing.T) {
	// when
	url, err := cli.GeneratePutURL(bucketID, "some_random_key", "", 30*time.Minute)
	if err != nil {
		t.Fatalf("failed to generate url: %s", err)
	}

	// then
	assert.NotEmpty(t, url)
}

func TestClient_GeneratePutURL_signingFailed(t *testing.T) {
	// when
	url, genErr := cli.GeneratePutURL(bucketID, "some_random_key", "", -30*time.Minute)

	// then
	isSigningFailed := func(err error) bool {
		type signingFailed interface {
			SigningFailed() bool
		}
		e, ok := err.(signingFailed)
		return ok && e.SigningFailed()
	}

	assert.Empty(t, url)
	assert.True(t, isSigningFailed(genErr))
}

func TestClient_GetObject_keyNotFound(t *testing.T) {
	// given & when
	_, getErr := cli.GetObject(bucketID, "some_random_key")

	// then
	isKeyNotFound := func(err error) bool {
		type keyNotFound interface {
			KeyNotFound() bool
		}
		e, ok := err.(keyNotFound)
		return ok && e.KeyNotFound()
	}

	assert.True(t, isKeyNotFound(getErr))
}

func TestMain(m *testing.M) {
	setupS3()
	code := m.Run()
	teardownS3()

	os.Exit(code)
}

func setupS3() {
	config, err := goaws.NewConfig(goaws.Region(endpoints.EuWest1RegionID))
	if err != nil {
		log.Fatalf("couldn't create config: %s", err)
	}

	cli = NewClient(config)
	bucketID = xid.New().String()

	if _, err := cli.s3.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketID),
	}); err != nil {
		if err == credentials.ErrNoValidProvidersFoundInChain {
			log.Fatalf("Test failed due to lack of AWS credentials in chain.")
		}
		log.Fatalf("couldn't create bucket: %s", err)
	}
}

func teardownS3() {
	if _, err := cli.s3.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketID),
	}); err != nil {
		log.Fatalf("couldn't delete bucket: %s", err)
	}
}
