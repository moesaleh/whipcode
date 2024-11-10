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
	"time"

	"golang.org/x/time/rate"
)

/**
 * Creates a new rate limiter.
 *
 * @return *RateLimiter Rate limiter object
 */
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*Client),
	}
}

/**
 * Checks if a client should be rate limited or
 * allowed to continue.
 *
 * @param ip string IP address of the client
 * @param burst int Burst rate
 * @param refill int Refill rate
 * @return *rate.Limiter Rate limiter object
 */
func (rl *RateLimiter) LimitClient(ip string, burst, refill int) *rate.Limiter {
	rl.mu.RLock()
	user, exists := rl.clients[ip]
	rl.mu.RUnlock()
	if !exists {
		rl.mu.Lock()
		defer rl.mu.Unlock()
		limiter := rate.NewLimiter(rate.Every(time.Duration(refill)*time.Second), burst)
		rl.clients[ip] = &Client{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}
	user.lastSeen = time.Now()
	return user.limiter
}

/**
 * Starts a cleanup routine to remove old clients
 * from the rate limiter. Run as a goroutine.
 */
func (rl *RateLimiter) StartCleanup() {
	go func() {
		for {
			time.Sleep(time.Minute)
			rl.mu.Lock()
			for ip, v := range rl.clients {
				if time.Since(v.lastSeen) >= 2*time.Minute {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
}

/**
 * Checks if a client should be rate limited or
 * allowed to continue. This is a wrapper around
 * LimitClient and Allow.
 *
 * @param ip string IP address of the client
 * @param burst int Burst rate
 * @param refill int Refill rate
 * @return bool True if the client should be rate limited
 */
func (rl *RateLimiter) CheckClient(ip string, burst, refill int) bool {
	limiter := rl.LimitClient(ip, burst, refill)
	return limiter.Allow()
}
