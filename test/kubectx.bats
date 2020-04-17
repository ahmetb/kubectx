#!/usr/bin/env bats

COMMAND="${COMMAND:-$BATS_TEST_DIRNAME/../kubectx}"

load common

@test "--help should not fail" {
  run ${COMMAND} --help
  echo "$output"
  [ "$status" -eq 0 ]
}

@test "-h should not fail" {
  run ${COMMAND} -h
  echo "$output"
  [ "$status" -eq 0 ]
}

@test "switch to previous context when no one exists" {
  use_config config1

  run ${COMMAND} -
  echo "$output"
  [ "$status" -eq 1 ]
  [[ $output = *"no previous context found" ]]
}

@test "list contexts when no kubeconfig exists" {
  run ${COMMAND}
  echo "$output"
  [ "$status" -eq 0 ]
  [[ "$output" = "warning: kubeconfig file not found" ]]
}

@test "get one context and list contexts" {
  use_config config1

  run ${COMMAND}
  echo "$output"
  [ "$status" -eq 0 ]
  [[ "$output" = "user1@cluster1" ]]
}

@test "get two contexts and list contexts" {
  use_config config2

  run ${COMMAND}
  echo "$output"
  [ "$status" -eq 0 ]
  [[ "$output" = *"user1@cluster1"* ]]
  [[ "$output" = *"user2@cluster1"* ]]
}

@test "get two contexts and select contexts" {
  use_config config2

  run ${COMMAND} user1@cluster1
  echo "$output"
  [ "$status" -eq 0 ]
  echo "$(get_context)"
  [[ "$(get_context)" = "user1@cluster1" ]]

  run ${COMMAND} user2@cluster1
  echo "$output"
  [ "$status" -eq 0 ]
  echo "$(get_context)"
  [[ "$(get_context)" = "user2@cluster1" ]]
}

@test "get two contexts and switch between contexts" {
  use_config config2

  run ${COMMAND} user1@cluster1
  echo "$output"
  [ "$status" -eq 0 ]
  echo "$(get_context)"
  [[ "$(get_context)" = "user1@cluster1" ]]

  run ${COMMAND} user2@cluster1
  echo "$output"
  [ "$status" -eq 0 ]
  echo "$(get_context)"
  [[ "$(get_context)" = "user2@cluster1" ]]

  run ${COMMAND} -
  echo "$output"
  [ "$status" -eq 0 ]
  echo "$(get_context)"
  [[ "$(get_context)" = "user1@cluster1" ]]

  run ${COMMAND} -
  echo "$output"
  [ "$status" -eq 0 ]
  echo "$(get_context)"
  [[ "$(get_context)" = "user2@cluster1" ]]
}

@test "get one context and switch to non existent context" {
  use_config config1

  run ${COMMAND} "unknown-context"
  echo "$output"
  [ "$status" -eq 1 ]
}

@test "-c/--current fails when no context set" {
  use_config config1

  run "${COMMAND}" -c
  echo "$output"
  [ $status -eq 1 ]
  run "${COMMAND}" --current
  echo "$output"
  [ $status -eq 1 ]
}

@test "-c/--current prints the current context" {
  use_config config1

  run "${COMMAND}" user1@cluster1
  [ $status -eq 0 ]

  run "${COMMAND}" -c
  echo "$output"
  [ $status -eq 0 ]
  [[ "$output" = "user1@cluster1" ]]
  run "${COMMAND}" --current
  echo "$output"
  [ $status -eq 0 ]
  [[ "$output" = "user1@cluster1" ]]
}

@test "rename context" {
  use_config config2

  run ${COMMAND} "new-context=user1@cluster1"
  echo "$output"
  [ "$status" -eq 0 ]

  run ${COMMAND}
  echo "$output"
  [ "$status" -eq 0 ]
  [[ ! "$output" = *"user1@cluster1"* ]]
  [[ "$output" = *"new-context"* ]]
  [[ "$output" = *"user2@cluster1"* ]]
}

@test "rename current context" {
  use_config config2

  run ${COMMAND} user2@cluster1
  echo "$output"
  [ "$status" -eq 0 ]

  run ${COMMAND} new-context=.
  echo "$output"
  [ "$status" -eq 0 ]

  run ${COMMAND}
  echo "$output"
  [ "$status" -eq 0 ]
  [[ ! "$output" = *"user2@cluster1"* ]]
  [[ "$output" = *"user1@cluster1"* ]]
  [[ "$output" = *"new-context"* ]]
}

@test "delete context" {
  use_config config2

  run ${COMMAND} -d "user1@cluster1"
  echo "$output"
  [ "$status" -eq 0 ]

  run ${COMMAND}
  echo "$output"
  [ "$status" -eq 0 ]
  [[ ! "$output" = "user1@cluster1" ]]
  [[ "$output" = "user2@cluster1" ]]
}

@test "delete current context" {
  use_config config2

  run ${COMMAND} user2@cluster1
  echo "$output"
  [ "$status" -eq 0 ]

  run ${COMMAND} -d .
  echo "$output"
  [ "$status" -eq 0 ]

  run ${COMMAND}
  echo "$output"
  [ "$status" -eq 0 ]
  [[ ! "$output" = "user2@cluster1" ]]
  [[ "$output" = "user1@cluster1" ]]
}

@test "delete non existent context" {
  use_config config1

  run ${COMMAND} -d "unknown-context"
  echo "$output"
  [ "$status" -eq 1 ]
}

@test "delete several contexts" {
  use_config config2

  run ${COMMAND} -d "user1@cluster1" "user2@cluster1"
  echo "$output"
  [ "$status" -eq 0 ]

  run ${COMMAND}
  echo "$output"
  [ "$status" -eq 0 ]
  [[ "$output" = "" ]]
}

@test "delete several contexts including a non existent one" {
  use_config config2

  run ${COMMAND} -d "user1@cluster1" "non-existent" "user2@cluster1"
  echo "$output"
  [ "$status" -eq 1 ]

  run ${COMMAND}
  echo "$output"
  [ "$status" -eq 0 ]
  [[ "$output" = "user2@cluster1" ]]
}

@test "unset selected context" {
  use_config config2

  run ${COMMAND} user1@cluster1
  [ "$status" -eq 0 ]

  run ${COMMAND} -u
  [ "$status" -eq 0 ]

  run ${COMMAND} -c
  [ "$status" -ne 0 ]
}
