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

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"

	"whipcode/control"
	"whipcode/podman"
	"whipcode/server"
)

func (l *LanguageID) UnmarshalJSON(b []byte) error {
	var intValue int
	if err := json.Unmarshal(b, &intValue); err == nil {
		l.value = strconv.Itoa(intValue)
		return nil
	}

	var strValue string
	if err := json.Unmarshal(b, &strValue); err == nil {
		l.value = strValue
		return nil
	}

	return json.Unmarshal(b, &l.value)
}

func Run(w http.ResponseWriter, r *http.Request) {
	masterKey := r.Header.Get("X-Master-Key")

	if masterKey == "" {
		log.Warn("Blocked the last request", "Reason", "missing master key")
		server.Send(w, http.StatusUnauthorized, []byte(`{"detail": "unauthorized"}`))
		return
	}

	ks, _ := r.Context().Value(server.KeyStoreContextKey).(*control.KeyStore)
	if !ks.CheckKey(masterKey, r.Context().Value(server.MasterKeyContextKey).([]string)) {
		log.Warn("Blocked the last request", "Reason", "invalid master key")
		server.Send(w, http.StatusUnauthorized, []byte(`{"detail": "unauthorized"}`))
		return
	}

	mimeType := r.Header.Get("Content-Type")
	if strings.Split(mimeType, ";")[0] != "application/json" {
		server.Send(w, http.StatusUnsupportedMediaType, []byte(`{"detail": "unsupported media type"}`))
		return
	}

	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		server.Send(w, http.StatusBadRequest, []byte(`{"detail": "invalid request format"}`))
		return
	}

	langMap, _ := r.Context().Value(server.LangMapContextKey).(server.LangMap)
	langConfig, exists := langMap[user.LanguageID.value]
	if !exists {
		server.Send(w, http.StatusBadRequest, []byte(`{"detail": "invalid value for parameter language_id, refer to the documentation"}`))
		return
	}

	codeBytes, err := base64.StdEncoding.DecodeString(user.Code)
	if err != nil || user.Code == "" {
		server.Send(w, http.StatusBadRequest, []byte(`{"detail": "invalid value for parameter code, must be a base64 encoded string"}`))
		return
	}

	entry := langConfig["entry"]
	ext := langConfig["ext"]
	img := langConfig["image"]

	ex, _ := r.Context().Value(server.ExecutorContextKey).(podman.Executor)
	status, result := ex.RunCode(string(codeBytes), entry, user.Args, ext, img, r.Context().Value(server.EnableCacheContextKey).(bool))
	resultBytes, _ := json.Marshal(result)

	server.Send(w, status, resultBytes)
}
