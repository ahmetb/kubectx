#!/usr/bin/env bats

COMMAND="${BATS_TEST_DIRNAME}/../kubens"
export KUBECTL="$BATS_TEST_DIRNAME/../mock-kubectl"

load common

@test "--help should not fail" {
  run ${COMMAND} --help
  echo "$output">&2
  [ "$status" -eq 0 ]
}

@test "-h should not fail" {
  run ${COMMAND} -h
  echo "$output">&2
  [ "$status" -eq 0 ]
}

@test "list namespaces when no kubeconfig exists" {
  run ${COMMAND}
  echo "$output"
  [ "$status" -eq "1" ]
  [[ "$output" = *"error: current-context is not set"* ]]
}

@test "list namespaces" {
  use_config config1
  switch_context user1@cluster1

  run ${COMMAND}
  echo "$output"
  [ "$status" -eq 0 ]
  [[ "$output" = *"default"* ]]
  [[ "$output" = *"kube-public"* ]]
  [[ "$output" = *"kube-system"* ]]
}

@test "switch to existent namespace" {
  use_config config1
  switch_context user1@cluster1

  run ${COMMAND} "kube-public"
  echo "$output"
  [ "$status" -eq 0 ]
  [[ "$output" = *'Active namespace is "kube-public"'* ]]
}

@test "switch to non existent namespace" {
  use_config config1
  switch_context user1@cluster1

  run ${COMMAND} "unknown-context"
  echo "$output"
  [ "$status" -eq 1 ]
  [[ "$output" = 'error: no namespace exists with name "unknown-context".' ]]
}

@test "switch between namespaces" {
  use_config config1
  switch_context user1@cluster1

  run ${COMMAND} kube-public
  echo "$output"
  [ "$status" -eq 0 ]
  echo "$(get_namespace)"
  [[ "$(get_namespace)" = "kube-public" ]]

  run ${COMMAND} kube-system
  echo "$output"
  [ "$status" -eq 0 ]
  echo "$(get_namespace)"
  [[ "$(get_namespace)" = "kube-system" ]]

  run ${COMMAND} -
  echo "$output"
  [ "$status" -eq 0 ]
  echo "$(get_namespace)"
  [[ "$(get_namespace)" = "kube-public" ]]

  run ${COMMAND} -
  echo "$output"
  [ "$status" -eq 0 ]
  echo "$(get_namespace)"
  [[ "$(get_namespace)" = "kube-system" ]]
}

@test "switch to previous namespace when none exists" {
  use_config config1
  switch_context user1@cluster1

  run ${COMMAND} -
  echo "$output"
  [ "$status" -eq 1 ]
  [[ "$output" = "error: No previous namespace found for current context." ]]
}

@test "switch to namespace when current context is empty" {
  use_config config1

  run ${COMMAND} -
  echo "$output"
  [ "$status" -eq 1 ]
  [[ "$output" = *"error: current-context is not set"* ]]
}
