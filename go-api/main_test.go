package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateAndFollow(t *testing.T) {
	srv := httptest.NewServer(NewStore().Handler())
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/links", "application/json",
		strings.NewReader(`{"url":"https://example.com"}`))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("want 201, got %d", resp.StatusCode)
	}
	var l Link
	if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
		t.Fatal(err)
	}

	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	r2, err := client.Get(srv.URL + "/" + l.Code)
	if err != nil {
		t.Fatal(err)
	}
	if r2.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("want 307, got %d", r2.StatusCode)
	}

	r3, _ := http.Get(srv.URL + "/links/" + l.Code + "/stats")
	var got Link
	_ = json.NewDecoder(r3.Body).Decode(&got)
	if got.Clicks != 1 {
		t.Fatalf("want 1 click, got %d", got.Clicks)
	}
}

func TestUnknownCode(t *testing.T) {
	srv := httptest.NewServer(NewStore().Handler())
	defer srv.Close()
	resp, _ := http.Get(srv.URL + "/nope")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("want 404, got %d", resp.StatusCode)
	}
}
