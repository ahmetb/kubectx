#!/usr/bin/env bats

COMMAND="$BATS_TEST_DIRNAME/../kubectx"

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
  add_cluster cluster1
  add_user user1
  add_context user1@cluster1 user1 cluster1

  run ${COMMAND} -
  echo "$output"
  [ "$status" -eq 1 ]
  [[ "$output" = "error: No previous context found." ]]
}

@test "create one context and list contexts" {
  add_cluster cluster1
  add_user user1  
  add_context user1@cluster1 user1 cluster1

  run ${COMMAND}
  echo "$output"
  [ "$status" -eq 0 ]
  [[ "$output" = "user1@cluster1" ]]
}

@test "create two contexts and list contexts" {
  add_cluster cluster1
  add_user user1
  add_user user2
  add_context user1@cluster1 user1 cluster1
  add_context user2@cluster1 user2 cluster1

  run ${COMMAND}
  echo "$output"
  [ "$status" -eq 0 ]
  [[ "$output" = *"user1@cluster1"* ]]
  [[ "$output" = *"user2@cluster1"* ]]
}

@test "create two contexts and select contexts" {
  add_cluster cluster1
  add_user user1
  add_user user2
  add_context user1@cluster1 user1 cluster1
  add_context user2@cluster1 user2 cluster1

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

@test "create two contexts and switch between contexts" {
  add_cluster cluster1
  add_user user1
  add_user user2
  add_context user1@cluster1 user1 cluster1
  add_context user2@cluster1 user2 cluster1

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
