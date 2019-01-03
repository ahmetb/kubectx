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
  cp "$BATS_TEST_DIRNAME/testdata/$1" $KUBECONFIG
}

# wrappers around "kubectl config" command

get_namespace() {
  local cur_ctx

  cur_ctx="$(get_context)" || exit_err "error getting current context"
  ns="$(kubectl config view -o=jsonpath="{.contexts[?(@.name==\"${cur_ctx}\")].context.namespace}")" \
     || exit_err "error getting current namespace"

  if [[ -z "${ns}" ]]; then
    echo "default"
  else
    echo "${ns}"
  fi
}

get_context() {
  kubectl config view -o=jsonpath='{.current-context}'
}

exit_err() {
  echo >&2 "${1}"
  exit 1
}

switch_context() {
  kubectl config use-context "${1}"
}
