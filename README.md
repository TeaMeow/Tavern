# Tavern [![GoDoc](https://godoc.org/github.com/teacat/tavern?status.svg)](https://godoc.org/github.com/teacat/tavern) [![Coverage Status](https://coveralls.io/repos/github/teacat/tavern/badge.svg?branch=master)](https://coveralls.io/github/teacat/tavern?branch=master) [![Build Status](https://travis-ci.org/teacat/tavern.svg?branch=master)](https://travis-ci.org/teacat/tavern) [![Go Report Card](https://goreportcard.com/badge/github.com/teacat/tavern)](https://goreportcard.com/report/github.com/teacat/tavern)

Tavern is a _struct-tag free_ validation library with custom validator supported.

## Why?

Tired of all the form validation that based on the struct tags. With Tavern you can validate the values without a struct.

## Indexes

## Installation

```bash
$ go get github.com/teacat/tavern
```

## Example

```go
err := tavern.Validate([]tavern.Rule{
    {
        Value: "Hello, world!",
        Validators: []tavern.Validator{
            tavern.WithRequired(),
            tavern.WithMinLength(5),
        },
    },
})
if err != nil {
    panic(err)
}
```

## Validators

A validation contains rules, and a rule requires a value and the validators. Validators can be chained and communicates by passing the context.

Here are the few built-in validators as example: `WithRequired`, `WithLength`, `WithRange`, etc. Check [GoDoc](https://pkg.go.dev/github.com/teacat/tavern) to see more built-in validators.

### Custom Validators

It's possible to create your own validators.

```go
type Validator func(ctx context.Context, value interface{}) (context.Context, error)
```

The context argument can be used as a communication between the validators, with `ctx.Value(tavern.KeyRequired).(bool)` you are able to get a boolean that states the value is required or not.

## Custom Errors

By default, Tavern returns simple built-in errors such as `ErrRequired`, `ErrLength` might not be what as you wanted. But you are able to create your own custom error for each validator by using `WithCustomError(validator Validator, err error)` function.

It's also a validator but returns your own custom error when the passed-in validator failed.

```go
err := tavern.Validate([]tavern.Rule{
    {
        Value: "",
        Validators: []tavern.Validator{
            tavern.WithCustomError(tavern.WithRequired(), errors.New("nani the fuck")),
            tavern.WithMinLength(5),
        },
    },
})
if err != nil {
    panic(err) // output: nani the fuck
}
```
