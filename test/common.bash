#!/usr/bin/env bats

# bats setup function
setup() {
  export KUBECONFIG=$(mktemp)
  export XDG_CACHE_HOME="$(mktemp -d)"
}

# bats teardown function
teardown() {
  rm -f $KUBECONFIG
  rm -f $XDG_CACHE_HOME/kubectx
  rmdir $XDG_CACHE_HOME
}

# wrappers around "kubectl config" command

add_cluster() {
  kubectl config set-cluster ${1}
}

add_user() {
  kubectl config set-credentials ${1}
}

add_context() {
    kubectl config set-context ${1} --user=${2} --cluster=${3}
}

get_context() {
    kubectl config current-context
}
