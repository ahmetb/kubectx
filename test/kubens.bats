#!/usr/bin/env bats

COMMAND="${BATS_TEST_DIRNAME}/../kubens"

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
