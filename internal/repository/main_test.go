package repository

import (
	"basic-trade/pkg/config"
	"basic-trade/pkg/db"
	"context"
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

const useDocker = true

func TestMain(m *testing.M) {
	cfg := config.Database{
		Name:         "basic-trade-test",
		Host:         "localhost",
		Port:         2024,
		Password:     "secret",
		User:         "admin",
		Timezone:     "Asia/Jakarta",
		SslMode:      "disable",
		MigrationURL: "file://../../migration",
	}

	var test int
	if useDocker {
		setUpDocketTestEnv(cfg)
		test = m.Run()
		tearDownDockerTestEnv()
	} else {
		setUpDatabase(cfg)
		test = m.Run()
		tearDownDockerTestEnv()
	}

	os.Exit(test)
}

func setUpDocketTestEnv(cfg config.Database) {
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
		Tag:        "16-alpine",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", cfg.User),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", cfg.Password),
			fmt.Sprintf("POSTGRES_DB=%s", cfg.Name),
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
	cfg.Port = port

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

func setUpDatabase(cfg config.Database) error {
	var err error
	testDB, err = db.InitDB(cfg)
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

func tearDown(ctx context.Context) {
	tx, err := testDB.Begin(ctx)
	if err != nil {
		log.Fatal("Couldn't teardown: ", err)
	}

	defer tx.Rollback(ctx)

	tx.Exec(ctx, "DELETE FROM variants")
	tx.Exec(ctx, "DELETE FROM products")
	tx.Exec(ctx, "DELETE FROM admins")

	tx.Commit(ctx)
}
