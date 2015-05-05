// Copyright (c) 2015 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"log"
	"net/http"
	"os"
)

var config Config

func init() {
	log.Println("Configuring...")
	config = Config{}
	LoadConfig(&config, os.Getenv("UBERALLS_CONFIG"))
	LoadConfig(&config, os.Getenv("UBERALLS_SECRETS"))
}

func main() {
	db, err := config.DB()
	if err != nil {
		log.Fatalf("Unable to initialize DB connection: %v", err)
	}

	if err := config.Automigrate(); err != nil {
		log.Fatalf("Could not establish database connection: %v", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/health", HealthHandler{})
	mux.Handle("/metrics", NewMetricsHandler(db))

	listenString := config.ConnectionString()
	log.Printf("Listening on %s... ", listenString)

	log.Fatal(http.ListenAndServe(listenString, mux))
}