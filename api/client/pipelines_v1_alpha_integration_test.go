package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	models "github.com/semaphoreci/cli/api/models"
)

func makePipelines(ids ...string) []models.PplListElemV1Alpha {
	var pipelines []models.PplListElemV1Alpha
	for _, id := range ids {
		pipelines = append(pipelines, models.PplListElemV1Alpha{
			Id:    id,
			Name:  "pipeline-" + id,
			State: "DONE",
			Label: "main",
		})
	}
	return pipelines
}

// setupTestServer creates a TLS test server and configures http.DefaultTransport
// to trust it. Returns the server and a cleanup function.
func setupTestServer(t *testing.T, handler http.Handler) (*httptest.Server, func()) {
	t.Helper()
	server := httptest.NewTLSServer(handler)
	originalTransport := http.DefaultTransport
	http.DefaultTransport = server.Client().Transport
	cleanup := func() {
		server.Close()
		http.DefaultTransport = originalTransport
	}
	return server, cleanup
}

func newTestPipelinesAPI(server *httptest.Server) PipelinesApiV1AlphaApi {
	host := strings.TrimPrefix(server.URL, "https://")
	return PipelinesApiV1AlphaApi{
		BaseClient:           NewBaseClient("test-token", host, "v1alpha"),
		ResourceNamePlural:   "pipelines",
		ResourceNameSingular: "pipeline",
	}
}

func TestListPplWithOptions_SinglePage(t *testing.T) {
	pipelines := makePipelines("p1", "p2", "p3")
	server, cleanup := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := json.Marshal(pipelines)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}))
	defer cleanup()

	api := newTestPipelinesAPI(server)
	result, err := api.ListPplWithOptions("proj-1", ListOptions{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got models.PipelinesListV1Alpha
	if err := json.Unmarshal(result, &got); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 pipelines, got %d", len(got))
	}
	if got[0].Id != "p1" {
		t.Errorf("expected first pipeline id 'p1', got '%s'", got[0].Id)
	}
}

func TestListPplWithOptions_MultiplePages(t *testing.T) {
	page1 := makePipelines("p1", "p2")
	page2 := makePipelines("p3", "p4")
	page3 := makePipelines("p5")

	server, cleanup := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		w.Header().Set("Content-Type", "application/json")

		switch page {
		case "1", "":
			w.Header().Set("Link", `<http://example.com/api/v1alpha/pipelines?page=2>; rel="next"`)
			data, _ := json.Marshal(page1)
			w.Write(data)
		case "2":
			w.Header().Set("Link", `<http://example.com/api/v1alpha/pipelines?page=3>; rel="next"`)
			data, _ := json.Marshal(page2)
			w.Write(data)
		case "3":
			data, _ := json.Marshal(page3)
			w.Write(data)
		default:
			t.Errorf("unexpected page requested: %s", page)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer cleanup()

	api := newTestPipelinesAPI(server)
	result, err := api.ListPplWithOptions("proj-1", ListOptions{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got models.PipelinesListV1Alpha
	if err := json.Unmarshal(result, &got); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}
	if len(got) != 5 {
		t.Errorf("expected 5 pipelines, got %d", len(got))
	}
	expectedIDs := []string{"p1", "p2", "p3", "p4", "p5"}
	for i, expected := range expectedIDs {
		if got[i].Id != expected {
			t.Errorf("pipeline[%d]: expected id '%s', got '%s'", i, expected, got[i].Id)
		}
	}
}

func TestListPplWithOptions_EmptyFirstPage(t *testing.T) {
	server, cleanup := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}))
	defer cleanup()

	api := newTestPipelinesAPI(server)
	result, err := api.ListPplWithOptions("proj-1", ListOptions{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got models.PipelinesListV1Alpha
	if err := json.Unmarshal(result, &got); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 pipelines, got %d", len(got))
	}
}

func TestListPplWithOptions_QueryParams(t *testing.T) {
	var capturedQueries []string

	server, cleanup := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQueries = append(capturedQueries, r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}))
	defer cleanup()

	api := newTestPipelinesAPI(server)
	_, err := api.ListPplWithOptions("proj-1", ListOptions{
		CreatedAfter:  1000,
		CreatedBefore: 2000,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(capturedQueries) != 1 {
		t.Fatalf("expected 1 request, got %d", len(capturedQueries))
	}

	q := capturedQueries[0]
	for _, expected := range []string{"project_id=proj-1", "created_after=1000", "created_before=2000", "page_size=200", "page=1"} {
		if !strings.Contains(q, expected) {
			t.Errorf("query %q missing expected param %q", q, expected)
		}
	}
}

func TestListPplWithOptions_HTTPError_NonRetryable(t *testing.T) {
	callCount := 0
	server, cleanup := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("access denied"))
	}))
	defer cleanup()

	api := newTestPipelinesAPI(server)
	_, err := api.ListPplWithOptions("proj-1", ListOptions{})

	if err == nil {
		t.Fatal("expected error but got none")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("expected error to contain status code 403, got: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call (no retries for 403), got %d", callCount)
	}
}

func TestListPplWithOptions_ServerError_Retried(t *testing.T) {
	callCount := 0
	server, cleanup := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(makePipelines("p1"))
		w.Write(data)
	}))
	defer cleanup()

	api := newTestPipelinesAPI(server)
	result, err := api.ListPplWithOptions("proj-1", ListOptions{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 3 {
		t.Errorf("expected 3 calls (2 failures + 1 success), got %d", callCount)
	}

	var got models.PipelinesListV1Alpha
	if err := json.Unmarshal(result, &got); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 pipeline, got %d", len(got))
	}
}

func TestListPplWithOptions_DeserializationError_NonRetryable(t *testing.T) {
	callCount := 0
	server, cleanup := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("not valid json"))
	}))
	defer cleanup()

	api := newTestPipelinesAPI(server)
	_, err := api.ListPplWithOptions("proj-1", ListOptions{})

	if err == nil {
		t.Fatal("expected error but got none")
	}
	if !strings.Contains(err.Error(), "deserialize") {
		t.Errorf("expected deserialization error, got: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call (no retries for deserialization), got %d", callCount)
	}
}

func TestListPplWithOptions_RetryWithStaleHeaders(t *testing.T) {
	page1 := makePipelines("p1")
	page2 := makePipelines("p2")
	page2Attempt := 0

	server, cleanup := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		w.Header().Set("Content-Type", "application/json")

		switch page {
		case "1", "":
			w.Header().Set("Link", `<http://example.com?page=2>; rel="next"`)
			data, _ := json.Marshal(page1)
			w.Write(data)
		case "2":
			page2Attempt++
			if page2Attempt == 1 {
				w.Header().Set("Link", `<http://example.com?page=3>; rel="next"`)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("temporary failure"))
				return
			}
			data, _ := json.Marshal(page2)
			w.Write(data)
		case "3":
			t.Error("should not request page 3 — stale headers from failed attempt should be cleared")
			w.WriteHeader(http.StatusBadRequest)
		default:
			t.Errorf("unexpected page: %s", page)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer cleanup()

	api := newTestPipelinesAPI(server)
	result, err := api.ListPplWithOptions("proj-1", ListOptions{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got models.PipelinesListV1Alpha
	if err := json.Unmarshal(result, &got); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 pipelines (p1 + p2), got %d", len(got))
	}
}

func TestListPplWithOptions_PageIncrements(t *testing.T) {
	var requestedPages []string

	totalPages := 3
	server, cleanup := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		requestedPages = append(requestedPages, page)
		w.Header().Set("Content-Type", "application/json")

		pageNum := len(requestedPages)
		if pageNum < totalPages {
			w.Header().Set("Link", fmt.Sprintf(`<http://example.com?page=%d>; rel="next"`, pageNum+1))
		}
		data, _ := json.Marshal(makePipelines(fmt.Sprintf("p%s", page)))
		w.Write(data)
	}))
	defer cleanup()

	api := newTestPipelinesAPI(server)
	_, err := api.ListPplWithOptions("proj-1", ListOptions{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"1", "2", "3"}
	if len(requestedPages) != len(expected) {
		t.Fatalf("expected %d pages requested, got %d", len(expected), len(requestedPages))
	}
	for i, exp := range expected {
		if requestedPages[i] != exp {
			t.Errorf("request %d: expected page=%s, got page=%s", i, exp, requestedPages[i])
		}
	}
}
