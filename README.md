# hyper-mux
A very simplified Mux focused for ease of use.

```shell
go get github.com/hyperjumptech/hyper-mux
```

## How to use the Mux

The following is how you going to use Hyper-Mux with vanilla http.Server

```go
import (
    mux "github.com/hyperjumptech/hyper-mux"
    "net/http"
)

var (
    hmux := mux.NewHyperMux()	
)

...
theServer := &http.Server{
    Addr:     "0.0.0.0:8080",
    Handler:  hmux,
}
err := theServer.ListenAndServe()
if err != nil {
    panic(err.Error())
}
...
```

## Adding Routes

First we create the route's handler function

```go
func HandleHelloHyperMux(w http.ResponseWriter, r *http.Request) {
    hmux.WriteString(w, http.StatusOK, "Hello Hyper-Mux")
}
```

Then we map the `HandleFunc` function to the route

```go
hmux.AddRoute("/", hmux.MethodGet, HandleHelloHyperMux)
```


## Using Middleware

First we create the middleware

```go
func ContextSetter(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        next.ServeHTTP(w,r)
        fmt.Println("POST CALL")
    })
}
```
Then we add it to our middleware chain

```go
hmux.UseMiddleware(ContextSetter)
```