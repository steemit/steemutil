package jsonrpc2

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// slowHandler sleeps far longer than any test timeout, so the only way a
// request completes quickly is transport-level cancellation/timeout.
func slowHandler(sleep time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(sleep)
		w.WriteHeader(http.StatusOK)
	})
}

func TestSendWithContext_CancelAbortsInFlightRequest(t *testing.T) {
	srv := httptest.NewServer(slowHandler(5 * time.Second))
	defer srv.Close()

	client := NewClient(srv.URL)
	if err := client.BuildSendData("condenser_api.get_block", []any{1}); err != nil {
		t.Fatalf("BuildSendData: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	_, err := client.SendWithContext(ctx)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected error from canceled request, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled in error chain, got: %v", err)
	}
	if elapsed > time.Second {
		t.Errorf("cancel should abort promptly, took %v", elapsed)
	}
}

func TestWithHTTPClient_InjectedTimeoutApplies(t *testing.T) {
	srv := httptest.NewServer(slowHandler(5 * time.Second))
	defer srv.Close()

	injected := &http.Client{Timeout: 100 * time.Millisecond}
	client := NewClientWithOptions(srv.URL, WithHTTPClient(injected))
	if err := client.BuildSendData("condenser_api.get_block", []any{1}); err != nil {
		t.Fatalf("BuildSendData: %v", err)
	}

	start := time.Now()
	_, err := client.Send()
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected timeout error from injected client, got nil")
	}
	if elapsed > time.Second {
		t.Errorf("injected 100ms timeout should apply, took %v", elapsed)
	}
}

func TestWithHTTPClient_DecodesResult(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// A get_block beyond head returns {"result": null} on real nodes.
		_, _ = w.Write([]byte(`{"id":1,"jsonrpc":"2.0","result":null}`))
	}))
	defer srv.Close()

	client := NewClientWithOptions(srv.URL, WithHTTPClient(&http.Client{}))
	if err := client.BuildSendData("condenser_api.get_block", []any{111000000}); err != nil {
		t.Fatalf("BuildSendData: %v", err)
	}

	res, err := client.Send()
	if err != nil {
		t.Fatalf("Send: %v", err)
	}
	if res.Result != nil {
		t.Errorf("expected nil result, got %v", res.Result)
	}
}

// TestSendWithContext_LegacyClientShape pins the fallback behavior: without
// WithHTTPClient the request still succeeds through the per-call 30s client.
func TestSendWithContext_LegacyClientShape(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":1,"jsonrpc":"2.0","result":{"ok":true}}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL)
	if err := client.BuildSendData("condenser_api.get_dynamic_global_properties", []any{}); err != nil {
		t.Fatalf("BuildSendData: %v", err)
	}

	res, err := client.SendWithContext(context.Background())
	if err != nil {
		t.Fatalf("SendWithContext: %v", err)
	}
	if res.Result == nil {
		t.Error("expected non-nil result")
	}
}
