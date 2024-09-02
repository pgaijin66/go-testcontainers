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
	testCases := []struct {
		name          string
		setupDatabase func(conn *pgx.Conn, ctx context.Context) error
		expectedData  []map[string]interface{}
	}{
		{
			name: "Handle multiple users",
			setupDatabase: func(conn *pgx.Conn, ctx context.Context) error {
				_, err := conn.Exec(ctx, "CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(50));")
				if err != nil {
					return err
				}
				_, err = conn.Exec(ctx, "INSERT INTO users (name) VALUES ('Alice'), ('Bob');")
				return err
			},
			expectedData: []map[string]interface{}{
				{"id": float64(1), "name": "Alice"},
				{"id": float64(2), "name": "Bob"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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

			// Setup database and insert test data
			err = testCase.setupDatabase(conn, ctx)
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

			// Convert the "data" field from []interface{} to []map[string]interface{} for comparison
			actualData, ok := response["data"].([]interface{})
			assert.True(t, ok)

			var actualDataConverted []map[string]interface{}
			for _, item := range actualData {
				actualDataConverted = append(actualDataConverted, item.(map[string]interface{}))
			}

			// Assert the JSON response
			assert.Equal(t, "success", response["message"])
			assert.Equal(t, float64(200), response["status_code"])
			assert.Equal(t, testCase.expectedData, actualDataConverted)
		})
	}
}
