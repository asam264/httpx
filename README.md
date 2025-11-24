# httpx

A simple, elegant, and powerful HTTP client library for Go with built-in retry, middleware support, and JSON handling.

## Features

- 🚀 **Simple API**: Fluent, chainable request builder
- 🔄 **Automatic Retry**: Configurable retry with exponential backoff and jitter
- 🔌 **Middleware Support**: Extensible middleware system for logging, metrics, and more
- 📦 **JSON First**: Built-in JSON request/response handling
- ⏱️ **Timeout Control**: Request-level and client-level timeout support
- 🎯 **Error Handling**: Structured HTTP error types with status code checking
- 🔧 **Flexible Configuration**: Rich options for customization

## Installation

```bash
go get github.com/asam264/httpx
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "github.com/asam264/httpx"
)

func main() {
    client := httpx.New()
    
    var result map[string]interface{}
    err := client.GetJSON(context.Background(),
        "https://api.github.com/users/octocat", &result)
    
    if err != nil {
        panic(err)
    }
}
```

### Using Request Builder

```go
client := httpx.New()

var response map[string]interface{}
err := client.NewRequest().
    Get("https://api.example.com/users").
    Query("page", "1").
    Query("limit", "10").
    Header("Authorization", "Bearer token").
    Do(context.Background()).
    Into(&response)
```

### POST Request with JSON

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

client := httpx.New()

var created User
err := client.PostJSON(context.Background(),
    "https://api.example.com/users",
    User{Name: "John", Email: "john@example.com"},
    &created)
```

## Configuration

### With Options

```go
client := httpx.New(
    httpx.WithTimeout(5*time.Second),
    httpx.WithRetry(3),
    httpx.WithHeader("User-Agent", "my-app/1.0"),
    httpx.WithMiddleware(httpx.LoggingMiddleware()),
)
```

### With Base URL

```go
client := httpx.New().
    WithBaseURL("https://api.example.com/v1")

// All requests will use this base URL
client.GetJSON(ctx, "/users", &result)
```

### With Retry Configuration

```go
client := httpx.New().
    WithRetry(3).
    WithRetryBackoff(100*time.Millisecond, 5*time.Second).
    WithRetryIf(func(resp *http.Response, err error) bool {
        // Custom retry logic
        return err != nil || (resp != nil && resp.StatusCode >= 500)
)
```

## Middleware

### Built-in Middleware

```go
// Logging middleware
client := httpx.New(
    httpx.WithMiddleware(httpx.LoggingMiddleware()),
)

// Metrics middleware
client := httpx.New(
    httpx.WithMiddleware(httpx.MetricsMiddleware("my-service")),
)

// Timeout middleware
client := httpx.New(
    httpx.WithMiddleware(httpx.TimeoutMiddleware(2*time.Second)),
)
```

### Custom Middleware

```go
func CustomMiddleware() httpx.Middleware {
    return func(next http.RoundTripper) http.RoundTripper {
        return httpx.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
            // Before request
            req.Header.Set("X-Custom-Header", "value")
            
            // Execute request
            resp, err := next.RoundTrip(req)
            
            // After request
            // Do something with response
            
            return resp, err
        })
    }
}
```

## Response Handling

### JSON Response

```go
var data MyStruct
err := client.NewRequest().
    Get("https://api.example.com/data").
    Do(ctx).
    Into(&data)
```

### Raw Response

```go
resp, err := client.NewRequest().
    Get("https://api.example.com/data").
    Do(ctx).
    Raw()

if err != nil {
    return err
}
defer resp.Body.Close()

// Handle raw response
```

### Bytes/String Response

```go
// Get as bytes
bytes, err := client.NewRequest().
    Get("https://api.example.com/data").
    Do(ctx).
    Bytes()

// Get as string
text, err := client.NewRequest().
    Get("https://api.example.com/data").
    Do(ctx).
    String()
```

## Error Handling

```go
err := client.GetJSON(ctx, url, &result)
if err != nil {
    // Check if it's an HTTP error
    if httpx.IsHTTPError(err) {
        httpErr, _ := httpx.GetHTTPError(err)
        fmt.Printf("Status: %d, Body: %s\n", httpErr.StatusCode, httpErr.Body)
    }
    
    // Check for specific status code
    if httpx.IsStatusCode(err, 404) {
        fmt.Println("Resource not found")
    }
    
    // Check for timeout
    if httpx.IsTimeout(err) {
        fmt.Println("Request timeout")
    }
}
```

## Global Client

```go
// Use global default client
httpx.GetJSON(ctx, "https://api.example.com/data", &result)
httpx.PostJSON(ctx, "https://api.example.com/data", reqBody, &result)

// Set custom default client
httpx.SetDefault(httpx.New(httpx.WithTimeout(10*time.Second)))
```

## Advanced Examples

### Complex Request

```go
client := httpx.New().
    WithBaseURL("https://api.example.com").
    WithHeader("Authorization", "Bearer token")

var result Response
err := client.NewRequest().
    Post("/users").
    Header("Content-Type", "application/json").
    Query("validate", "true").
    JSONBody(User{Name: "John"}).
    Do(ctx).
    Into(&result)
```

### With Custom Transport

```go
transport := &http.Transport{
    MaxIdleConns:        100,
    IdleConnTimeout:     90 * time.Second,
    TLSHandshakeTimeout: 10 * time.Second,
}

client := httpx.New(httpx.WithTransport(transport))
```

## License

Apache 2.0

