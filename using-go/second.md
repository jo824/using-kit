# Why you should be writing your microservices in Go
Keep it simple, stupid. The KISS principle is one of my favorites. Often we are guilty of making systems unnecessarily  complex. This creates a miserable cycle of working with these things that we create. Does it have to be this way? Go was designed with simplicity and ease of use top of mind. I'm claiming by using Go to build your next server/service the benefits will go beyond the performance of your software. You will develop software faster, be happier, and maybe even bored with how easy it is to use Go.

This will be the first article in a series exploring building systems with Go. In this post we'll be looking at:
- origins of Go
- notes on Go syntax and Interfaces
- the `net/http` package from the standard library
- what features we expect in a production ready microservice #TODO
- a few iterations/different approaches to building a server
- a standard library approach
- swapping out the default serveMux (router)
- Go kit approach to organizing our microservice

By the end it is my hope that you will have enough information to accurately assess using Go in your own work. Also, enough knowledge of Go to confidently build your own service from scratch.

## Origins of Go.
Picking the right tool for the job is important. Go is the language of the cloud, that's building the modern cloud
(i.e. docker, kubernetes). Before Go it didn't seem like any one language existed that checked all the boxes
required for current day development challenges. Why can't we have efficient compilation,
efficient execution, and ease of programming? In [a talk from 2012](https://talks.golang.org/2012/splash.article) Rob Pike describes why they created Go for use at Google.
> Go was designed and developed to make working in this environment more productive. Besides its better-known aspects  such as built-in concurrency and garbage collection, Go's design considerations include rigorous dependency management,  the adaptability of software architecture as systems grow, and robustness across the boundaries between components.

A programming language built to deal with the challenges of software engineering? Specifically working with distributed systems.
Alright, enough of this and lets actually look at some code.

### Quick look at syntax
Go was designed with syntax simplicity and readability in mind. We wanted a feeling of familiarity to existing languages declaration syntax is closer to Pascal's than to C's. The declared name appears before  the type and there are more keywords:
Here we declare a function named 'fn' , and a struct of custom type `T`
``
var fn func([]int) int
type T struct { a, b int }
```

Here's the same declarations in C.
```
int (*fn)(int[]);
struct T { int a, b; }
```
Declarations introduced by keyword are easier to parse both for people and for computers, and having the type syntax not be the expression syntax as it is in C has a significant effect on parsing: it adds grammar but eliminates ambiguity. But there is a nice side effect, too: for initializing.

Here's an example of explicit vs derived initialization
```
// NewT is a method on type T above that returns a newly create T struct
// a method is just a function with a special receiver, in this case type T.
func (t T) NewT(a,b int)(T){
}

var myStruct T = t.newT(val1,val2) //explicit - assuming type T has this method associated with it.This isn't some built in constructor
mystruct := t.newT(val1,val2)      //derived
```

To me, Go feels like a dynamically typed language but I'm still getting the benefits and speed of a statically typed,
compiled language.
[More on syntax in Go](https://go.dev/blog/declaration-syntax)


### Interfaces
Interfaces in Go are one of, if not the best feature of the language. Go's interfaces let you use [duck typing](https://en.wikipedia.org/wiki/Duck_typing) like you would in a purely dynamic
language like Python but still have the compiler catch obvious mistakes. Go encourages composition over inheritance,
using simple, often one-method interfaces to define trivial behaviors that serve as clean, comprehensible
boundaries between components.

In the next section we'll see concrete examples of using interfaces and how they relate to building our first server.

[more on composition over inheritance](https://go.dev/talks/2012/splash.article#TOC_15.)
[more on types](https://go.dev/doc/faq#types)
[interfaces history](https://research.swtch.com/interfaces)

### Building an HTTP server using Go's standard Library.
Go’s true power comes from the fact that the language is small and the standard library is large. It enables newcomers to ramp up quickly.
Why? It limits the number of ways you can do something and has a very opinionated view of the world at the compiler level. The goal of this section is to not only build our Go http server, but understand the structs and abstractions involved, and how they fit together.
We will build off these ideas as we evolve our server. Lets dig into the <code>`net/http`</code> package.

The first piece we'll need is the `Handler` interface.

```
type Handler interface {
ServeHTTP(ResponseWriter, *Request)
}
```

`Handler` is an interface that contains a single method `ServeHTTP`. `ServeHTTP` takes 2 values:
1. `ResponseWriter` which is an interface is used by an HTTP handler to construct an HTTP response.
2. A pointer to an `http.Request` struct.

We'll look at the Request and Response objects as they're represented in Go. With some basics covered, more specifically Interfaces, we are ready to build an
<code>HTTP.Handler</code>.

```go
//main.go

type MyFirstHandler struct{}

func (g MyFirstHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Println("MyFirstHandler type implements http.Handler")
}

func main() {
    var handler MyFirstHandler


    // ListenAndServe listens on the TCP network address addr and then calls
    // Serve with handler to handle requests on incoming connections.
    // Accepted connections are configured to enable TCP keep-alives.
    //
    // The handler is typically nil, in which case the DefaultServeMux is used.
    //
    // ListenAndServe always returns a non-nil error.
    err := http.ListenAndServe(":8833",handler)
    if err != nil {
        fmt.Println("error while attempting to listen for incoming connections", err)
    }
}
```

Lets break this down a bit. `ListenAndServe` tells our app to listen on a specific port/network address that we provide.
The second parameter is our handler. If you read the comment above you'll see that http library provides us with a default  handler - `DefaultServeMux` if we do not provide our own. What is a ServeMux?
From the standard library comments:

```
// ServeMux is an HTTP request multiplexer.
// It matches the URL of each incoming request against a list of registered
// patterns and calls the handler for the pattern that
// most closely matches the URL.
```

* If you've ever used an http framework in another language that leveraged an MVC-pattern, handlers are similar to controllers. They perform your application logic and write response information (headers, body, etc).
* servemux is also referred to as a router. Primary function of a router is to store a mapping between url paths and their handlers that we define. Router features, and implementation will vary, but they all implement `ServeHTTP` interface.

Lets start to improve our server with some small changes. In the code snippet below:
    * We've edited the handler to return the current time instead of printing a static string.
    * Declaring a serveMux aka router in our main function.
    * Register our handler to a specific route `/use-handler`

```go
//main.go
func (g MyFirstHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    tm := time.Now().Format(time.RFC3339)
    w.Write([]byte("The time is: " + tm))
}

func main() {
    var firstHandler MyFirstHandler

    // declare our serveMux in main
    router :=http.NewServeMux()
    // register handler we defined - it now responds to any request to path use-handler
    router.Handle("/use-handler", firstHandler)

    err := http.ListenAndServe(":8833", router)
    if err != nil {
        fmt.Println("error while attempting to listen for incoming connections", err)
    }
}
```

Now lets swap out the default router with a library from outside the standard library. This new router gives us a convenient method for setting named path parameters and helper function for accessing them.
Go's default ServeMux is limited to static routes and does not support parameters in the route pattern.
I may be guilty of deferring to use gorilla out of familiarity, but I like the features it provides and won't optimize this until it becomes an issue for me. If you're curious to learn more about [different routers in Go here's a nice writeup](https://github.com/julienschmidt/go-http-routing-benchmark)
on the state of routers at the time and some benchmark numbers.

Here's what our code looks like now:

```
//main.go

import (
"fmt"
"net/http"
"time"

"github.com/gorilla/mux"
)


type MyFirstHandler struct{}

func (g MyFirstHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    //get URL path param name from the request
    name := mux.Vars(r)["name"]
    tm := time.Now().Format(time.RFC3339)
    w.Write([]byte(fmt.Sprintf("Hello %s, the time is: %s\n", name, tm)))
}
```

For the first time we are importing a package outside the standard library. We're now using the [gorilla/mux](https://github.com/gorilla/mux#readme).
We also use the `mux.Vars`  helper function which takes the route params and puts them in a map for us to access.

Here's the updated route:
```
router.Handle("/use-handler/{name:[a-zA-Z]+}", firstHandler).Methods("GET")

//main.go
```
We set the allowable http verb for this route with `.METHODS("GET")`,set expected route param value, and naming it name.  I think we're ready to turn this into a real service.


## Intro to Go kit
In this section, we'll talk about Go-kit, which I'll refer to as kit for the rest of this post. Quick note: Having a solid understanding of [HTTP](https://developer.mozilla.org/en-US/docs/Web/HTTP/Overview) in general will help with the digestion of this content. At a minimum knowing a little bit about what you can expect
from [request/response objects](https://developer.mozilla.org/en-US/docs/Web/HTTP/Messages).

### What is *kit* and why should I use it?
In their own words:
> kit is a collection of Go (golang) packages (libraries) that help you build robust, reliable, maintainable microservices.
>You should use kit if you know you want to adopt the microservices pattern in your organization. Go kit will help you structure and build out your services, avoid common pitfalls, and write code that grows with grace. Go kit de-risks both Go and microservices by providing mature patterns and idioms, written and maintained by a large group of experienced contributors, and validated in production environments.

We won't dive into all things kit here, but lets take a look at how kit gives us a template for structuring our service logic.

![Path of a request](req-path.png  "Request Path")

After taking a look at this diagram there are some similar pieces to our original simple server implementation. We should be familiar with everything except what is contained inside the purple area(inside handler). It's not that it couldn't have existed, but now we have a clearer picture for organizing our server logic.
The major benefit of *kit* is that it provides some nice abstractions that assist in structuring your service. They group a service into these 3 layers:

1. Transport layer
2. Endpoint layer
3. Service layer

#### Transport layer
Transport layer comes from the [OSI model](https://en.wikipedia.org/wiki/OSI_model#Layer_4:_Transport_layer). The transport layer as it relates to OSI is defined as
``` means of transferring variable-length data sequences from a source to a destination host, while maintaining the quality of service functions. ```
HTTP isn't actually a transport layer, but our HTTP.server relies on TCP which falls in the transport layer.

Here you have the flexibility of implementing one or more transports(example HTTP/gRPC). In our example service we
are using HTTP encoding a json response. The following code is defined in <code>rawkit/server/main.go</code> in our project and we'll look at
that first.

Our example code won't explore the benefits of multiple transports, but work only with HTTP & JSON.


```go
// using-kit/server/main.go

    svc := rawkit.NewThingSvc(logger)

    getThingHandler := httptransport.NewServer(
    loggingMiddleware(log.With(logger, "method", "get-a-thing"))(rawkit.GetAThingEP(svc)),
                    rawkit.DecodeGetThingRequest,
                    httptransport.EncodeJSONResponse,
    )

    r := mux.NewRouter()
    r.Handle("/things/{id:[a-zA-Z]+}", getThingHandler).Methods("GET")

    http.ListenAndServe(DEFAULT_PORT, r)
```

Breaking down our new main function:

1. Again we are using the same [gorilla/mux setup for our router](https://github.com/gorilla/mux).
2. `httptransport.NewServer` is from the package <code>github.com/go-kit/kit/transport/http</code> and it creates a `kit.Server`
                that wraps an endpoint, a decoder for the request. The <code>kit/transport/http.Server</code>code> type implements http.Handler.

Lets peel this back another layer and look at the Endpoint function, decoders, and encoders that live inside this kit defined type, server, that acts as a wrapper.

An `Endpoint` is the fundamental building block of servers and clients. It represents a single RPC method.
Endpoint type is a function that takes in an interface request and returns an interface response. The decoder/encoder
and endpoint func is where your safety and anti-fragile logic will live. Looking back to our simple server example
we interacted with `http.Request` type and we will again here. The first stop on the request path
is our decoder.

```go
//using-kit/service/endpoints.go

func DecodeGetThingRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req GetThingRequest
    req.ID = mux.Vars(r)["id"]
    if len(req.ID) == 0 {
        return nil, errors.New("missing ID route param")
    }
    return req, nil
}
```
What should be inside a decoder function:
- Here we interact with an `http.Request` struct just like our simple server handlers. The goal is to convert this into a different request struct that our
underlying service(s) expect.
    - This abstracts away the http transport bit and is a point where we can switch easily for something else if needed.
- Here we interact with http.Request object. Taking it and converting the request into a struct that our service will interact with going forward.
- We also validate the request. Deserializing the payload body(if it was a post/put request). In this case it's a get request. We grab the route parameter 'id' which is a string, and we confirm that it's not an empty string.

This is a good place for any type of validation and conversion or types we expect to work with.

[kit-faq](https://gokit.io/faq/#what-is-go-kit)
[kit-examples](https://github.com/go-kit/examples)