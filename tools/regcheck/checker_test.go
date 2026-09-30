package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestVerifyPipDependencies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/pypi/dploot/json" {
			http.NotFound(w, request)
			return
		}
		_, _ = fmt.Fprint(w, `{"releases":{"3.1.3":{}}}`)
	}))
	defer server.Close()

	checker := newChecker(time.Second, "")
	checker.pypiBase = server.URL
	if err := checker.verifyPipDependencies([]string{"dploot==3.1.3"}); err != nil {
		t.Fatalf("verifyPipDependencies: %v", err)
	}
	if err := checker.verifyPipDependencies([]string{"dploot==4.1.1"}); err == nil {
		t.Fatal("expected unpublished dependency version to fail")
	}
}
