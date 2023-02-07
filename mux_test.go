package mux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMux(t *testing.T) {
	m := New()

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

func TestNew(t *testing.T) {
	m := New()
	if m.routes == nil {
		t.Error("Expected routes to be initialized")
	}
}

func TestConnect(t *testing.T) {
	m := New()
	m.Connect("/", func(w http.ResponseWriter, r *http.Request) {})
	if len(m.routes["CONNECT"]) != 1 {
		t.Error("Expected route to be added")
	}
}

func TestDelete(t *testing.T) {
	m := New()
	m.Delete("/", func(w http.ResponseWriter, r *http.Request) {})
	if len(m.routes["DELETE"]) != 1 {
		t.Error("Expected route to be added")
	}
}

func TestGet(t *testing.T) {
	m := New()
	m.Get("/", func(w http.ResponseWriter, r *http.Request) {})
	if len(m.routes["GET"]) != 1 {
		t.Error("Expected route to be added")
	}
}

func TestGetWithParams(t *testing.T) {
	mux := New()
	expectedParamValue := "testvalue"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paramValue := Param(r, "testparam")
		if paramValue != expectedParamValue {
			t.Errorf("expected param value %q but got %q", expectedParamValue, paramValue)
		}
	})

	mux.Get("/test/:testparam", handler)

	req, err := http.NewRequest("GET", "/test/"+expectedParamValue, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHead(t *testing.T) {
	m := New()
	m.Head("/", func(w http.ResponseWriter, r *http.Request) {})
	if len(m.routes["HEAD"]) != 1 {
		t.Error("Expected route to be added")
	}
}

func TestOptions(t *testing.T) {
	m := New()
	m.Options("/", func(w http.ResponseWriter, r *http.Request) {})
	if len(m.routes["OPTIONS"]) != 1 {
		t.Error("Expected route to be added")
	}
}

func TestPatch(t *testing.T) {
	m := New()
	m.Patch("/", func(w http.ResponseWriter, r *http.Request) {})
	if len(m.routes["PATCH"]) != 1 {
		t.Error("Expected route to be added")
	}
}

func TestPost(t *testing.T) {
	m := New()
	m.Post("/", func(w http.ResponseWriter, r *http.Request) {})
	if len(m.routes["POST"]) != 1 {
		t.Error("Expected route to be added")
	}
}

func TestPut(t *testing.T) {
	m := New()
	m.Put("/", func(w http.ResponseWriter, r *http.Request) {})
	if len(m.routes["PUT"]) != 1 {
		t.Error("Expected route to be added")
	}
}

func TestTrace(t *testing.T) {
	m := New()
	m.Trace("/", func(w http.ResponseWriter, r *http.Request) {})
	if len(m.routes["TRACE"]) != 1 {
		t.Error("Expected route to be added")
	}
}

func TestUse(t *testing.T) {
	m := New()
	m.Use(func(next http.Handler) http.Handler {
		return next
	})
	if len(m.middlewares) != 1 {
		t.Error("Expected middleware to be added")
	}
}

func TestNoMatchingRoute(t *testing.T) {
	mux := New()
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	r, _ := http.NewRequest("GET", "/not-found", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func BenchmarkMux(b *testing.B) {
	mux := New()
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		b.Fatalf("Error creating request: %v", err)
	}

	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, r)
	}
}

func BenchmarkWithMiddleware(b *testing.B) {
	mux := New()
	mux.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)
		})
	})
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		b.Fatalf("Error creating request: %v", err)
	}

	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, r)
	}
}

func BenchmarkMatchingRoute(b *testing.B) {
	mux := New()
	mux.Get("/users/:id", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r, err := http.NewRequest("GET", "/users/123", nil)
	if err != nil {
		b.Fatalf("Error creating request: %v", err)
	}

	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, r)
	}
}

func BenchmarkNoMatchingRoute(b *testing.B) {
	mux := New()
	mux.Get("/users/:id", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r, err := http.NewRequest("GET", "/products/123", nil)
	if err != nil {
		b.Fatalf("Error creating request: %v", err)
	}

	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, r)
	}
}
