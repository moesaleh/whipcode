//
//  Copyright 2024 whipcode.app (AnnikaV9)
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//          http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on
//  an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific
//  language governing permissions and limitations under the License.
//

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/charmbracelet/log"

	"whipcode/config"
	"whipcode/control"
	"whipcode/podman"
	"whipcode/routes"
	"whipcode/server"
	"whipcode/utils"
)

const VERSION = "1.5.1"

func main() {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      "2006-01-02 15:04:05",
	})
	log.SetDefault(logger)

	fileConfig := config.LoadConfig("config.toml")

	var version, enableTLS, enableCache, enablePing, standalone, genKey, selfTest bool
	var port, maxBytesSize, rlBurst, rlRefill, timeout int
	var keyFile, proxy, podmanPath, tlsDir, langMap string

	flag.Usage = func() {
		fmt.Printf("usage: %s [options]\n", os.Args[0])
		fmt.Println(`
commands:
    --gen-key                 generate a master key
    --self-test               run self test`)
		fmt.Println(`
options:
    -h, --help                print this help message
    -v, --version             print version information
    -p, --port       PORT     port to listen on
    -m, --max-bytes  BYTES    max bytes to accept
    -t, --timeout    SECONDS  timeout for execution
    -k, --key        FILE     master key file
    --lang-map       FILE     language map file
    --podman-path    PATH     path to podman
    --proxy          ADDR     reverse proxy address
    --cache                   enable execution cache
    --tls                     enable tls
    --tls-dir        DIR      directory with cert and key
    --ping                    enable /ping endpoint
    --standalone              enable rate limiting (CHECK README)
    --burst          COUNT    rate limit burst
    --refill	     SECONDS  rate limit refill time`)
		fmt.Println("\nsee config.toml for default values")
	}
	flag.BoolVar(&genKey, "gen-key", false, "")
	flag.BoolVar(&selfTest, "self-test", false, "")
	flag.BoolVar(&version, "version", false, "")
	flag.BoolVar(&version, "v", false, "")
	flag.IntVar(&port, "port", fileConfig.Port, "")
	flag.IntVar(&port, "p", fileConfig.Port, "")
	flag.IntVar(&maxBytesSize, "max-bytes", fileConfig.MaxBytes, "")
	flag.IntVar(&maxBytesSize, "m", fileConfig.MaxBytes, "")
	flag.IntVar(&timeout, "timeout", fileConfig.Timeout, "")
	flag.IntVar(&timeout, "t", fileConfig.Timeout, "")
	flag.StringVar(&keyFile, "key", fileConfig.Key, "")
	flag.StringVar(&keyFile, "k", fileConfig.Key, "")
	flag.StringVar(&langMap, "lang-map", fileConfig.LangMap, "")
	flag.StringVar(&podmanPath, "podman-path", fileConfig.PodmanPath, "")
	flag.StringVar(&proxy, "proxy", fileConfig.Proxy, "")
	flag.BoolVar(&enableCache, "cache", fileConfig.Cache, "")
	flag.BoolVar(&enableTLS, "tls", fileConfig.TLS, "")
	flag.StringVar(&tlsDir, "tls-dir", fileConfig.TLSDir, "")
	flag.BoolVar(&enablePing, "ping", fileConfig.Ping, "")
	flag.BoolVar(&standalone, "standalone", fileConfig.Standalone, "")
	flag.IntVar(&rlBurst, "burst", fileConfig.Burst, "")
	flag.IntVar(&rlRefill, "refill", fileConfig.Refill, "")
	flag.Parse()

	switch {
	case version:
		fmt.Println(VERSION)
		return

	case genKey:
		utils.GenKey()
		return

	case selfTest:
		utils.SelfTest()
		return
	}

	if _, err := os.Stat(podmanPath); os.IsNotExist(err) {
		log.Fatal("Podman binary not found", "Error", err)
	}

	if err := os.MkdirAll(filepath.Join(".", "run"), 0755); err != nil {
		log.Fatal("Could not create temp dir", "Error", err)
	}

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-exitChan
		podman.Cleanup()
		os.Exit(0)
	}()

	keyStore, keyAndSalt := control.InitializeKeystore(keyFile)

	scopedParams := server.ScopedMiddleWareParams{
		LangMap:      *config.LoadLangs(langMap),
		EnableCache:  enableCache,
		KeyAndSalt:   keyAndSalt,
		KeyStore:     keyStore,
		MaxBytesSize: maxBytesSize,
		Executor:     *podman.NewExecutor(timeout, podmanPath),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		server.Send(w, http.StatusNotFound, []byte(`{"detail": "not found"}`))
	})

	http.HandleFunc("POST /run", server.ScopedMiddleWare(routes.Run, scopedParams))
	http.HandleFunc("/run", func(w http.ResponseWriter, _ *http.Request) {
		server.Send(w, http.StatusMethodNotAllowed, []byte(`{"detail": "method not allowed"}`))
	})

	if enablePing {
		http.HandleFunc("/ping", routes.Ping)
	}

	rateLimiter := control.NewRateLimiter()
	if standalone {
		rateLimiter.StartCleanup()
	}

	params := server.MiddleWareParams{
		RateLimiter: rateLimiter,
		Standalone:  standalone,
		RlBurst:     rlBurst,
		RlRefill:    rlRefill,
		Proxy:       proxy,
	}

	handler := server.MiddleWare(http.DefaultServeMux, params)
	server.StartServer(port, handler, enableTLS, tlsDir, timeout)
}
