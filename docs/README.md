<img alt="Whipcode" src="https://github.com/user-attachments/assets/b68d8164-cbbe-45cb-9f66-43618a0e8544"><br/>

<a href="https://go.dev"><img height="20px" alt="Go badge" src="https://github.com/user-attachments/assets/c5115760-24b7-4272-8b48-7b8071e5053d"></a> <a href="https://github.com/containers/podman"><img height="20px" alt="Podman badge" src="https://github.com/user-attachments/assets/89c586f0-6932-49f3-a12f-bbc9b52c2c4f"></a>
 <a href="/LICENSE"><img height="20px" alt="Apache badge" src="https://github.com/user-attachments/assets/2c52fd74-66d1-45a4-a825-2f51c72eedf8"></a>

REST API for executing untrusted code with Podman.

Implemented as a scalable, stateless microservice with no user management or authentication.

> Don't want to self host? Head over to [whipcode.app](https://whipcode.app) to access the live endpoint.
>
> Drop us an email at [hello@whipcode.app](mailto:hello@whipcode.app) if you'd like to get in touch.

<details>
  <summary>Click to see default languages</summary>
  <br/>

| ID  | Language       | Environment   |
| --- | -------------- | ------------- |
| 1   | Python         | cpython       |
| 2   | JavaScript     | node.js       |
| 3   | Bash           | -             |
| 4   | Perl           | -             |
| 5   | Lua            | -             |
| 6   | Ruby           | -             |
| 7   | C              | gcc           |
| 8   | C++            | gcc           |
| 9   | Rust           | -             |
| 10  | Fortran        | gfortran      |
| 11  | Haskell        | runghc        |
| 12  | Java           | openjdk       |
| 13  | Go             | gccgo         |
| 14  | TypeScript     | swc > node.js |
| 15  | Common Lisp    | sbcl          |
| 16  | Racket         | -             |
| 17  | Crystal        | -             |
| 18  | Clojure        | -             |
| 19  | x86 Assembly   | nasm          |
| 20  | Zig            | -             |
| 21  | Nim            | -             |
| 22  | D              | gdc           |
| 23  | C#             | mono          |
| 24  | Rscript        | -             |
| 25  | Dart           | -             |
| 26  | VB.NET         | mono          |
| 27  | F#             | mono          |
| 28  | PHP            | -             |

To add languages, see:
- [images/build.toml](/images/build.toml)
- [images/extra_setup/](/images/extra_setup/)
- [entry/](/entry/)
- [langmap.toml](/langmap.toml)
- [tests/tests.toml](tests/tests.toml)

</details>

## Table of contents
- [Setting up](#setting-up)
  - [Environment setup](#environment-setup)
  - [Building](#building)
- [Starting the service](#starting-the-service)
- [Systemd](#systemd)
- [CLI options](#cli-options)
- [API reference](#api-reference)
  - [Headers](#headers)
  - [Body](#body)
  - [Response](#response)
  - [Example request](#example-request)
  - [Example response](#example-response)
- [Tasks](#tasks)
- [Contributing](#contributing)
- [Credits](#credits)
- [License](#license)

## Setting up
**Please set this up on a SELinux-enabled system.**

Tested on Fedora Server 41, Go 1.23

### Environment setup
Go:
```bash
# Download Go tarball (> 1.23)
wget https://go.dev/dl/go1.<version>.linux-amd64.tar.gz

# Clean up previous installations
sudo rm -rf /usr/local/go

# Extract the tarball into /usr/local
sudo tar -C /usr/local -xzf go1.<version>.linux-amd64.tar.gz

# Add to PATH
sudo echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile

# Load new PATH
source /etc/profile
```

go-task:
```bash
# Download go-task rpm
wget https://github.com/go-task/task/releases/latest/download/task_linux_amd64.rpm

# Install the rpm
sudo dnf install task_linux_amd64.rpm
```

Podman:
```bash
# Install podman
sudo dnf install podman

# Allocate uids and gids
sudo echo "$USER:100000:65536" | sudo tee /etc/subuid /etc/subgid

# Reset podman
podman system reset

# Should return 'overlay'
podman info | grep graphDriverName
```

SELinux:
```bash
# Should return 'Enforcing'
sudo getenforce

# Install container specific policies
sudo dnf install container-selinux

# Load whipcode's own policy
sudo semodule -i selinux/whipcode.cil selinux/base_container.cil
```

### Building
Use `task <command>` to run predefined build actions:

| Command          | Action                                        |
| ---------------- | --------------------------------------------- |
| `all`            | Build everything.                             |
| `build`          | Build only whipcode.                          |
| `build-images`   | Build only the container images.              |
| `rebuild-images` | Clean rebuild images.                         |
| `update`         | Update (git pull), build whipcode and images. |

See the [Tasks](#tasks) section for more non-build actions.

## Starting the service

> [!WARNING]
> **Do not** run whipcode without a reverse proxy or API gateway in front of it. While there is a standalone mode for per IP rate limiting, it is not meant to be used in production. Use an API gateway like Kong, Tyk and WSO2 to enforce rate limits, policies and authentication. Configure your gateway to add a `X-Master-Key` header to every request with the secret defined below. **Do not** host the gateway on the same system. **Do not** run whipcode as root, or with SELinux disabled/permissive.

1. Save your master key's argon2 hash to *.masterkey*:  `task key`

2. Copy the configuration template:  `task config-init`

3. Run the service:
   ```bash
   # Default port 8000:
   task run

   # Port 6060 with /ping enabled:
   task run -- --ping --port 6060
   ```
   The endpoint will be available at `/run`

4. Test the service:  `task test`\
   If every response has "Success!" in the `stdout` field, the service is working correctly.

## Systemd
Install and enable the systemd user service:  `task systemd-install`

View status or logs:
```bash
task status
task logs        # only logs from whipcode
task logs-full   # logs including podman
```

## CLI options
> [!NOTE]
> The default values are not hardcoded, but specified in the [configuration file](/config.default.toml).

- `-a` `--addr` `ADDR`\
  The address to listen on. (default: none [listen on all interfaces])

- `-p` `--port` `PORT`\
  The port to listen on. May not always work with authbind when attempting to bind to ports < 1024. (default: 8000)

- `-b` `--max-bytes` `BYTES`\
  The maximum size of the request body in bytes. Requests larger than this will be rejected. (default: 1000000)

- `-t` `--timeout` `SECONDS`\
  The maximum time allowed for code execution. (default: 10)

- `-k` `--key` `FILE`\
  Path to the file containing the master key's argon2 hash and salt. (default: .masterkey)

- `-m` `--lang-map` `FILE`\
  Path to the file containing the language map. (default: langmap.toml)

- `--podman-path` `PATH`\
  Path to the podman binary. (default: /usr/bin/podman)

- `--proxy` `ADDR`\
  The address of the reverse proxy or API gateway in front of whipcode. Requests not originating from this address will be rejected. (default: none)

- `--cache`\
  Enables an LRU cache for code executions. This will speed up responses for repeated requests. (default: false)\
  **Note:** The cache is not persistent and will be lost on restart. While this feature is intended to reduce server load and latency, in some situations it may end up worsening it. Memory usage will also increase.

- `--tls`\
  Enables TLS.

- `--tls-dir` `DIR`\
  The directory containing cert.pem and key.pem. (default: tls)

- `--ping`\
  Enables the /ping endpoint. Replies with "pong".

- `--standalone`\
  Enables per IP rate limiting, without the need for a reverse proxy or API gateway. This is NOT RECOMMENDED in production. (default: false)

- `--burst` `COUNT`     (Requires --standalone)\
  The number of requests allowed in a burst. (default: 3)

- `--refill` `SECONDS`  (Requires --standalone)\
  The number of seconds for each request to refill in the burst bucket. (default: 1)

## API reference

`POST /run`

### Headers
- `Content-Type: application/json`
- `X-Master-Key: $MASTER_KEY`

### Body
| Name          | Required | Type                 | Description                                    |
| ------------- | -------- | -------------------- | ---------------------------------------------- |
| `code`        | yes      | `string`             | The source code, base64 encoded.               |
| `language_id` | yes      | `integer` `string`   | Language ID of the submitted code.             |
| `args`        | no       | `string`             | Compiler/interpreter args separated by spaces. |
| `timeout`     | no       | `integer` `string`   | Timeout in seconds for the code to run. Capped at the timeout set in whipcode's configuration. |
| `stdin`       | no       | `string`             | Standard input.                                |
| `env`         | no       | `object`             | Environment variables.                         |

### Response
`200 OK`
| Name            | Type     | Description                                                     |
| --------------- | -------- | --------------------------------------------------------------- |
| `stdout`        | `string` | All data captured from stdout.                                  |
| `stderr`        | `string` | All data captured from stderr.                                  |
| `container_age` | `float`  | Duration the container allocated for your code ran, in seconds. |
| `timeout`       | `bool`   | Boolean value depending on whether your container lived past the timeout period. A reply from a timed-out request will not have any data in stdout and stderr.|

`400` `401` `403` `404` `405` `415` `429` `500`
| Name     | Type     | Description                                       |
| -------- | -------- | ------------------------------------------------- |
| `detail` | `string` | Details about why the request failed to complete. |

### Example request
```bash
lang=2  # javascript
code='console.log("Hello world!");'
timeout=5

curl -s -X POST $ENDPOINT \
    -H 'Content-Type: application/json' \
    -H "X-Master-Key: $MASTER_KEY" \
    -d '{
        "language_id": "'$lang'",
        "code": "'$(echo -n $code | base64)'",
        "timeout": "'$timeout'"
    }' | jq
```


### Example response
```json
{
  "stdout": "Hello world!\n",
  "stderr": "",
  "container_age": 0.335837,
  "timeout": false
}
```

## Tasks
The provided [Taskfile](/Taskfile.yml) has the following tasks defined:
| Task                | Action                                                       |
| ------------------- | ------------------------------------------------------------ |
| `run`               | Run whipcode with the provided with CLI arguments.           |
| `build`             | Build only whipcode.                                         |
| `build-images`      | Build only container images.                                 |
| `rebuild-images`    | Remove existing images and rebuild them.                     |
| `clean`             | Remove all files in the build directory.                     |
| `all`               | Clean and build the project + images.                        |
| `update`            | Pull the latest changes and run the `all` task.              |
| `config-init`       | Copy the default configuration and open it in an editor.     |
| `key`               | Generate a key using `--gen-key`                             |
| `test`              | Run a self-test using `--self-test`                          |
| `systemd-install`   | Install and enable the systemd service for the current user. |
| `status`            | Display the status of the systemd service.                   |
| `logs`              | Show the recent logs for the systemd service using its PID.  |
| `logs-full`         | Show the full logs for the systemd service.                  |


## Contributing
Please read the [Contributing Guidelines](/.github/CONTRIBUTING.md) and [Code of Conduct](/.github/CODE_OF_CONDUCT.md) before opening a pull request.

## Credits
External libraries used:
- [BurntSushi/toml](https://github.com/BurntSushi/toml)
- [karlseguin/ccache](https://github.com/karlseguin/ccache)
- [charmbracelet/log](https://github.com/charmbracelet/log)
- [charmbracelet/huh](https://github.com/charmbracelet/huh)
- [fatih/color](https://github.com/fatih/color)

## License
This project is licensed under the [Apache License, Version 2.0](/LICENSE).
