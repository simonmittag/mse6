![](mse6.png)

mse6 is a http server implementing mock responses for testing of edge cases, such as slow writes.

[![Github Workflows](https://github.com/simonmittag/mse6/workflows/Go/badge.svg)](https://github.com/simonmittag/mse6/actions)
[![Github Issues](https://img.shields.io/github/issues/simonmittag/mse6)](https://github.com/simonmittag/mse6/issues)
[![Github Activity](https://img.shields.io/github/commit-activity/m/simonmittag/mse6)](https://img.shields.io/github/commit-activity/m/simonmittag/mse6)  
[![Go Version](https://img.shields.io/github/go-mod/go-version/simonmittag/jabba)](https://img.shields.io/github/go-mod/go-version/simonmittag/jabba)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Features
* Configurable HTTP/1.1 port
* Configurable wait time for slow responses

## Routes
```bash
/mse6/getting 
/mse6/posting
/mse6/slowheader
/mse6/slowbody
/mse6/gzip
```

## Up and running
```bash
go install github.com/simonmittag/mse6/cmd/mse6 && mse6
```

## Usage
```bash
λ mse6 -help
Usage of mse6:
  -port int
        the http port (default 8080)
  -wait int
        wait time for server to respond in seconds (default 3)
```

## Contributions
The mse6 team welcomes all [contributors](https://github.com/simonmittag/mse6/blob/master/CONTRIBUTING.md). Everyone interacting with the project's codebase, issue trackers, chat rooms and mailing lists
is expected to follow the [code of conduct](https://github.com/simonmittag/mse6/blob/master/CODE_OF_CONDUCT.md)
