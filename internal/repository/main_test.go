package repository

import (
	"basic-trade/pkg/config"
	"basic-trade/pkg/db"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

var (
	testDB          *pgxpool.Pool
	testAdminRepo   AdminRepository
	testProductRepo ProductRepository
	testVariantRepo VariantRepository
	pool            *dockertest.Pool
	resource        *dockertest.Resource
)

const useDocker = false

func TestMain(m *testing.M) {
	cfg := config.LoadTestConfig("../../app.yaml")

	var test int
	if useDocker {
		setUpDocketTestEnv(cfg)
		test = m.Run()
		tearDownDockerTestEnv()
	} else {
		setUpDatabase(cfg)
		test = m.Run()
	}

	os.Exit(test)
}

func setUpDocketTestEnv(cfg config.Config) {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatal("Couldn't construct pool: ", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatal("Couldn't connect to Docker: ", err)
	}

	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", cfg.Database.User),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", cfg.Database.Password),
			fmt.Sprintf("POSTGRES_DB=%s", cfg.Database.Name),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		log.Fatal("Couldn't start resource: ", err)
	}

	rcPort := resource.GetPort("5432/tcp")
	port, err := strconv.Atoi(rcPort)
	if err != nil {
		log.Fatal("Couldn't set port: ", err)
	}
	cfg.Database.Port = port

	// Takes a few seconds to start up
	// Increase the delay if it fails to connect
	time.Sleep(3 * time.Second)

	if err := pool.Retry(func() error {
		var err error
		err = setUpDatabase(cfg)

		if err != nil {
			for retries := 1; retries <= 3; retries++ {
				time.Sleep(5 * time.Second)
				err = setUpDatabase(cfg)

				if err == nil {
					break
				}
			}
		}

		return err
	}); err != nil {
		tearDownDockerTestEnv()
		log.Fatal("Couldn't connect to DB: ", err)
	}
}

func setUpDatabase(cfg config.Config) error {
	var err error
	testDB, err = db.InitDB(cfg.Database)
	if err != nil {
		return err
	}

	testAdminRepo = NewAdminRepository(testDB)
	testProductRepo = NewProductRepository(testDB)
	testVariantRepo = NewVariantRepository(testDB)

	return nil
}

func tearDownDockerTestEnv() {
	if err := pool.Purge(resource); err != nil {
		log.Fatal("Could not purge Docker: ", err)
	}
}
