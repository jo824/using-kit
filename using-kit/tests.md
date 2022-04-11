## Testing Go applications

A major benefit of using Go is the tooling provide to us as part of the `Go` command. We're going to look at the `go test` tool.
Let's start by running `go help test`. This will give us our first look at what we can expect from test functions and how to set them up.

```
usage: go test [build/test flags] [packages] [build/test flags & test binary flags]

'Go test' automates testing the packages named by the import paths.
It prints a summary of the test results in the format:

ok   archive/tar   0.011s
FAIL archive/zip   0.022s
ok   compress/gzip 0.033s
...

followed by detailed output for each failed package.

'Go test' recompiles each package along with any files with names matching
the file pattern "*_test.go".
These additional files can contain test functions, benchmark functions, and
example functions. See 'go help testfunc' for more.

```
Takeaways from the above:

1. A test is created by writing a function with a name beginning with "Test" followed by a function name.
2. The `go test` tool matches the file pattern "*_test.go".
3. Provides a summary output with pass/fail/time to run information.
4. In addition to tests we can create example functions, or use benchmarks to monitor the performance of a function over time.


### Honorable mention: Go Vet
`go vet` performs a static analysis of your code and warns you of things that wouldn't be picked up by the compiler. You can run `go tool vet help` to see the different analyzers that are available to run.


Go standard library provides us with a complete package for testing, named ... `testing`. We'll take a look at using the testing package in [the microservice]() we created in a [previous post]().



## Testing package
The testing package provides the tools we need to write unit tests. Testing code typically lives in the same package as the code it tests. Looking at our example code we have `service.go` and `service_test.go`
Let's look at `service_test.go` and the testing function signatures.

### Service layer tests
```go
//service_test.go

package service

import (
    "testing"
    )

func TestGetAThing(t *testing.T) {
    ...
}

...

```

The only parameter must be `t *Testing.T`, which is a pointer to `type T`. Functions I use from the testing package.
- `t.Fail()` which marks the function as failed but continues on with execution of other tests.
- `t.Log()` works the same as fmt.Printf. The text will be printed only if the test fails or the -test.v flag is set.
- `t.Errorf` is equivalent to Logf followed by Fail.


#TODO where to show a full test example?

### Transport layer.
We can find more tests in our transport layer. Here I like to test the payload bodies, http status codes of each response, and the body of each response.

```go

    func TestHTTP(t *testing.T) {
        svc := NewThingSvc()
        h := BuildHTTPHandler(svc, log.NewNopLogger())
        testSrv := httptest.NewServer(h)
        defer testSrv.Close()

        addReq := postThingRequest{
            ID:        "thing1",
            Available: true,
        }
        body, _ := json.Marshal(addReq)

        for _, tc := range []struct {
            m, url, b string
            expected      int
        }{
            {"GET", testSrv.URL + "/thing/yik", "", http.StatusOK},
            {"GET", testSrv.URL + "/thing/exists", "", http.StatusNotFound},
            {"POST", testSrv.URL + "/thing", string(body), http.StatusOK},
            {"POST", testSrv.URL + "/thing", string(body), http.StatusBadRequest},
            {"POST", testSrv.URL + "/thing", "", http.StatusBadRequest},
        } {
            req, _ := http.NewRequest(tc.m, tc.url, strings.NewReader(tc.b))
            res, _ := http.DefaultClient.Do(req)
            if tc.expected != res.StatusCode {
                t.Errorf("%s %s %s: expected %d have %d", tc.m, tc.url, tc.b, tc.expected, res.StatusCode)
            }
        }
    }
```

In the example just above we see a table-driven style to testing. We create an anonymous struct in our for-loop. Scoped to tc(testCase) inside our loop. This structure is defined with 4 fields.
1. HTTP request method.
2. URL of our test server.
3. Payload body.
4. Expected HTTP status code.

The first 3 fields are to build our HTTP request object. The 4th field is an integer that represents the HTTP status code we expect to be returned from our endpoint. Lets run our tests in the service package and view the results.
We run our tests with the command `GOFLAGS="-count=1" go test 	using-kit/using-kit/service/. -v`.
- GOFLAGS="-count=1" to disable caching of our test results.
- '-v' for verbose output.

```
=== RUN   TestHTTP
--- PASS: TestHTTP (0.00s)
=== RUN   TestGetAThing
=== RUN   TestGetAThing/yik
=== RUN   TestGetAThing/yak
=== RUN   TestGetAThing/nope
--- PASS: TestGetAThing (0.00s)
--- PASS: TestGetAThing/yik (0.00s)
--- PASS: TestGetAThing/yak (0.00s)
--- PASS: TestGetAThing/nope (0.00s)
=== RUN   TestGetAllThings
--- PASS: TestGetAllThings (0.00s)
=== RUN   TestAddThing
--- PASS: TestAddThing (0.00s)
PASS
ok  	using-kit/using-kit/service	0.175s

```

## TODO mention of Go code coverage tooling using ccov.png graphic
## Provide examples function usage in a test + benchmarking a function.
## show a test that fails to demonstrate output and fix.
