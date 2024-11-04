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

import (
	"whipcode/server"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/log"
)

func LoadConfig(path string) *Config {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Fatal("Could not load config", "File", path, "Error", err)
	}
	return &config
}

func LoadLangs(path string) *server.LangMap {
	var langs = make(server.LangMap)
	if _, err := toml.DecodeFile(path, &langs); err != nil {
		log.Fatal("Could not load language map", "File", path, "Error", err)
	}
	return &langs
}
