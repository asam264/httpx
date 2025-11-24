package examples

import (
	"context"
	"github.com/asam264/httpx"
	"testing"
	"time"
)

func TestBasicUsage(t *testing.T) {
	client := httpx.New()

	var result map[string]any
	err := client.GetJSON(context.Background(),
		"https://api.github.com/users/octocat", &result)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Result: %+v", result)
}

func TestWithConfig(t *testing.T) {
	client := httpx.New(
		httpx.WithTimeout(5*time.Second),
		httpx.WithRetry(3),
		httpx.WithMiddleware(httpx.LoggingMiddleware()),
	)

	var result any
	err := client.PostJSON(context.Background(),
		"https://httpbin.org/post",
		map[string]string{"hello": "world"},
		&result)

	if err != nil {
		t.Fatal(err)
	}
}
