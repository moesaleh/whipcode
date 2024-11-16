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

package routes

/**
 * Struct for decoding requests to the /run endpoint.
 *
 * @field Code string Code to run
 * @field LanguageID StrInt ID of the language
 * @field Args string Compiler/interpreter arguments
 * @field Timeout StrInt Execution timeout
 */
type User struct {
	Code       string            `json:"code"`
	LanguageID StrInt            `json:"language_id"`
	Args       string            `json:"args"`
	Timeout    StrInt            `json:"timeout"`
	Stdin      string            `json:"stdin"`
	Env        map[string]string `json:"env"`
}

/**
 * Struct for string + integer values.
 *
 * @field value string Value
 */
type StrInt struct {
	value string
}
