package http

import (
	"net/http"
	"testing"
)

func TestPagingParamsFromRequest(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://localhost:8080/users?limit=20", nil)
	if err != nil {
		t.Fatalf("error creating request for test: %v", err)
	}
	params, err := PagingParamsFromRequest(r)
	if err != nil {
		t.Fatalf("error parsing valid query params: %v", err)
	}
	if params.Limit != 20 {
		t.Fatalf("expected limit to be 20 but it is %d", params.Limit)
	}
	if params.AfterId != 0 {
		t.Fatalf("expected after_id to be 0 but it is %d", params.AfterId)
	}
	r, err = http.NewRequest(http.MethodGet, "http://localhost:8080/users?limit=20&after_id=20", nil)
	if err != nil {
		t.Fatalf("error creating request for test: %v", err)
	}
	params, err = PagingParamsFromRequest(r)
	if err != nil {
		t.Fatalf("error parsing valid query params: %v", err)
	}
	if params.Limit != 20 {
		t.Fatalf("expected limit to be 20 but it is %d", params.Limit)
	}
	if params.AfterId != 20 {
		t.Fatalf("expected after_id to be 0 but it is %d", params.AfterId)
	}
	r, err = http.NewRequest(http.MethodGet, "http://localhost:8080/users", nil)
	if err != nil {
		t.Fatalf("error creating request for test: %v", err)
	}
	params, err = PagingParamsFromRequest(r)
	if err != nil {
		t.Fatalf("error parsing valid query params: %v", err)
	}
	if params.Limit != 0 {
		t.Fatalf("expected limit to be 0 but it is %d", params.Limit)
	}
	if params.AfterId != 0 {
		t.Fatalf("expected after_id to be 0 but it is %d", params.AfterId)
	}
}
