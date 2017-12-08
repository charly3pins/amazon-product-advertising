package amazon

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_ok(t *testing.T) {
	criteria := Criteria{}
	expectedResponse := ItemSearchResponse{}

	server := httptest.NewServer(&myHandler{func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Error("unexpected method:", req.Method)
			return
		}
		w.WriteHeader(http.StatusOK)
	}})
	defer server.Close()

	client := NewClient(ClientConfig{AWSEndpoint: server.URL})
	res, err := client.ItemSearch(criteria)
	if err != nil {
		t.Error(err)
		return
	}
	if res.Items.TotalPages != expectedResponse.Items.TotalPages {
		t.Error("Response unexpected")
		return
	}
}

func TestNewClient_ko(t *testing.T) {
	criteria := Criteria{}

	server := httptest.NewServer(&myHandler{func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Error("unexpected method:", req.Method)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}})
	defer server.Close()

	client := NewClient(ClientConfig{AWSEndpoint: server.URL})
	_, err := client.ItemSearch(criteria)
	if err != ErrBadStatusCode {
		t.Error("unexpected error:", err)
		return
	}
}

type myHandler struct {
	serveHTTP func(http.ResponseWriter, *http.Request)
}

func (h myHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.serveHTTP(w, req)
}
