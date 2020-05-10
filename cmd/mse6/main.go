package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/simonmittag/mse6"
	"net"
	"os"
)

type Mode uint8

const (
	Server Mode = 1 << iota
	Test
	Version
)

func main() {
	mode := Server
	port := flag.Int("p", 8081, "the http port")
	waitSecs := flag.Int("w", 3, "wait time for server to respond in seconds")
	tM := flag.Bool("t", false, "server self test")
	vM := flag.Bool("v", false, "print the server version")
	flag.Parse()
	if *tM {
		mode = Test
	}
	if *vM {
		mode = Version
	}

	switch mode {
	case Server:
		mse6.Bootstrap(*port, float64(*waitSecs))
	case Test:
		printSelftest(*port)
	case Version:
		printVersion()
	}
}

func printVersion() {
	fmt.Printf("mse6 %s\n", mse6.Version)
	os.Exit(0)
}

func printSelftest(port int) {
	_, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err == nil {
		log.Info().Msgf("mse6 %s self test pass. port %d available", mse6.Version, port)
		os.Exit(0)
	} else {
		log.Error().Msgf("mse6 %s self test fail. port %d unavailable", mse6.Version, port)
		os.Exit(1)
	}
}
