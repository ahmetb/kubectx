#!/usr/bin/env bats

# bats setup function
setup() {
  export XDG_CACHE_HOME="$(mktemp -d)"
  export KUBECONFIG="${XDG_CACHE_HOME}/config"
}

# bats teardown function
teardown() {
  rm -rf "$XDG_CACHE_HOME"
}

use_config() {
  cp "testdata/$1" $KUBECONFIG
}

# wrappers around "kubectl config" command

get_context() {
    kubectl config current-context
}
