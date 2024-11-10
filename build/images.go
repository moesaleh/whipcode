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

package build

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
)

/**
 * Builds a Containerfile for a given language.
 *
 * @param string lang
 * @param string setup
 * @return string Containerfile content
 */
func ContainerFile(lang, setup string) string {
	header := "FROM docker.io/alpine:latest"
	prefix := "RUN apk update --no-cache && apk upgrade --no-cache && apk add --no-cache libc-dev musl-dev "
	suffix := "apk --purge del apk-tools && rm -rf /var/cache/apk /var/lib/apk /lib/apk /etc/apk /sbin/apk /usr/share/apk /usr/lib/apk /usr/sbin/apk /usr/local/apk /usr/bin/apk /usr/local/bin/apk /usr/local/sbin/apk /usr/local/lib/apk /usr/local/share/apk /usr/local/libexec/apk /usr/local/etc/apk"

	setupScript := fmt.Sprintf("images/extra_setup/%s.sh", strings.ToLower(lang))
	if _, err := os.Stat(setupScript); err == nil {
		header += fmt.Sprintf("\nCOPY %s /tmp/setup.sh\n", setupScript)
		suffix = "sh /tmp/setup.sh && rm -f /tmp/setup.sh && " + suffix
	}

	return fmt.Sprintf("%s\n%s%s && %s", header, prefix, setup, suffix)
}

/**
 * Builds images for all languages.
 */
func BuildImages() {
	tempFile := "TEMP_CONTAINERFILE"
	defer func() {
		if err := os.Remove(tempFile); err != nil {
			color.Red("Error removing temporary file: %v", err)
			os.Exit(1)
		}
	}()

	var builds Builds

	if _, err := toml.DecodeFile("images/build.toml", &builds); err != nil {
		color.Red("Could not load build configuration: %v", err)
		os.Exit(1)
	}

	i := 0
	for lang, setup := range builds {
		i++
		containerContent := ContainerFile(lang, setup.Setup)

		if err := os.WriteFile(tempFile, []byte(containerContent), 0644); err != nil {
			color.Red("Error writing Containerfile for %s: %v", lang, err)
			os.Exit(1)
		}

		cmd := exec.Command("podman", "build", "-t", fmt.Sprintf("whipcode-%s", strings.ToLower(lang)), "-f", tempFile, ".")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			color.Red("Error creating stdout pipe for %s: %v", lang, err)
			os.Exit(1)
		}

		if err := cmd.Start(); err != nil {
			color.Red("Error starting build for %s: %v", lang, err)
			os.Exit(1)
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Printf("[%d/%d] [%s] %s\n", i, len(builds), strings.ToLower(lang), scanner.Text())
		}

		if err := cmd.Wait(); err != nil {
			color.Red("Error building image for %s: %v", lang, err)
			os.Exit(1)
		}
	}
	color.Green("All images built successfully.")
}
