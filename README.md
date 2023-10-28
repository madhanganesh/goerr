# goerr
This package defines an error type that maintains a stack of nested errors and gives a human readable stack trace for logging.

# Problem
In typical go `error` type, you don't have the stack trace of complete call chain. The only possibility is to log the stack trace at every function in the call chain. This will have multiple log entries made from single call and also the log are very cumbersome. 
Also the default stack trace contains lot many details that becomes difficult to read.

# Principle
This package is to facilitate the well-known idiom

"throw error multiple times, but log once at the top most level"

For eg., if the call chain is `controller` -> `service` -> `repository`. 

Then with `goerr` you return error from repository to service to controller and log the same in controller.

# Output
If you return `goerr` in all methods and get the stack at top most level, it gives below nicely formatted, easily readable stack

```shell
controller failed [/Users/madhan.ganesh/src/github.com/madhanganesh/goerr/samplesrc/samples.go:11 (samplesrc.Controller)]
    service failed [/Users/madhan.ganesh/src/github.com/madhanganesh/goerr/samplesrc/samples.go:19 (samplesrc.Service)]
        error from database [/Users/madhan.ganesh/src/github.com/madhanganesh/goerr/samplesrc/samples.go:26 (samplesrc.Repository)]
```

# Usage
Whenver you wanted to return an error just use
```go
err := goerr.New(nil, "error in here")
```

if you have to nest the error, just pass it in New
```go
err1 := goerr.New(err, "error in here")
```

# Sample code that log in nested methods
```go
func Controller() error{
	err := Service()
	if err != nil {
		return goerr.New(err, "controller failed")
	}
	return nil
}

func Service() error {
	err := Repository()
	if err != nil {
		return goerr.New(err, "service failed")
	}
	return err
}

func Repository() error {
	err := errors.New("error from database")
	return goerr.New(nil, err.Error())
}
```
# Code to get the stack at the top most level
```go
        err := samplesrc.Controller()
	if err != nil {
		t.Logf("error in controller: %s", goerr.Stack(err))
	}
```

# Compatibility
- `goerr` implements standard `error` interface, so can be assigned where ever error is used
- `goerr.Error()` will give the error text of the top most error object
- `goerr.Stack(err)` will give stack trace of call chain
- `goerr.Stack(err)` can be called for error type as well, in which case it will just return `Error()`

# Installation
```shell
go install github.com/madhanganesh/goerr
```
