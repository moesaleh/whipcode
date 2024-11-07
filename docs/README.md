<img alt="Whipcode" src="https://github.com/user-attachments/assets/b68d8164-cbbe-45cb-9f66-43618a0e8544"><br/>

REST service for executing untrusted code with Podman.

Implemented as a scalable stateless microservice with no user management or authentication, minimizing damage from potential zero-day breakouts.

[API reference](/docs/reference.md)

<details>
  <summary>Click to see default languages</summary>
  <br/>

| ID  | Language       | RT            |
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
- [scripts/images.sh](/scripts/images.sh)
- [scripts/extra_setup/](/scripts/extra_setup/)
- [entry/](/entry/)
- [languages.toml](/languages.default.toml)
- [scripts/tests.sh](/scripts/tests.sh)

</details>


## Setting up
**Please set this up on a SELinux-enabled system.**

Tested on Fedora Server 41.

### Environment setup
Go:
```bash
# Download Go tarball
wget https://go.dev/dl/go1.<version>.linux-amd64.tar.gz

# Clean up old installations
sudo rm -rf /usr/local/go

# Install
sudo tar -C /usr/local -xzf go1.<version>.linux-amd64.tar.gz
sudo echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile

# Load new PATH
source /etc/profile
```

go-task:
```bash
# Download go-task rpm
wget https://github.com/go-task/task/releases/latest/download/task_linux_amd64.rpm

# Install
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
Use `task <action>` to run predefined build actions:

- **all**\
  *Build everything*

- **build**\
  *Build only the service*

- **build-images**\
  *Build only the container images*

- **rebuild-images**\
  *Clean rebuild images*

- **update**\
  *Update (git pull), build service and images*

## Starting the service

> [!WARNING]
> **Do not** run this service without a reverse proxy or API gateway in front of it. While whipcode does have a standalone mode for per IP rate limiting, it is not meant to be used in production. Use an API gateway like Kong, Tyk and WSO2 to enforce rate limits, policies and authentication. Configure your gateway to add a `X-Master-Key` header to every request with the secret defined below. **Do not** host the gateway on the same system.
>
> **Do not** run whipcode as root, or with SELinux disabled/permissive.

1. Save your master key's argon2 hash to *.masterkey*:
   ```bash
   # We don't want the shell storing the key
   export HISTFILE=/dev/null
   task key -- <KEY>
   ```

2. Copy the configuration templates:  `task config-init`

3. Run the service:
   ```bash
   # Default port 8000:
   task run

   # Port 6060 with /ping enabled:
   task run -- --ping --port 6060
   ```

4. Test the service:  `task test`

   If every response has "Success!" in the `stdout` field, the service is working correctly.

## Systemd
Install the systemd user service:  `task systemd-install`

View status/logs:
```bash
task status
task logs        # only logs from whipcode
task logs-full   # logs including podman
```

## CLI options
CLI options:

- **-p, --port  PORT**\
  *The port to listen on. May not always work with authbind when attempting to bind to ports < 1024. (default: 8000)*

- **-m, --max  BYTES**\
  *The maximum size of the request body in bytes. Requests larger than this will be rejected. (default: 1000000)*

- **-t, --timeout  SECONDS**\
  *The maximum time allowed for code execution. Should be set lower than the server's write timeout, which is 20 seconds. (default: 10)*

- **-k, --key  FILE**\
  *Path to the file containing the master key's argon2 hash and salt. (default: .masterkey)*

- **--proxy  ADDR**\
  *The address of the reverse proxy or API gateway in front of whipcode. Requests not originating from this address will be rejected. (default: none)*

- **--cache**\
  *Enables an LRU cache for code executions. This will speed up responses for repeated requests. (default: false)*

  ***Note:** The cache is not persistent and will be lost on restart. While this feature is intended to reduce server load and latency, in some situations it may end up worsening it. Memory usage will also increase.*

- **--tls**\
  *Enables TLS. Requires tls/cert.pem and tls/key.pem to be present.*

- **--ping**\
  *Enables the /ping endpoint. Replies with "pong".*

- **--standalone**\
  *Enables per IP rate limiting, without the need for a reverse proxy or API gateway. This is NOT RECOMMENDED in production. (default: false)*

- **--burst  COUNT**     (Requires --standalone)\
  *The number of requests allowed in a burst. (default: 3)*

- **--refill  SECONDS**  (Requires --standalone)\
  *The number of seconds for each request to refill in the burst bucket. (default: 1)*

## License
This project is licensed under the[ Apache License, Version 2.0](/LICENSE)
