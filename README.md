# goamz - An Amazon Library for Go

[![Build Status](http://travis-ci.org/czos/goamz.png?branch=master)](https://travis-ci.org/czos/goamz)

The _goamz_ package enables Go programs to interact with Amazon Web Services.

This is a fork of the version [developed within Canonical](https://wiki.ubuntu.com/goamz) with additional functionality and services from [a number of contributors](https://github.com/goamz/goamz/contributors)!

The API of AWS is very comprehensive, though, and goamz doesn't even scratch the surface of it. That said, it's fairly well tested, and is the foundation in which further calls can easily be integrated. We'll continue extending the API as necessary - Pull Requests are _very_ welcome!

The following packages are available at the moment:

```
github.com/czos/goamz/autoscaling
github.com/czos/goamz/aws
github.com/czos/goamz/cloudfront
github.com/czos/goamz/cloudwatch
github.com/czos/goamz/dynamodb
github.com/czos/goamz/ec2
github.com/czos/goamz/elb
github.com/czos/goamz/iam
github.com/czos/goamz/rds
github.com/czos/goamz/route53
github.com/czos/goamz/s3
github.com/czos/goamz/sqs

github.com/czos/goamz/exp/mturk
github.com/czos/goamz/exp/sdb
github.com/czos/goamz/exp/sns
```

Packages under `exp/` are still in an experimental or unfinished/unpolished state.

## API documentation

The API documentation is currently available at:

[http://godoc.org/github.com/czos/goamz](http://godoc.org/github.com/czos/goamz)

## How to build and install goamz

Just use `go get` with any of the available packages. For example:

* `$ go get github.com/czos/goamz/ec2`
* `$ go get github.com/czos/goamz/s3`

## Running tests

To run tests, first install gocheck with:

`$ go get github.com/motain/gocheck`

Then run go test as usual:

`$ go test github.com/czos/goamz/...`

_Note:_ running all tests with the command `go test ./...` will currently fail as tests do not tear down their HTTP listeners.

If you want to run integration tests (costs money), set up the EC2 environment variables as usual, and run:

$ gotest -i
