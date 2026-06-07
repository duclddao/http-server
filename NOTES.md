## Middleware

Example
```go
func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
```

Then
```go
mux.Handle("/app/", middlewareLog(handler))
```

## Context
The context package in Go's standard library is used to pass request-scoped information through your program.

In HTTP servers, the most important parts are cancellation and timeouts. When a client disconnects, a request times out, or the server shuts down, the request's context can tell the rest of your code to stop working on that request.

### Request Context
Every http.Request has a context:

```go
ctx := r.Context()
```

That context belongs to the current HTTP request. If the request is canceled, ctx is canceled too.

###Database Calls
Many database APIs accept a context.Context as their first argument. SQLC-generated methods are no exception:

```go
user, err := cfg.db.CreateUser(ctx, params.Email)
```

By passing the request context to the database call, the database work is tied to the lifetime of the HTTP request. If the client gives up before the query finishes, the query can be canceled instead of wasting server resources.

In a handler, this usually means passing r.Context() directly:

user, err := cfg.db.CreateUser(r.Context(), params.Email)

###Background Context
You'll also see [context.Background()](https://pkg.go.dev/context#Background) in Go code. It's useful when a Context is expected but there's no incoming request or parent operation to start from – like in startup code or a background job.

For web handlers, prefer r.Context(). It carries the cancellation signal for the specific request you're handling.