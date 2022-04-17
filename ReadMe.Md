![](mse6.png)

mse6 is a mock HTTP/TLS server for integration tests and edge cases. It supports slow HTTP responses, bad
content encoding, sudden crashes, network termination and other abnormal
behaviour that is otherwise difficult to observe or reproduce under test conditions.

[![Circleci Builds](https://circleci.com/gh/simonmittag/mse6.svg?style=shield)](https://circleci.com/gh/simonmittag/mse6)
[![Docker Build Status](https://img.shields.io/docker/cloud/build/simonmittag/mse6)](https://hub.docker.com/r/simonmittag/mse6/builds)
[![Github Workflows](https://github.com/simonmittag/mse6/workflows/Go/badge.svg)](https://github.com/simonmittag/mse6/actions)
[![Github Issues](https://img.shields.io/github/issues/simonmittag/mse6)](https://github.com/simonmittag/mse6/issues)
[![Github Activity](https://img.shields.io/github/commit-activity/m/simonmittag/mse6)](https://img.shields.io/github/commit-activity/m/simonmittag/mse6)  
[![Go Report](https://goreportcard.com/badge/github.com/simonmittag/mse6)](https://goreportcard.com/report/github.com/simonmittag/mse6)
[![Codeclimate Maintainability](https://api.codeclimate.com/v1/badges/362bc41f687169d50e6b/maintainability)](https://codeclimate.com/github/simonmittag/mse6/maintainability)
[![Codeclimate Test Coverage](https://api.codeclimate.com/v1/badges/362bc41f687169d50e6b/test_coverage)](https://codeclimate.com/github/simonmittag/mse6/test_coverage)
[![Go Version](https://img.shields.io/github/go-mod/go-version/simonmittag/jabba)](https://img.shields.io/github/go-mod/go-version/simonmittag/jabba)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Docker Pulls](https://img.shields.io/docker/pulls/simonmittag/mse6)](https://img.shields.io/docker/pulls/simonmittag/mse6)
[![Docker Image Size](https://img.shields.io/docker/image-size/simonmittag/mse6?sort=date)](https://img.shields.io/docker/image-size/simonmittag/mse6?sort=date)
[![Version](https://img.shields.io/badge/version-0.4.3-orange)](https://github.com/simonmittag/mse6/releases/tag/v0.4.3)

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
    -s self-signed ssl mode
    -t	server self test
    -u string
    	the path prefix (default "/mse6/")
    -v	print the server version
    -w int
      	wait time for server to respond in seconds (default 3)
```

## Routes
`GET /mse6/badcontentlength`
Sends invalid content length header, too large for response

`GET /mse6/badgzip`
Sends a response gzip content encoding header and garbled binary

`GET /mse6/brotli`
Sends a response with br content encoding header and brotli encoded binary response

`GET /mse6/choose`
Sends a HTTP response to the client with one of the following content encodings: `br`, `gzip`, `deflate` or `identity` 
Content encoding preference is in above order and depends on values found in `Accept-Encoding` header found on request. 

`GET /mse6/chunked`
Sends a chunked HTTP/1.1 response to the client

`DELETE /mse6/delete`
Standard json response with status code 204

`GET /mse6/die`
Kills the server process without sending a response

`GET /mse6/echoheader`
echoes all request headers sent on response body for testing

`GET /mse6/echoquery`
echoes all query string values sent with request on response body for testing

`GET /mse6/jwksmix`
sends a list of mixed algorithm Jwks keys

`GET /mse6/jwksbad`
sends illegally formatted Jwks key

`GET /mse6/jwksbadrotate?rc=0`
sends a rotating Jwks key that alternates every request. Sends malformed keys. Reset good behaviour with rc=0

`GET /mse6/get`
Standard json response with status code 200

`HEAD /mse6/getorhead?cl=y`
`GET /mse6/getorhead`
Sends a HTTP HEAD response or a GET response, depending on request method.
If query param cl is appended with any value, the Content-Length header of 
the body otherwise sent by GET is also included in the HEAD response

`GET /mse6/gzip`
Sends a gzipped response with proper content encoding

`GET /mse6/hangupduringheader`
Sends a partial header only response, waits 2s, then closes the TCP connection.

`GET /mse6/hangupafterheader`
Sends a complete header only response message , waits 2s, then closes the TCP connection.

`GET /mse6/hangupduringbody`
Sends a complete header message, then some of the body, waits 2s, then closes the TCP connection.

`OPTIONS /mse6/options?code=n`
Sends a HTTP OPTIONS response as per RFC7231, section 4.3.7. Contains Allow headers and a status code.
Legal status codes are 200 and 204, the rest is undefined.

`GET /mse6/nocontentenc`
Sends a HTTP response without a content encoding header set

`POST /mse6/post`
Standard json response with status code 201

`PATCH /mse6/patch`
Standard json response with status code 200

`PUT /mse6/put`
Standard json response with status code 200

`GET /mse6/send?code=nnn&url=http%3A%2F%2Fwww.google.com`
Sends arbitrary response code between 100 and 999. For redirects, you can supply a custom
location parameter. Don't forget to URL encode your params.

`GET /mse6/slowheader?wait=n`
Sends headers but only after waiting for n seconds. 
Alternatively configure default with -w=n on cli

`GET /mse6/slowbody?wait=n`
Sends body after initial lag of n/2s, then sends remaining body without chunking after n/2s. 
Alternatively configure default with -w=n on cli

`TRACE /mse6/trace`
Standard json response with status code 200 and "message/http" content type

`GET /mse6/tiny`
Tiny JSON response with content encoding identity

`GET /mse6/tinygzip`
Tiny JSON response with content encoding gzip

`GET /mse6/unknowncontentenc`
Sends unknown content-encoding header with json response.


`GET /mse6/websocket?n=1&c=true&c1=true&c2=true`
Upgrade the client connection to a HTTP/1.1 websocket that echoes messages sent to it.
Specify n for how often the echo should repeat. Specify c to orderly close the connection
from the server side after sending all echo responses. Specify c1 to only send websocket
protocol close. Specify c2 to only hang up on TCP connection, without respecting websocket protocol.

## Contributions
The mse6 team welcomes all [contributors](https://github.com/simonmittag/mse6/blob/master/CONTRIBUTING.md). Everyone interacting with the project's codebase, issue trackers, chat rooms and mailing lists
is expected to follow the [code of conduct](https://github.com/simonmittag/mse6/blob/master/CODE_OF_CONDUCT.md)
