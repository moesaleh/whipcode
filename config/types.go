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

package config

/**
 * Struct for defining configuration options.
 *
 * @field Port int Port to listen on
 * @field Addr string Address to listen on
 * @field MaxBytes int Max bytes to accept
 * @field Proxy string Reverse proxy address
 * @field TLS bool Enable tls
 * @field TLSDir string Directory with cert and key
 * @field Ping bool Enable /ping endpoint
 * @field LangMap string Path to the language map
 * @field PodmanPath string Path to podman
 * @field Timeout int Timeout for executions
 * @field Key string Master key file
 * @field Cache bool Enable execution cache
 * @field Standalone bool Enable rate limiting
 * @field Burst int Burst for the rate limiter
 * @field Refill int Refill for the rate limiter
 */
type Config struct {
	Port       int
	Addr       string
	MaxBytes   int
	Proxy      string
	TLS        bool
	TLSDir     string
	Ping       bool
	LangMap    string
	PodmanPath string
	Timeout    int
	Key        string
	Cache      bool
	Standalone bool
	Burst      int
	Refill     int
}
