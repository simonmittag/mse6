![](mse6.png)

mse6 is a mock server doing everything in it's power to get HTTP and TLS all wrong.

[![Circleci Builds](https://circleci.com/gh/simonmittag/mse6.svg?style=shield)](https://circleci.com/gh/simonmittag/mse6)
[![Docker Build Status](https://img.shields.io/docker/cloud/build/simonmittag/mse6)](https://hub.docker.com/r/simonmittag/mse6/builds)
[![Github Workflows](https://github.com/simonmittag/mse6/workflows/Go/badge.svg)](https://github.com/simonmittag/mse6/actions)
[![Github Issues](https://img.shields.io/github/issues/simonmittag/mse6)](https://github.com/simonmittag/mse6/issues)
[![Github Activity](https://img.shields.io/github/commit-activity/m/simonmittag/mse6)](https://img.shields.io/github/commit-activity/m/simonmittag/mse6)  
[![Codeclimate Maintainability](https://api.codeclimate.com/v1/badges/362bc41f687169d50e6b/maintainability)](https://codeclimate.com/github/simonmittag/mse6/maintainability)
[![Codeclimate Test Coverage](https://api.codeclimate.com/v1/badges/362bc41f687169d50e6b/test_coverage)](https://codeclimate.com/github/simonmittag/mse6/test_coverage)
[![Go Version](https://img.shields.io/github/go-mod/go-version/simonmittag/jabba)](https://img.shields.io/github/go-mod/go-version/simonmittag/jabba)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Docker Pulls](https://img.shields.io/docker/pulls/simonmittag/mse6)](https://img.shields.io/docker/pulls/simonmittag/mse6)
[![Docker Image Size](https://img.shields.io/docker/image-size/simonmittag/mse6?sort=date)](https://img.shields.io/docker/image-size/simonmittag/mse6?sort=date)


## Up and running
### Docker
```
docker pull simonmittag/mse6 &&
docker run -p8081:8081 simonmittag/mse6
```

### Osx
```
brew install simonmittag/cli/mse6 &&
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
λ mse6 -h
  Usage of mse6:
    -p int
      	the http port (default 8081)
    -t	server self test
    -u string
    	the path prefix (default "/mse6/")
    -v	print the server version
    -w int
      	wait time for server to respond in seconds (default 3)
```

## Routes
`GET /mse6/die`
Kills the server process without sending a response

`GET /mse6/get`
Standard json response with status code 200

`POST /mse6/post`
Standard json response with status code 200

`GET /mse6/echoheader`
echoes all request headers sent on response body for testing

`GET /mse6/slowheader?wait=n`
Sends headers but only after waiting for n seconds. 
Alternatively configure default with -w=n on cli

`GET /mse6/slowbody?wait=n`
Sends body after initial lag of n/2s, then sends remaining body without chunking after n/2s. 
Alternatively configure default with -w=n on cli

`GET /mse6/badcontentlength`
Sends invalid content length header, too large for response

`GET /mse6/send?code=nnn&url=http%3A%2F%2Fwww.google.com`
Sends arbitrary response code between 100 and 999. For redirects, you can supply a custom
location parameter. Don't forget to URL encode your params.

`GET /mse6/gzip`
Sends a gzipped response with proper content encoding

`GET /mse6/badgzip`
Sends a response gzip content encoding header and garbled binary

## Contributions
The mse6 team welcomes all [contributors](https://github.com/simonmittag/mse6/blob/master/CONTRIBUTING.md). Everyone interacting with the project's codebase, issue trackers, chat rooms and mailing lists
is expected to follow the [code of conduct](https://github.com/simonmittag/mse6/blob/master/CODE_OF_CONDUCT.md)
