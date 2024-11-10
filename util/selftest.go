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

package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
)

func SelfTest() {
	var key string
	var port string

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Master key").
				Value(&key),
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
		color.Red("Could not run self-test: %v", err)
		os.Exit(1)
	}

	key = strings.TrimSpace(key)
	port = strings.TrimSpace(port)

	if port == "" {
		port = "8000"
	}

	if key == "" {
		color.Red("Master key is required")
		return
	}

	for id, code := range Tests {
		payload := map[string]string{
			"language_id": fmt.Sprintf("%d", id),
			"code":        base64.StdEncoding.EncodeToString([]byte(code)),
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			color.Red("Error encoding JSON for language %d: %v", id, err)
			continue
		}

		url := fmt.Sprintf("http://0.0.0.0:%s/run", port)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			color.Red("Error creating request for language %d: %v", id, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Master-Key", key)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			color.Red("Error sending request for language %d: %v", id, err)
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

		Print("%d %s", id, responseBody)
	}
}
