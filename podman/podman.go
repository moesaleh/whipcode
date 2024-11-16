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
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/karlseguin/ccache/v3"
)

/**
 * Creates a new LRU cache for caching exec results
 * and a new Executor instance.
 *
 * @param timeout int Timeout for execution
 * @param podmanPath string Path to podman executable
 * @return *Executor New Executor instance
 */
func NewExecutor(timeout int, podmanPath string) *Executor {
	cache := ccache.New(ccache.Configure[map[string]interface{}]().MaxSize(100).ItemsToPrune(10))
	return &Executor{execCache: cache, timeout: timeout, podmanPath: podmanPath}
}

/**
 * Cleans up the temp directory. Called on SIGINT and
 * SIGTERM.
 */
func Cleanup() {
	if err := os.RemoveAll(filepath.Join(".", "run")); err != nil {
		log.Error("Could not clean up temp dir", "Error", err)
	}
}

/**
 * Sanitizes the given string from shell injection.
 *
 * @param String string Code to sanitize
 * @return string Sanitized code
 */
func Sanitize(String string) string {
	slices := strings.Fields(String)
	for i, slice := range slices {
		slices[i] = "'" + strings.ReplaceAll(slice, "'", "'\\''") + "'"
	}
	return strings.Join(slices, " ")
}

/**
 * Runs the given code in a podman container. The code is
 * dumped into a temp file, which is then mounted into the
 * container.
 *
 * @param code string Code to run
 * @param entry string Entry point for the container
 * @param cArgs string Args for the interpreter or compiler
 * @param ext string File extension of the code
 * @param timeout int Timeout for execution
 * @param enableCache bool Enable caching of exec results
 * @return int HTTP status code
 * @return map[string]interface{} Response body
 */
func (ex *Executor) RunCode(opt ExecutionOptions) (int, map[string]interface{}) {
	cArgs, stdin := Sanitize(opt.Args), Sanitize(opt.Stdin)

	if opt.EnableCache {
		if item := ex.execCache.Get(cArgs + opt.Entry + opt.Code); item != nil {
			go item.Extend(time.Hour * 24)
			return http.StatusOK, item.Value()
		}
	}

	boxID := strconv.Itoa(rand.Intn(9000000) + 1000000)
	srcFileName := fmt.Sprintf("run%s.%s", boxID, opt.Ext)
	srcFilePath := filepath.Join(".", "run", srcFileName)

	if err := os.WriteFile(srcFilePath, []byte(opt.Code), 0644); err != nil {
		log.Error("Could not write to temp file", "Error", err)
		return http.StatusInternalServerError, map[string]interface{}{
			"detail": "internal server error",
		}
	}
	defer os.Remove(srcFilePath)

	thisTimeout := opt.Timeout
	if opt.Timeout == 0 || opt.Timeout > ex.timeout {
		thisTimeout = ex.timeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(thisTimeout)*time.Second)
	defer cancel()

	var stdout, stderr bytes.Buffer
	args := []string{
		"run",
		"--rm",
		"--read-only",
		"--no-hosts",
		"--hostname", "box" + boxID,
		"--network", "none",
		"--timeout", strconv.Itoa(thisTimeout + 1),
		"--cap-drop", "ALL",
		"--memory", "512m",
		"--memory-reservation", "128m",
		"--cpus", "1.0",
		"--pids-limit", "32",
		"--user", "nobody",
		"--tmpfs", "/tmp:rw,size=64m,mode=1777",
		"--tmpfs", "/var/tmp:ro,size=32m,mode=1777",
		"--security-opt", "no-new-privileges",
		"--security-opt", "mask=/run:/sys:/var",
		"--security-opt", "label=type:whipcode.process",
		"--security-opt", "proc-opts=hidepid=2,subset=pid",
		"--unsetenv", "container",
		"--volume", fmt.Sprintf("./entry/%s.sh:/entry.sh:z,ro", opt.Entry),
		"--volume", fmt.Sprintf("./run/%s:/source.%s:Z,ro", srcFileName, opt.Ext),
	}
	for k, v := range opt.Env {
		args = append(args, "--env", k+"="+v)
	}
	args = append(
		args,
		"whipcode-"+opt.Entry,
		"sh", "-c", fmt.Sprintf("echo stdout-start && echo stderr-start >&2 && echo %s | sh entry.sh %s", stdin, cArgs),
	)

	cmdExec := exec.CommandContext(ctx, ex.podmanPath, args...)
	cmdExec.Stdout = &stdout
	cmdExec.Stderr = &stderr

	startTime := time.Now()
	_ = cmdExec.Run()
	duration := time.Since(startTime).Seconds()

	if ctx.Err() == context.DeadlineExceeded {
		result := map[string]interface{}{
			"stdout":        "",
			"stderr":        "",
			"container_age": duration,
			"timeout":       true,
		}

		if opt.EnableCache {
			go ex.execCache.Set(cArgs+opt.Entry+opt.Code, result, time.Hour*24)
		}

		return http.StatusOK, result
	}

	stdoutStr := stdout.String()
	stderrStr := stderr.String()
	if !strings.HasPrefix(stdoutStr, "stdout-start") || !strings.HasPrefix(stderrStr, "stderr-start") {
		log.Warn("Caught unsafe output", "STDOUT", stdoutStr, "STDERR", stderrStr)
		return http.StatusInternalServerError, map[string]interface{}{
			"detail": "internal server error",
		}
	}

	result := map[string]interface{}{
		"stdout":        strings.TrimPrefix(stdoutStr, "stdout-start\n"),
		"stderr":        strings.TrimPrefix(stderrStr, "stderr-start\n"),
		"container_age": duration,
		"timeout":       false,
	}

	if opt.EnableCache {
		go ex.execCache.Set(cArgs+opt.Entry+opt.Code, result, time.Hour*24)
	}

	return http.StatusOK, result
}
