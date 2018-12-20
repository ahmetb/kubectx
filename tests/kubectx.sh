#!/usr/bin/env bats

clean() {
  rm -f /tmp/config
  rm -f $HOME/.kube/kubectx
}

add_cluster() {
  kubectl config --kubeconfig=/tmp/config set-cluster ${1}
}

add_user() {
  kubectl config --kubeconfig=/tmp/config set-credentials ${1}
}

add_context() {
    kubectl config --kubeconfig=/tmp/config set-context ${1} --user=${2} --cluster=${3}

}

@test "switch to no previous context" {
  clean
  add_cluster cluster1
  add_user user1
  add_context user1@cluster1 user1 cluster1

  KUBECONFIG=/tmp/config run ../kubectx -
  [ "$status" -eq 1 ]
  [ "$output" = "error: No previous context found." ]
}

@test "list one context" {
  clean
  add_cluster cluster1
  add_user user1  
  add_context user1@cluster1 user1 cluster1

  result="$(KUBECONFIG=/tmp/config ../kubectx)"
  [ "$result" = "user1@cluster1" ]
}

@test "list several contexts" {
  clean
  add_cluster cluster1
  add_user user1
  add_user user2
  add_context user1@cluster1 user1 cluster1
  add_context user2@cluster1 user2 cluster1

  result="$(KUBECONFIG=/tmp/config ../kubectx)"
  [ "$result" = "\
user1@cluster1
user2@cluster1" ]
}

@test "switch to context" {
  clean
  add_cluster cluster1
  add_user user1
  add_user user2
  add_context user1@cluster1 user1 cluster1
  add_context user2@cluster1 user2 cluster1

  KUBECONFIG=/tmp/config ../kubectx user1@cluster1
  result="$(kubectl config --kubeconfig=/tmp/config current-context)"
  [ "$result" = "user1@cluster1" ]

  KUBECONFIG=/tmp/config ../kubectx user2@cluster1
  result="$(kubectl config --kubeconfig=/tmp/config current-context)"
  [ "$result" = "user2@cluster1" ]
}

@test "switch to the previous context" {
  clean
  add_cluster cluster1
  add_user user1
  add_user user2
  add_context user1@cluster1 user1 cluster1
  add_context user2@cluster1 user2 cluster1

  KUBECONFIG=/tmp/config ../kubectx user1@cluster1
  result="$(kubectl config --kubeconfig=/tmp/config current-context)"
  [ "$result" = "user1@cluster1" ]

  KUBECONFIG=/tmp/config ../kubectx user2@cluster1
  result="$(kubectl config --kubeconfig=/tmp/config current-context)"
  [ "$result" = "user2@cluster1" ]

  KUBECONFIG=/tmp/config ../kubectx -
  result="$(kubectl config --kubeconfig=/tmp/config current-context)"
  [ "$result" = "user1@cluster1" ]

  KUBECONFIG=/tmp/config ../kubectx -
  result="$(kubectl config --kubeconfig=/tmp/config current-context)"
  [ "$result" = "user2@cluster1" ]
}
