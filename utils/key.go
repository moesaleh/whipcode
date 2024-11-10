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
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"golang.org/x/crypto/argon2"
)

func RandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		color.Red("Could not generate random string: %v", err)
		os.Exit(1)
	}
	return hex.EncodeToString(bytes)
}

func KeyForm() (string, string) {
	var salt string
	var key string

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Salt").
				Placeholder("Leave empty for random salt").
				Value(&salt),
			huh.NewInput().
				Title("Key").
				Placeholder("Leave empty for random key").
				Value(&key),
		),
	).WithShowHelp(true).WithTheme(huh.ThemeBase()).Run()

	if err != nil {
		if err == huh.ErrUserAborted {
			os.Exit(0)
		}
		color.Red("Could not generate key: %v", err)
		os.Exit(1)
	}

	return strings.TrimSpace(salt), strings.TrimSpace(key)
}

func GenKey() {
	salt, key := KeyForm()

	if salt == "" {
		salt = RandomString(16)
	}

	if key == "" {
		key = RandomString(32)
		fmt.Println("This is your master key:", key)
		fmt.Println("It won't be shown again, so make sure to save it somewhere safe.")
	}

	hash := argon2.IDKey([]byte(key), []byte(salt), 1, 4096, 1, 32)

	file, err := os.Create(".masterkey")
	if err != nil {
		color.Red("Could not create .masterkey: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	file.WriteString(hex.EncodeToString(hash) + "\n" + salt)
	color.Green("Hash and salt saved to .masterkey")
}
