# web
implements a Go web server that uses an object to carry request-scoped state

# premise
- context.Context is ugly; minimize its surface area
- objects carry _request-scoped_ state when given proper method recievers
- this applies to an `http.ServeMux`

# method
- Create an object
- Implement an initial handler method as a value reciever
- Assign request scoped data to the value
- Service the request, invoking value or pointer reciever methods from the initial handler if necessary

# references
- golang/go#225602: [relax recommendation against putting Contexts in structs](https://github.com/golang/go/issues/22602)
