![](mse6.png)

mse6 is a http server implementing mock responses for testing of edge cases, such as slow writes.

[![Github Workflows](https://github.com/simonmittag/mse6/workflows/Go/badge.svg)](https://github.com/simonmittag/mse6/actions)
[![Github Issues](https://img.shields.io/github/issues/simonmittag/mse6)](https://github.com/simonmittag/mse6/issues)
[![Github Activity](https://img.shields.io/github/commit-activity/m/simonmittag/mse6)](https://img.shields.io/github/commit-activity/m/simonmittag/mse6)  
[![Codeclimate Maintainability](https://api.codeclimate.com/v1/badges/362bc41f687169d50e6b/maintainability)](https://codeclimate.com/github/simonmittag/mse6/maintainability)
[![Codeclimate Test Coverage](https://api.codeclimate.com/v1/badges/362bc41f687169d50e6b/test_coverage)](https://codeclimate.com/github/simonmittag/mse6/test_coverage)
[![Go Version](https://img.shields.io/github/go-mod/go-version/simonmittag/jabba)](https://img.shields.io/github/go-mod/go-version/simonmittag/jabba)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Features
* Configurable HTTP/1.1 port
* Configurable wait time for slow responses

## Routes
### /mse6/die
kills the server process without sending a response

### /mse6/get
standard json response with status code 200

### /mse6/post
standard json response with status code 200

### /mse6/slowheader
sends headers but only after waiting for 3 seconds. Configure with -wait on cli

### /mse6/slowbody
sends body after initial lag, then sends some more body without chunking. Configure with -wait on cli

### /mse6/badcontentlength
sends invalid content length header, too large for response

### /mse6/send?code=nnn
sends arbitrary response code between 100 and 999

### /mse6/gzip
sends a gzipped response with proper content encoding

## Up and running
### Osx
```
brew tap simonmittag/homebrew-mse6 &&
brew install mse6 &&
mse6
```

### Golang
```bash
go get github.com/simonmittag/mse6 && 
go install github.com/simonmittag/mse6/cmd/mse6 && 
mse6
```

## Usage
```
λ mse6 -help
Usage of mse6:
  -port int
        the http port (default 8080)
  -test
        server self test
  -wait int
        wait time for server to respond in seconds (default 3)
```

## Contributions
The mse6 team welcomes all [contributors](https://github.com/simonmittag/mse6/blob/master/CONTRIBUTING.md). Everyone interacting with the project's codebase, issue trackers, chat rooms and mailing lists
is expected to follow the [code of conduct](https://github.com/simonmittag/mse6/blob/master/CODE_OF_CONDUCT.md)
