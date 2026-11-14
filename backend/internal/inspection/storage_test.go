package inspection

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAzurePrivateStorage(t *testing.T) {
	var saved string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/container/image.png" || r.URL.Query().Get("sig") != "secret" {
			t.Errorf("incorrect storage request")
		}
		switch r.Method {
		case http.MethodPut:
			if r.Header.Get("If-None-Match") != "*" || r.Header.Get("x-ms-blob-type") != "BlockBlob" {
				t.Error("missing protected upload headers")
			}
			data, err := io.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
			}
			saved = string(data)
			w.WriteHeader(http.StatusCreated)
		case http.MethodGet:
			_, _ = io.WriteString(w, saved)
		case http.MethodDelete:
			saved = ""
			w.WriteHeader(http.StatusAccepted)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	t.Cleanup(server.Close)
	storage, err := NewStorage("", server.URL+"/container?sig=secret")
	if err != nil {
		t.Fatal(err)
	}
	storage.client = server.Client()
	if err := storage.Put(context.Background(), "image.png", "image/png", []byte("image")); err != nil {
		t.Fatal(err)
	}
	data, err := storage.Get(context.Background(), "image.png")
	if err != nil || string(data) != "image" {
		t.Fatalf("get: %q %v", data, err)
	}
	if err := storage.Delete(context.Background(), "image.png"); err != nil || saved != "" {
		t.Fatalf("delete: %v", err)
	}
	if _, err := NewStorage("", "http://localhost/container?sig=secret"); err == nil {
		t.Fatal("insecure storage URL accepted")
	}
}

func TestMLUnavailableAndTimeout(t *testing.T) {
	client, err := NewMLClient("http://127.0.0.1:1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Predict(context.Background(), "tomato", "image/png", []byte("image"))
	var mlErr *MLError
	if !errors.As(err, &mlErr) || mlErr.Code != "ML_UNAVAILABLE" {
		t.Fatalf("unavailable: %v", err)
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	cancel()
	_, err = client.Predict(ctx, "tomato", "image/png", []byte("image"))
	if !errors.As(err, &mlErr) || mlErr.Code != "ML_TIMEOUT" || strings.Contains(err.Error(), "127.0.0.1") {
		t.Fatalf("context error: %v", err)
	}
}
