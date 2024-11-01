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
    [C]='gcc'
    [CPP]='g++'
    [FORTRAN]='gfortran'
    [GO]='gcc-go'
    [HASKELL]='ghc'
    [JAVA]='openjdk21'
    [LUA]='lua && ln -s /usr/bin/lua5* /usr/bin/lua'
    [PERL]='perl'
    [PYTHON]='python3'
    [RUBY]='ruby'
    [RUST]='rust'
    [TYPESCRIPT]='npm && npm i -g @swc/cli @swc/core'
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
