#!/usr/bin/env bats

COMMAND="${COMMAND:-$BATS_TEST_DIRNAME/../kubens}"

# TODO(ahmetb) remove this after bash implementations are deleted
export KUBECTL="$BATS_TEST_DIRNAME/../test/mock-kubectl"

# short-circuit namespace querying in kubens go implementation
export _MOCK_NAMESPACES=1

load common

@test "--help should not fail" {
  run ${COMMAND} --help
  echo "$output">&2
  [[ "$status" -eq 0 ]]
}

@test "-h should not fail" {
  run ${COMMAND} -h
  echo "$output">&2
  [[ "$status" -eq 0 ]]
}

@test "list namespaces when no kubeconfig exists" {
  run ${COMMAND}
  echo "$output"
  [[ "$status" -eq 1 ]]
}

@test "list namespaces" {
  use_config config1
  switch_context user1@cluster1

  run ${COMMAND}
  echo "$output"
  [[ "$status" -eq 0 ]]
  [[ "$output" = *"ns1"* ]]
  [[ "$output" = *"ns2"* ]]
}

@test "switch to existing namespace" {
  use_config config1
  switch_context user1@cluster1

  run ${COMMAND} "ns1"
  echo "$output"
  [[ "$status" -eq 0 ]]
  [[ "$output" = *'Active namespace is "ns1"'* ]]
}

@test "switch to non-existing namespace" {
  use_config config1
  switch_context user1@cluster1

  run ${COMMAND} "unknown-namespace"
  echo "$output"
  [[ "$status" -eq 1 ]]
  [[ "$output" = *'no namespace exists with name "unknown-namespace"'* ]]
}

@test "switch between namespaces" {
  use_config config1
  switch_context user1@cluster1

  run ${COMMAND} ns1
  echo "$output"
  [[ "$status" -eq 0 ]]
  echo "$(get_namespace)"
  [[ "$(get_namespace)" = "ns1" ]]

  run ${COMMAND} ns2
  echo "$output"
  [[ "$status" -eq 0 ]]
  echo "$(get_namespace)"
  [[ "$(get_namespace)" = "ns2" ]]

  run ${COMMAND} -
  echo "$output"
  [[ "$status" -eq 0 ]]
  echo "$(get_namespace)"
  [[ "$(get_namespace)" = "ns1" ]]

  run ${COMMAND} -
  echo "$output"
  [[ "$status" -eq 0 ]]
  echo "$(get_namespace)"
  [[ "$(get_namespace)" = "ns2" ]]
}

@test "switch to previous namespace when none exists" {
  use_config config1
  switch_context user1@cluster1

  run ${COMMAND} -
  echo "$output"
  [[ "$status" -eq 1 ]]
  [[ "$output" = *"No previous namespace found for current context"* ]]
}

@test "switch to namespace when current context is empty" {
  use_config config1

  run ${COMMAND} -
  echo "$output"
  [[ "$status" -eq 1 ]]
  [[ "$output" = *"current-context is not set"* ]]
}

@test "-c/--current works when no namespace is set on context" {
  use_config config1
  switch_context user1@cluster1

  run ${COMMAND} "-c"
  echo "$output"
  [[ "$status" -eq 0 ]]
  [[ "$output" = "default" ]]
  run ${COMMAND} "--current"
  echo "$output"
  [[ "$status" -eq 0 ]]
  [[ "$output" = "default" ]]
}

@test "-c/--current prints the namespace after it is set" {
  use_config config1
  switch_context user1@cluster1
  ${COMMAND} ns1

  run ${COMMAND} "-c"
  echo "$output"
  [[ "$status" -eq 0 ]]
  [[ "$output" = "ns1" ]]
  run ${COMMAND} "--current"
  echo "$output"
  [[ "$status" -eq 0 ]]
  [[ "$output" = "ns1" ]]
}

@test "-c/--current fails when current context is not set" {
  use_config config1
  run ${COMMAND} -c
  echo "$output"
  [[ "$status" -eq 1 ]]

  run ${COMMAND} --current
  echo "$output"
  [[ "$status" -eq 1 ]]
}
