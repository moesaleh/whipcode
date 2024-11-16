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

package podman

import (
	"github.com/karlseguin/ccache/v3"
)

/**
 * Struct for defining an executor.
 *
 * @field timeout int Timeout for the executor
 * @field podmanPath string Path to the podman executable
 * @field execCache *ccache.Cache[map[string]interface{}]
 *   Cache for the executor
 */
type Executor struct {
	timeout    int
	podmanPath string
	execCache  *ccache.Cache[map[string]interface{}]
}

/**
 * Struct for defining execution options.
 *
 * @field Code string Code to run
 * @field Entry string Entry point
 * @field Args string Compiler/interpreter arguments
 * @field Stdin string Standard input
 * @field Ext string File extension
 * @field Timeout int Execution timeout
 * @field Env map[string]string Environment variables
 * @field EnableCache bool Enable cache
 */
type ExecutionOptions struct {
	Code        string
	Entry       string
	Args        string
	Stdin       string
	Ext         string
	Timeout     int
	Env         map[string]string
	EnableCache bool
}
