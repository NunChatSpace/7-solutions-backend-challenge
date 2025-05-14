package users_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"atomicgo.dev/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database/mongo/repositories"
	adapterHTTP "github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/http"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/savsgio/atreugo/v11"
)

var client *mongo.Client
var db *mongo.Database

func TestMain(m *testing.M) {
	os.Setenv("APP__ENV", "test")
	clientOptions := options.Client().ApplyURI("mongodb://mongou:mongop@localhost:27017/?authSource=admin")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db = client.Database("_test")
	err = db.Drop(context.Background())
	if err != nil {
		log.Fatalf("Failed to drop test database: %v", err)
	}
	err = db.CreateCollection(context.Background(), "users")
	if err != nil {
		log.Fatalf("Failed to create test collection: %v", err)
	}

	m.Run()

	err = client.Disconnect(context.Background())
	if err != nil {
		log.Fatalf("Failed to disconnect MongoDB: %v", err)
	}

}

func TestUserAPIs(t *testing.T) {
	go func() {
		cfg, err := config.LoadConfig()
		if err != nil {
			panic(err)
		}

		config := atreugo.Config{
			Addr:             cfg.App.Port,
			GracefulShutdown: true,
		}

		deps := di.NewDependency(cfg)
		repositories.ProvideRepositories(deps)
		services.ProvideServices(deps)
		server := adapterHTTP.NewServer(deps, config)
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	time.Sleep(1 * time.Second)
	// Test: POST /users
	t.Run("POST /users", func(t *testing.T) {
		jsonBody := []byte(`
			{
				"name": "John Doe",
				"email": "john@doe.com",
				"password": "password123"
			}
		`)

		resp, err := http.Post("http://localhost:8888/api/v1/users", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		if resp.StatusCode > 299 || resp.StatusCode < 200 {
			t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
		}

		defer resp.Body.Close()
		if _, err := io.ReadAll(resp.Body); err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
	})

	userID := ""
	t.Run("GET /users", func(t *testing.T) {
		loginBody := []byte(`
			{
				"email": "john@doe.com",
				"password": "password123"
			}
		`)
		resp1, err := http.Post("http://localhost:8888/api/v1/sessions", "application/json", bytes.NewBuffer(loginBody))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp1.Body.Close()
		resp1body := make(map[string]interface{})
		if _resp1body, err := io.ReadAll(resp1.Body); err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		} else {
			if err := json.Unmarshal(_resp1body, &resp1body); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}
		}
		if resp1.StatusCode > 299 || resp1.StatusCode < 200 {
			t.Fatalf("Expected status code 200, got %d", resp1.StatusCode)
		}

		req, err := http.NewRequest("GET", "http://localhost:8888/api/v1/users", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resp1body["access_token"]))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		var users []map[string]interface{}
		err = json.Unmarshal(body, &users)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(users) == 0 {
			t.Fatalf("No user returned in response")
		}

		userID = users[0]["id"].(string)
		if resp.StatusCode > 299 || resp.StatusCode < 200 {
			t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
		}

	})

	t.Run("GET /users/:id", func(t *testing.T) {
		loginBody := []byte(`
			{
				"email": "john@doe.com",
				"password": "password123"
			}
		`)
		resp1, err := http.Post("http://localhost:8888/api/v1/sessions", "application/json", bytes.NewBuffer(loginBody))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp1.Body.Close()
		resp1body := make(map[string]interface{})
		if _resp1body, err := io.ReadAll(resp1.Body); err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		} else {
			if err := json.Unmarshal(_resp1body, &resp1body); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}
		}
		if resp1.StatusCode > 299 || resp1.StatusCode < 200 {
			t.Fatalf("Expected status code 200, got %d", resp1.StatusCode)
		}

		url := fmt.Sprintf("http://localhost:8888/api/v1/users/%s", userID)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resp1body["access_token"]))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		if resp.StatusCode > 299 || resp.StatusCode < 200 {
			t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PATCH /users/:id", func(t *testing.T) {
		loginBody := []byte(`
			{
				"email": "john@doe.com",
				"password": "password123"
			}
		`)
		resp1, err := http.Post("http://localhost:8888/api/v1/sessions", "application/json", bytes.NewBuffer(loginBody))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp1.Body.Close()
		resp1body := make(map[string]interface{})
		if _resp1body, err := io.ReadAll(resp1.Body); err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		} else {
			if err := json.Unmarshal(_resp1body, &resp1body); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}
		}
		if resp1.StatusCode > 299 || resp1.StatusCode < 200 {
			t.Fatalf("Expected status code 200, got %d", resp1.StatusCode)
		}
		jsonBody := []byte(`
			{
				"email": "Doe@John.com"
			}`)
		url := fmt.Sprintf("http://localhost:8888/api/v1/users/%s", userID)
		req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resp1body["access_token"]))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		var user map[string]interface{}
		err = json.Unmarshal(body, &user)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if resp.StatusCode > 299 || resp.StatusCode < 200 {
			t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
		}

		// Assert response code
		assert.Equal(user["email"], "Doe@John.com")
	})

	t.Run("DELETE /users/:id", func(t *testing.T) {
		loginBody := []byte(`
			{
				"email": "Doe@John.com",
				"password": "password123"
			}
		`)
		resp1, err := http.Post("http://localhost:8888/api/v1/sessions", "application/json", bytes.NewBuffer(loginBody))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp1.Body.Close()
		resp1body := make(map[string]interface{})
		if _resp1body, err := io.ReadAll(resp1.Body); err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		} else {
			if err := json.Unmarshal(_resp1body, &resp1body); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}
		}
		if resp1.StatusCode > 299 || resp1.StatusCode < 200 {
			t.Fatalf("Expected status code 200, got %d", resp1.StatusCode)
		}

		url := fmt.Sprintf("http://localhost:8888/api/v1/users/%s", userID)
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resp1body["access_token"]))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode > 299 || resp.StatusCode < 200 {
			t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
		}
	})

	t.Run("TestCreateAndLoginAndAuthenticateUser", func(t *testing.T) {
		jsonBody := []byte(`
			{
				"name": "John Doe",
				"email": "john@doe123.com",
				"password": "password123"
			}
		`)

		resp1, err := http.Post("http://localhost:8888/api/v1/users", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp1.Body.Close()
		if _, err := io.ReadAll(resp1.Body); err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		if resp1.StatusCode > 299 || resp1.StatusCode < 200 {
			t.Fatalf("Expected status code 200, got %d", resp1.StatusCode)
		}
		loginBody := []byte(`
			{
				"email": "john@doe123.com",
				"password": "password123"
			}
		`)

		resp2, err := http.Post("http://localhost:8888/api/v1/sessions", "application/json", bytes.NewBuffer(loginBody))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp2.Body.Close()
		resp2body := make(map[string]interface{})
		if _resp2body, err := io.ReadAll(resp2.Body); err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		} else {
			if err := json.Unmarshal(_resp2body, &resp2body); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}
		}
		if resp2.StatusCode > 299 || resp2.StatusCode < 200 {
			t.Fatalf("Expected status code 200, got %d", resp2.StatusCode)
		}

		req, err := http.NewRequest("GET", "http://localhost:8888/api/v1/users", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resp2body["access_token"]))
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if _, err := io.ReadAll(resp.Body); err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
	})
}
