// +build local

package s3

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Ryanair/goaws"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func TestClient_GeneratePutURLWithMetadata_ok(t *testing.T) {
	// when
	meta := map[string]*string{"filename": aws.String("HappyFace.jpg")}
	url, err := cli.GeneratePutURLWithMetadata(bucketID, "some_random_key", "", 30*time.Minute, meta)
	if err != nil {
		t.Fatalf("failed to generate url with metada: %s", err)
	}

	// then
	assert.NotEmpty(t, url)
}

func TestClient_GeneratePutURLWithMetadata_signingFailed(t *testing.T) {
	// when
	meta := map[string]*string{"filename": aws.String("SadFace.png")}
	url, genErr := cli.GeneratePutURLWithMetadata(bucketID, "some_random_key", "", -30*time.Minute, meta)

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

func TestClient_PutObject_ok(t *testing.T) {
	// given & when
	putErr := cli.PutObject(bucketID, "some_random_key", bytes.NewReader([]byte("abc")))

	// then
	out, getErr := cli.GetObject(bucketID, "some_random_key")
	savedItem, readErr := ioutil.ReadAll(out)

	assert.NoError(t, putErr)
	assert.NoError(t, getErr)
	assert.NoError(t, readErr)
	assert.Equal(t, "abc", string(savedItem))
}

func TestClient_GetObject_keyNotFound(t *testing.T) {
	// given & when
	_, getErr := cli.GetObject(bucketID, "non_existing_key")

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
		log.Fatalf("couldn't create bucket %s due to: %s", bucketID, err)
	}
}

func teardownS3() {
	iter := s3manager.NewDeleteListIterator(cli.s3, &s3.ListObjectsInput{
		Bucket: aws.String(bucketID),
	})

	if err := s3manager.NewBatchDeleteWithClient(cli.s3).Delete(aws.BackgroundContext(), iter); err != nil {
		log.Fatalf("couldn't delete objects from bucket %s due to: %s", bucketID, err)
	}

	if _, err := cli.s3.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketID),
	}); err != nil {
		log.Fatalf("couldn't delete bucket %s due to: %s", bucketID, err)
	}
}
