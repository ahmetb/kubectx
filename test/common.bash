#!/usr/bin/env bats

# bats setup function
setup() {
  TEMP_HOME="$(mktemp -d)"
  export TEMP_HOME
  export HOME=$TEMP_HOME
  export KUBECONFIG="${TEMP_HOME}/config"
}

# bats teardown function
teardown() {
  rm -rf "$TEMP_HOME"
}

use_config() {
  cp "$BATS_TEST_DIRNAME/testdata/$1" $KUBECONFIG
}

# wrappers around "kubectl config" command

get_namespace() {
  kubectl config view -o=jsonpath="{.contexts[?(@.name==\"$(get_context)\")].context.namespace}"
}

get_context() {
  kubectl config current-context
}

switch_context() {
  kubectl config use-context "${1}"
}
