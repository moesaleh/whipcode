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

package control

import (
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

/**
 * Struct that holds the cached master key.
 *
 * @field cachedKey atomic.Value Cached master key
 */
type KeyStore struct {
	cachedKey atomic.Value
}

/**
 * Struct that holds rate limit status for clients.
 *
 * @field limiter *rate.Limiter Rate limiter
 * @field lastSeen time.Time Last seen time
 */
type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

/**
 * Struct that holds a map of clients.
 *
 * @field clients map[string]*Client Map of clients
 * @field mu sync.RWMutex Mutex for the map
 */
type RateLimiter struct {
	clients map[string]*Client
	mu      sync.RWMutex
}
