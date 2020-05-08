package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/simonmittag/mse6"
	"net"
	"os"
)

func main() {
	port := flag.Int("port", 8081, "the http port")
	waitSecs := flag.Int("wait", 3, "wait time for server to respond in seconds")
	testMode := flag.Bool("test", false, "server self test")
	flag.Parse()

	if !*testMode {
		mse6.Bootstrap(*port, float64(*waitSecs))
	} else {
		selftest(*port)
	}
}

func selftest(port int) {
	_, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err == nil {
		log.Info().Msgf("mse6 self test pass. port %d available", port)
		os.Exit(0)
	} else {
		log.Error().Msgf("mse6 self test fail. port %d unavailable", port)
		os.Exit(1)
	}
}
