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
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/charmbracelet/log"
)

const (
	LangMapContextKey     contextKey = "langMap"
	MasterKeyContextKey   contextKey = "masterKey"
	KeyStoreContextKey    contextKey = "keyStore"
	EnableCacheContextKey contextKey = "enableCache"
	ExecutorContextKey    contextKey = "executor"
)

/**
 * Middleware for the /run endpoint that caps the
 * request body size and passes various parameters
 * to the handler.
 *
 * @param f http.HandlerFunc Handler
 * @param params ScopedMiddleWareParams Parameters
 * @return http.HandlerFunc Handler
 */
func ScopedMiddleWare(f http.HandlerFunc, params ScopedMiddleWareParams) http.HandlerFunc {
	/**
	 * @param w http.ResponseWriter Response writer
	 * @param r *http.Request Request object
     */
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, int64(params.MaxBytesSize))

		ctx := r.Context()
		ctx = context.WithValue(ctx, LangMapContextKey, params.LangMap)
		ctx = context.WithValue(ctx, MasterKeyContextKey, params.KeyAndSalt)
		ctx = context.WithValue(ctx, KeyStoreContextKey, params.KeyStore)
		ctx = context.WithValue(ctx, EnableCacheContextKey, params.EnableCache)
		ctx = context.WithValue(ctx, ExecutorContextKey, params.Executor)

		f(w, r.WithContext(ctx))
	}
}

/**
 * Global middleware for all requests that performs
 * rate limiting and host checks.
 *
 * @param handler http.Handler Handler
 * @param params MiddleWareParams Parameters
 * @return http.Handler Handler
 */
func MiddleWare(handler http.Handler, params MiddleWareParams) http.Handler {
	/**
	 * @param w http.ResponseWriter - Response writer
	 * @param r *http.Request - Request object
	 */
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		details := fmt.Sprintf("%s %s %s", host, r.Method, r.URL)

		if params.Proxy != "" && host != params.Proxy {
			log.Warn(details, "Blocked", "host not allowed")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if params.Standalone && !params.RateLimiter.CheckClient(host, params.RlBurst, params.RlRefill) {
			log.Info(details, "Blocked", "rate limit exceeded")
			Send(w, http.StatusTooManyRequests, []byte(`{"detail": "you are sending too many requests"}`))
			return
		}

		log.Info(details)

		handler.ServeHTTP(w, r)
	})
}
