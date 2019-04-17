package internal

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/pkg/errors"
)

type DockerImage struct {
	Repo string
	Tag  string
	Env  []string
}

func DockerSetup(m *testing.M, img DockerImage, setup func(*dockertest.Resource) error, purge ...func(*dockertest.Resource) error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run(img.Repo, img.Tag, img.Env)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		if err := setup(resource); err != nil {
			return errors.Wrap(err, "Could not setup resource")
		}
		time.Sleep(100 * time.Millisecond)

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	for _, p := range purge {
		if err := p(resource); err != nil {
			log.Printf("Could not purge resource: %s", err)
		}
	}

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge pool resource: %s", err)
	}

	os.Exit(code)
}
