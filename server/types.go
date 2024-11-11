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
	"whipcode/control"
	"whipcode/podman"
)

type contextKey string

type LangMap map[string]map[string]string

/**
 * Struct that's used to pass options to the /run
 * endpoint handler.
 *
 * @field LangMap LangMap Language map
 * @field EnableCache bool Enable cache
 * @field KeyAndSalt []string Key and salt
 * @field MaxBytesSize int Maximum bytes size
 * @field KeyStore *control.KeyStore Cached Key store
 * @field Executor podman.Executor Podman executor
 */
type ScopedMiddlewareParams struct {
	LangMap      LangMap
	EnableCache  bool
	KeyAndSalt   []string
	MaxBytesSize int
	KeyStore     *control.KeyStore
	Executor     podman.Executor
}

/**
 * Struct that's used to pass options to the global
 * middleware
 *
 * @field RateLimiter *control.RateLimiter Rate limiter
 * @field Standalone bool Standalone mode
 * @field RlBurst int Rate limiter burst
 * @field RlRefill int Rate limiter refill
 * @field Proxy string Reverse proxy address
 */
type MiddlewareParams struct {
	RateLimiter *control.RateLimiter
	Standalone  bool
	RlBurst     int
	RlRefill    int
	Proxy       string
}
