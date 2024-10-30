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
	"log"
	"net/http"
	"time"
)

func StartServer(port int, handler http.Handler, enableTLS bool) {
	log.Printf("Starting whipcode on port: %d", port)

	addr := fmt.Sprintf(":%d", port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  90 * time.Second,
	}

	var err error
	if enableTLS {
		err = srv.ListenAndServeTLS("tls/cert.pem", "tls/key.pem")
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil {
		log.Fatalf("Fatal: %v", err)
	}
}
