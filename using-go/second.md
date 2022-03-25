# Why you should be writing your microservices in Go
Keep it simple, stupid. The KISS principle is one of my favorites. Often we are guilty of making systems unnecessarily
complex. This creates a miserable cycle of working with these things that we create. Does it have to be this way?
Go was designed with simplicity & ease of use top of mind. I'm claiming by using Go to build your next server/service
the benefits will go beyond the performance of your software. You will develop software faster, be happier, and maybe
even bored with how easy it is to use Go. This will be the first article in a series exploring building systems with Go.
In this tutorial we'll be looking at:
    * Origins of Go
    * Some language features & design principles
    * the `net/http` package from the standard library
    * what features we expect in a production ready microservice
    * a few iterations/different approaches to building a server
        * A standard library approach
        * Swapping out the default router
        * Go kit approach to organizing our microservice

By the end it is my hope that you will have enough information to accurately assess using Go in your own work. Also,
enough knowledge of Go to confidently build your own service from scratch.

//It feels like a dynamically
//Go was designed to help us do just that.
//It was created wi
//Most of us are working with [distributed systems](/#:~:text=Distributed%20System%20-%20Definition,order%20to%20achieve%20common%20goals.).
//Go is designed with that in mind.


## Origins of Go.
Picking the right tool for the job is important. Go is the language of the cloud, that's building the modern cloud
(i.e. docker, kubernetes). Before Go it didn't seem like any one language existed that checked all the boxes
required for current day development challenges. Why can't we have efficient compilation,
efficient execution, and ease of programming? In a talk from 2012 Rob Pike describes why they created Go for use at Google.
>>> Go was designed and developed to make working in this environment more productive. Besides its better-known aspects
such as built-in concurrency and garbage collection, Go's design considerations include rigorous dependency management,
the adaptability of software architecture as systems grow, and robustness across the boundaries between components.

A programming language built to deal with the challenges of software engineering? Specifically working with distributed
systems.

Alright, enough of this and lets actually look at some code.

### Quick look at syntax
Go was designed with syntax simplicity and readability in mind. We wanted a feeling of familiarity to existing languages
e declaration syntax is closer to Pascal's than to C's. The declared name appears before
the type and there are more keywords:
```
    var fn func([]int) int
    type T struct { a, b int }
```
Here we declare a function named 'fn' , and a struct of custom type `T`

Here's the same declarations in C.
```
int (*fn)(int[]);
struct T { int a, b; }
```
Declarations introduced by keyword are easier to parse both for people and for computers, and having the type syntax
not be the expression syntax as it is in C has a significant effect on parsing: it adds grammar but eliminates ambiguity.
But there is a nice side effect, too: for initializing.
Here's an example of explicit vs derived initialization
```
// NewT is a method on type T above that returns a newly create T struct
// a method is just a function with a special receiver, in this case T.
func (t T) NewT(a,b int)(T){
}

var myStruct T = t.newT(val1,val2) //explicit - assuming type T has this method associated with it.
mystruct := t.newT(val1,val2)      //derived
```

To me, Go feels like a dynamically typed language but I'm still getting the benefits and speed of a statically typed,
compiled language.
[More on syntax in Go](https://go.dev/blog/declaration-syntax)


### Interfaces
Interfaces in Go are one of, if no the best feature of the language
Go's interfaces let you use [duck typing](https://en.wikipedia.org/wiki/Duck_typing) like you would in a purely dynamic
language like Python but still have the compiler catch obvious mistakes. Go encourages composition over inheritance,
using simple, often one-method interfaces to define trivial behaviors that serve as clean, comprehensible
boundaries between components.

In the next section we'll see concrete examples of using interfaces and how they relate to building our first server.

[more on composition over inheritance](https://go.dev/talks/2012/splash.article#TOC_15.)
[more on types](https://go.dev/doc/faq#types)

### Building an HTTP server using Go's standard Library.
Go’s true power comes from the fact that the language is small and the standard library is large. It enables newcomers
to ramp up quickly. Why? It limits the number of ways you can do something and has a very opinionated view of the world
at the compiler level.
The goal of this section is to not only build our Go http server, but understand the structs and abstractions involved,
and how they fit together. We will build off these ideas as we evolve our server.

So lets dig into the `net/http` package. The first piece we'll need is the `Handler` interface.

```
type Handler interface {
ServeHTTP(ResponseWriter, *Request)
}
```

The `Handler` interface is contains a single method `ServeHTTP` that takes 2 values"
 1. `ResponseWriter` which is an interface is used by an HTTP handler to construct an HTTP response
 2. A pointer to an `http.Request`

We'll look at the Request and Response objects in more detail in a bit. We'll work each of them often.
Now lets build a handler knowing what we know about Go types and interfaces.

```
type MyFirstHandler string

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
//ListenAndServer func signature -- ListenAndServe(addr string, handler Handler) error
err := http.ListenAndServe(":8833",handler)
if err != nil {
fmt.Println("error while attempting to listen for incoming connections", err)
}
}

//main.go
```

Lets break this down a bit. `ListenAndServer` tells our app to listen on a specific port/network address that we provide.
The second parameter is our handler. If you read the comment above you'll see that http library provides us with a default
handler - `DefaultServeMux` if we do not provide our own. What is a ServeMux? From the standard library comments:

```
// ServeMux is an HTTP request multiplexer.
// It matches the URL of each incoming request against a list of registered
// patterns and calls the handler for the pattern that
// most closely matches the URL.

```
Still ServeMux? Handler?
* If you've ever used an http framework in another language that leveraged an MVC-pattern, handlers are similar to controllers.
They perform your application logic & write response information (headers, body, etc).
* servemux is also referred to as a router. [There are additional routers available to you](https://benhoyt.com/writings/go-routing/)
outside of the standard library. Primary function of a router is to store a mapping between url paths and their handlers
that we define. Router features, and implementation will vary, but they all implement `ServeHTTP` interface.

[benchmarks for routers](https://github.com/julienschmidt/go-http-routing-benchmark)

Knowing what we do now. Lets reshuffle our example a bit.


```
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

//main.go
```
What's changed?

* We've edited the handler to return the current time instead of printing a static string.
* Declaring a serveMux aka router in our main function.
* Register our handler to a specific route "/use-handler"







