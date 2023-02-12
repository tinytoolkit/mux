package mux_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tinytoolkit/mux"
)

func TestMux(t *testing.T) {
	m := mux.New()

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}

	m.Get("/", handler)

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create a new HTTP request: %v", err)
	}

	w := httptest.NewRecorder()

	m.ServeHTTP(w, r)

	if got, want := w.Code, http.StatusOK; got != want {
		t.Errorf("Unexpected status code: got %d, want %d", got, want)
	}

	if got, want := strings.TrimSpace(w.Body.String()), "Hello, World!"; got != want {
		t.Errorf("Unexpected body: got %q, want %q", got, want)
	}
}

func TestGetWithParams(t *testing.T) {
	m := mux.New()
	expectedParamValue := "testvalue"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paramValue := mux.Param(r, "testparam")
		if paramValue != expectedParamValue {
			t.Errorf("expected param value %q but got %q", expectedParamValue, paramValue)
		}
	})

	m.Get("/test/:testparam", handler)

	req, err := http.NewRequest("GET", "/test/"+expectedParamValue, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	m.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestNoMatchingRoute(t *testing.T) {
	m := mux.New()
	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	r, _ := http.NewRequest("GET", "/not-found", nil)
	w := httptest.NewRecorder()

	m.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func BenchmarkMux(b *testing.B) {
	m := mux.New()
	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		b.Fatalf("Error creating request: %v", err)
	}

	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.ServeHTTP(w, r)
	}
}

func BenchmarkWithMiddleware(b *testing.B) {
	m := mux.New()
	m.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)
		})
	})
	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		b.Fatalf("Error creating request: %v", err)
	}

	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.ServeHTTP(w, r)
	}
}

func BenchmarkMatchingRoute(b *testing.B) {
	m := mux.New()
	m.Get("/users/:id", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r, err := http.NewRequest("GET", "/users/123", nil)
	if err != nil {
		b.Fatalf("Error creating request: %v", err)
	}

	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.ServeHTTP(w, r)
	}
}

func BenchmarkNoMatchingRoute(b *testing.B) {
	m := mux.New()
	m.Get("/users/:id", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r, err := http.NewRequest("GET", "/products/123", nil)
	if err != nil {
		b.Fatalf("Error creating request: %v", err)
	}

	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.ServeHTTP(w, r)
	}
}
