package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func TestInitRedis(t *testing.T) {
	redisdb := redis.NewClient(&redis.Options{
		Addr:         ":6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	err := redisdb.Set("TestKey", "TestValue", 0).Err()
	if err != nil {
		t.Errorf("Redis Set Not working! Check DB connection")
	}

	TestValue, err := redisdb.Get("TestKey").Result()
	if err == redis.Nil {
		t.Errorf("Test Key Not exists!")

	} else if err != nil {
		t.Errorf("Redis Get Not working! Check DB connection")
	} else if TestValue != "TestValue" {
		t.Errorf("Test Value is not correct!Test Value: %v but it should be TestValue", TestValue)
	}
}

func TestRedirectEndpointHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.

	req, err := http.NewRequest("GET", "/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RedirectEndpoint)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	if strings.Contains(rr.Body.String(), "short_link") {
		t.Errorf("handler returned unexpected body, got: %v",
			rr.Body.String())
	}
}

func TestCreateShortLinkEndpointHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.

	var jsonStr = []byte(`{"url":"https://www.google.com"}`)

	req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateShortLinkEndpoint)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	if strings.Contains(rr.Body.String(), "short_link") {
		t.Errorf("handler returned unexpected body, got: %v",
			rr.Body.String())
	}
}
