#!/usr/bin/env bats

COMMAND="$BATS_TEST_DIRNAME/../kubectx"

@test "--help should not fail" {
  run ${COMMAND} --help
  [ "$status" -eq 0 ]
}

@test "-h should not fail" {
  run ${COMMAND} -h
  [ "$status" -eq 0 ]
}
