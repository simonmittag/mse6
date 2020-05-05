package main

import (
	"flag"
	"github.com/simonmittag/mse6"
)

func main() {
	port := flag.Int("port", 8080, "the http port")
	waitSecs := flag.Int("wait", 3, "wait time for server to respond in seconds")
	flag.Parse()
	mse6.Bootstrap(*port, float64(*waitSecs))
}
