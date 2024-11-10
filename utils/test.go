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

package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
)

/**
 * Generate a form for the user to input the master
 * key and port.
 *
 * @return string Master key
 * @return string Port
 */
func TestForm() (string, string) {
	var key string
	var port string

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Master key").
				Value(&key).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("Master key is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Port").
				Placeholder("8000").
				Value(&port).
				Validate(func(s string) error {
					if s == "" {
						return nil
					}
					if _, err := strconv.Atoi(s); err != nil {
						return fmt.Errorf("Invalid port")
					}
					if i, _ := strconv.Atoi(s); i < 1 || i > 65535 {
						return fmt.Errorf("Invalid port")
					}
					return nil
				}),
		),
	).WithShowHelp(true).WithTheme(huh.ThemeBase()).Run()

	if err != nil {
		if err == huh.ErrUserAborted {
			os.Exit(0)
		}
		color.Red("Could not run test: %v", err)
		os.Exit(1)
	}

	return strings.TrimSpace(key), strings.TrimSpace(port)
}

/**
 * Run the self-test for the application. This
 * will send a request for each language to the
 * server with a test payload.
 */
func SelfTest() {
	key, port := TestForm()

	if port == "" {
		port = "8000"
	}

	var tests Tests

	if _, err := toml.DecodeFile("tests/tests.toml", &tests); err != nil {
		color.Red("Could not load test configuration: %v", err)
		os.Exit(1)
	}

	for id, code := range tests {
		payload := Payload{
			"language_id": id,
			"code":        base64.StdEncoding.EncodeToString([]byte(code.Test)),
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			color.Red("Error encoding JSON for language %s: %v", id, err)
			continue
		}

		url := fmt.Sprintf("http://0.0.0.0:%s/run", port)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			color.Red("Error creating request for language %s: %v", id, err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Master-Key", key)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			color.Red("Error sending request for language %s: %v", id, err)
			continue
		}
		defer resp.Body.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		responseBody := buf.String()

		var Print = color.Red

		if strings.Contains(responseBody, "Success!") {
			Print = color.Green
		}

		Print("%s %s", id, responseBody)
	}
}
