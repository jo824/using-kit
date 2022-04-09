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
```
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
We can find more tests in our transport layer. 

