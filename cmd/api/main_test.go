package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSignupHandlerAcceptsJSON(t *testing.T) {
	mux := newServerMux()

	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(`{"email":"user@example.com","password":"secret"}`))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["email"] != "user@example.com" {
		t.Fatalf("expected email in response, got %q", body["email"])
	}

	if _, ok := body["password_hash"]; ok {
		t.Fatal("response should not include password_hash")
	}

	if _, ok := body["password"]; ok {
		t.Fatal("response should not include password")
	}
}

func TestSignupHandlerReturnsJSONErrorForInvalidJSON(t *testing.T) {
	mux := newServerMux()

	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(`{"email":`))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["error"] != "invalid JSON" {
		t.Fatalf("expected invalid JSON error, got %q", body["error"])
	}
}

func TestSignupHandlerReturnsJSONErrorForEmptyEmail(t *testing.T) {
	mux := newServerMux()

	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(`{"email":"","password":"secret"}`))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["error"] != "email is required" {
		t.Fatalf("expected email required error, got %q", body["error"])
	}
}

func TestSignupHandlerReturnsJSONErrorForEmptyPassword(t *testing.T) {
	mux := newServerMux()

	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(`{"email":"user@example.com","password":""}`))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["error"] != "password is required" {
		t.Fatalf("expected password required error, got %q", body["error"])
	}
}
