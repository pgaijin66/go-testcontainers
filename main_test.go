package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGetUsersHandler(t *testing.T) {
	ctx := context.Background()

	// Start a PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "user",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections"),
		),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	defer func() {
		err := pgContainer.Terminate(ctx)
		assert.NoError(t, err)
	}()

	mappedPort, err := pgContainer.MappedPort(ctx, "5432")
	assert.NoError(t, err)

	host, err := pgContainer.Host(ctx)
	assert.NoError(t, err)

	dbUrl := fmt.Sprintf("postgres://user:password@%s:%s/testdb", host, mappedPort.Port())

	// Connect to the PostgreSQL container
	var conn *pgx.Conn
	for i := 0; i < 5; i++ {
		conn, err = pgx.Connect(ctx, dbUrl)
		if err == nil {
			break
		}
	}
	assert.NoError(t, err)
	defer func() {
		if conn != nil {
			err := conn.Close(ctx)
			assert.NoError(t, err)
		}
	}()

	// Create the "users" table and insert test data
	_, err = conn.Exec(ctx, "CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(50));")
	assert.NoError(t, err)

	_, err = conn.Exec(ctx, "INSERT INTO users (name) VALUES ('Alice'), ('Bob');")
	assert.NoError(t, err)

	// Assign the global db variable for handler to use
	db = conn

	// Setup the HTTP request and response recorder
	reqHttp, err := http.NewRequest("GET", "/users", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUsersHandler)
	handler.ServeHTTP(rr, reqHttp)

	// Check the status code is what we expect
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body is what we expect
	var response map[string]interface{}
	body, err := ioutil.ReadAll(rr.Body)
	assert.NoError(t, err)

	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Define the expected response data
	expectedData := []interface{}{
		map[string]interface{}{"id": float64(1), "name": "Alice"},
		map[string]interface{}{"id": float64(2), "name": "Bob"},
	}

	// Assert the JSON response
	assert.Equal(t, "success", response["message"])
	assert.Equal(t, float64(200), response["status_code"])
	assert.Equal(t, expectedData, response["data"])
}
