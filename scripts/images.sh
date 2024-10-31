#!/usr/bin/env bash
#
#  Copyright 2024 whipcode.app (AnnikaV9)
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#          http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing,
#  software distributed under the License is distributed on
#  an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
#  either express or implied. See the License for the specific
#  language governing permissions and limitations under the License.
#

if ! command -v podman &> /dev/null
then
    echo "Podman is not installed."
    exit 1
fi

declare -A langs=(
    [BASH]='bash'
    [NODEJS]='nodejs'
    [C]='bash gcc && echo "#!/bin/bash" > /usr/bin/c-run && echo "gcc \$1 -o /tmp/run && ./tmp/run" >> /usr/bin/c-run && chmod +x /usr/bin/c-run'
    [CPP]='bash g++ && echo "#!/bin/bash" > /usr/bin/cpp-run && echo "g++ \$1 -o /tmp/run && ./tmp/run" >> /usr/bin/cpp-run && chmod +x /usr/bin/cpp-run'
    [FORTRAN]='bash gfortran && echo "#!/bin/bash" > /usr/bin/f90-run && echo "gfortran \$1 -o /tmp/run && ./tmp/run" >> /usr/bin/f90-run && chmod +x /usr/bin/f90-run'
    [GO]='bash gcc-go && echo "#!/bin/bash" > /usr/bin/go-run && echo "gccgo \$1 -o /tmp/run && ./tmp/run" >> /usr/bin/go-run && chmod +x /usr/bin/go-run'
    [HASKELL]="bash ghc && echo '#!/bin/bash' > /usr/bin/hs-run && echo 'runghc --ghc-arg=\"-v0\" \$1' >> /usr/bin/hs-run && chmod +x /usr/bin/hs-run"
    [JAVA]='openjdk21'
    [LUA]='lua && ln -s /usr/bin/lua5* /usr/bin/lua'
    [PERL]='perl'
    [PYTHON]='python3'
    [RUBY]='ruby'
    [RUST]='bash rust && echo "#!/bin/bash" > /usr/bin/rs-run && echo "rustc \$1 -o /tmp/run && ./tmp/run" >> /usr/bin/rs-run && chmod +x /usr/bin/rs-run'
    [TYPESCRIPT]='npm && npm i -g @swc/cli @swc/core && echo "#!/bin/bash" > /usr/bin/ts-run && echo "swc -q \$1 -o /tmp/run.js && node /tmp/run.js" >> /usr/bin/ts-run && chmod +x /usr/bin/ts-run'
    [LISP]='sbcl'
    [RACKET]='racket'
)

prefix="FROM docker.io/alpine:latest\nRUN apk update --no-cache && apk upgrade --no-cache && apk add --no-cache libc-dev musl-dev "
suffix="&& apk --purge del apk-tools && rm -rf /var/cache/apk /var/lib/apk /lib/apk /etc/apk /sbin/apk /usr/share/apk /usr/lib/apk /usr/sbin/apk /usr/local/apk /usr/bin/apk /usr/local/bin/apk /usr/local/sbin/apk /usr/local/lib/apk /usr/local/share/apk /usr/local/libexec/apk /usr/local/etc/apk"

trap "exit" INT
declare -i i=0
for lang in "${!langs[@]}"; do
    i+=1
    langs["$lang"]=$(echo -e "${prefix}${langs[$lang]} ${suffix}")
    podman build -t whipcode-${lang,,} - <<< "${langs[$lang]}" | while read line; do echo "[$i/16] [$lang] $line"; done
done
