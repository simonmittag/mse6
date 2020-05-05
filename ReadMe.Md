![](jabba.png)

mse6 is a test http server with configurable slow responses and other edge cases.

## Features
* Configurable HTTP/1.1 port
* Configurable wait time for slow responses
* /mse6/getting endpoint sends json response
* /mse6/posting endpoint sends json response
* /mse6/slowheader endpoint sends json response but waits before sending headers
* /mse6/slowbody endpoint sends json response but waits while sending body and sends increments without chunked encoding
* /mse6/gzip endpoint sends gzip body

## Up and running

```bash
go install github.com/simonmittag/mse6/cmd/mse6 && mse6
```

## Contributions
The mse6 team welcomes all [contributors](https://github.com/simonmittag/mse6/blob/master/CONTRIBUTING.md). Everyone interacting with the project's codebase, issue trackers, chat rooms and mailing lists
is expected to follow the [code of conduct](https://github.com/simonmittag/mse6/blob/master/CODE_OF_CONDUCT.md)
