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

package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

/**
 * Starts the server with the given port, handler, and
 * TLS settings.
 *
 * @param port int Port to use
 * @param handler http.Handler Handler to use
 * @param enableTLS bool Whether to enable TLS
 * @param tlsDir string Directory for the TLS files
 * @param timeout int Configured execution timeout
 */
func StartServer(port int, handler http.Handler, enableTLS bool, tlsDir string, timeout int) {
	log.Info("Starting whipcode", "Port", port, "TLS", enableTLS)

	addr := fmt.Sprintf(":%d", port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: (time.Duration(timeout) + 10) * time.Second,
		IdleTimeout:  20 * time.Second,
	}

	var Serve = srv.ListenAndServe

	if enableTLS {
		Serve = func() error {
			return srv.ListenAndServeTLS(tlsDir+"/cert.pem", tlsDir+"/key.pem")
		}
	}

	if err := Serve(); err != nil {
		log.Fatal("Server failed", "Error", err)
	}
}
