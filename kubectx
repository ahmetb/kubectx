#!/usr/bin/env bash
#
# kubectx(1) is a utility to manage and switch between kubectl contexts.

# Copyright 2017 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

[[ -n $DEBUG ]] && set -x

set -eou pipefail
IFS=$'\n\t'

SELF_CMD="$0"

KUBECTX="${XDG_CACHE_HOME:-$HOME/.kube}/kubectx"

usage() {
  local SELF
  SELF="kubectx"
  if [[ "$(basename "$0")" == kubectl-* ]]; then # invoked as plugin
    SELF="kubectl ctx"
  fi

  cat <<EOF
USAGE:
  $SELF                       : list the contexts
  $SELF <NAME>                : switch to context <NAME>
  $SELF -                     : switch to the previous context
  $SELF -c, --current         : show the current context name
  $SELF <NEW_NAME>=<NAME>     : rename context <NAME> to <NEW_NAME>
  $SELF <NEW_NAME>=.          : rename current-context to <NEW_NAME>
  $SELF -d <NAME> [<NAME...>] : delete context <NAME> ('.' for current-context)
                                  (this command won't delete the user/cluster entry
                                  that is used by the context)
  $SELF -u, --unset           : unset the current context

  $SELF -h,--help             : show this message
EOF
}

exit_err() {
   echo >&2 "${1}"
   exit 1
}

current_context() {
  $KUBECTL config view -o=jsonpath='{.current-context}'
}

get_contexts() {
  $KUBECTL config get-contexts -o=name | sort -n
}

list_contexts() {
  set -u pipefail
  local cur ctx_list
  cur="$(current_context)" || exit_err "error getting current context"
  ctx_list=$(get_contexts) || exit_err "error getting context list"

  local yellow darkbg normal
  yellow=$(tput setaf 3 || true)
  darkbg=$(tput setab 0 || true)
  normal=$(tput sgr0 || true)

  local cur_ctx_fg cur_ctx_bg
  cur_ctx_fg=${KUBECTX_CURRENT_FGCOLOR:-$yellow}
  cur_ctx_bg=${KUBECTX_CURRENT_BGCOLOR:-$darkbg}

  for c in $ctx_list; do
  if [[ -n "${_KUBECTX_FORCE_COLOR:-}" || \
       -t 1 && -z "${NO_COLOR:-}" ]]; then
    # colored output mode
    if [[ "${c}" = "${cur}" ]]; then
      echo "${cur_ctx_bg}${cur_ctx_fg}${c}${normal}"
    else
      echo "${c}"
    fi
  else
    echo "${c}"
  fi
  done
}

read_context() {
  if [[ -f "${KUBECTX}" ]]; then
    cat "${KUBECTX}"
  fi
}

save_context() {
  local saved
  saved="$(read_context)"

  if [[ "${saved}" != "${1}" ]]; then
    printf %s "${1}" > "${KUBECTX}"
  fi
}

switch_context() {
  $KUBECTL config use-context "${1}"
}

choose_context_interactive() {
  local choice
  choice="$(_KUBECTX_FORCE_COLOR=1 \
    FZF_DEFAULT_COMMAND="${SELF_CMD}" \
    fzf --ansi --no-preview || true)"
  if [[ -z "${choice}" ]]; then
    echo 2>&1 "error: you did not choose any of the options"
    exit 1
  else
    set_context "${choice}"
  fi
}

set_context() {
  local prev
  prev="$(current_context)" || exit_err "error getting current context"

  switch_context "${1}"

  if [[ "${prev}" != "${1}" ]]; then
    save_context "${prev}"
  fi
}

swap_context() {
  local ctx
  ctx="$(read_context)"
  if [[ -z "${ctx}" ]]; then
    echo "error: No previous context found." >&2
    exit 1
  fi
  set_context "${ctx}"
}

context_exists() {
  grep -q ^"${1}"\$ <($KUBECTL config get-contexts -o=name)
}

rename_context() {
  local old_name="${1}"
  local new_name="${2}"

  if [[ "${old_name}" == "." ]]; then
    old_name="$(current_context)"
  fi

  if ! context_exists "${old_name}"; then
    echo "error: Context \"${old_name}\" not found, can't rename it." >&2
    exit 1
  fi

  if context_exists "${new_name}"; then
    echo "Context \"${new_name}\" exists, deleting..." >&2
    $KUBECTL config delete-context "${new_name}" 1>/dev/null 2>&1
  fi

  $KUBECTL config rename-context "${old_name}" "${new_name}"
}

delete_contexts() {
  for i in "${@}"; do
    delete_context "${i}"
  done
}

delete_context() {
  local ctx
  ctx="${1}"
  if [[ "${ctx}" == "." ]]; then
    ctx="$(current_context)" || exit_err "error getting current context"
  fi
  echo "Deleting context \"${ctx}\"..." >&2
  $KUBECTL config delete-context "${ctx}"
}

unset_context() {
  echo "Unsetting current context." >&2
  $KUBECTL config unset current-context
}

main() {
  if [[ -z "${KUBECTL:-}" ]]; then
    if hash kubectl 2>/dev/null; then
      KUBECTL=kubectl
    elif hash kubectl.exe  2>/dev/null; then
      KUBECTL=kubectl.exe
    else
      echo >&2 "kubectl is not installed"
      exit 1
    fi
  fi

  if [[ "$#" -eq 0 ]]; then
    if [[ -t 1 &&  -z "${KUBECTX_IGNORE_FZF:-}" && "$(type fzf &>/dev/null; echo $?)" -eq 0 ]]; then
      choose_context_interactive
    else
      list_contexts
    fi
  elif [[ "${1}" == "-d" ]]; then
    if [[ "$#" -lt 2 ]]; then
      echo "error: missing context NAME" >&2
      usage
      exit 1
    fi
    delete_contexts "${@:2}"
  elif [[ "$#" -gt 1 ]]; then
    echo "error: too many arguments" >&2
    usage
    exit 1
  elif [[ "$#" -eq 1 ]]; then
    if [[ "${1}" == "-" ]]; then
      swap_context
    elif [[ "${1}" == '-c' || "${1}" == '--current' ]]; then
      # we don't call current_context here for two reasons:
      # - it does not fail when current-context property is not set
      # - it does not return a trailing newline
      kubectl config current-context
    elif [[ "${1}" == '-u' || "${1}" == '--unset' ]]; then
      unset_context
    elif [[ "${1}" == '-h' || "${1}" == '--help' ]]; then
      usage
    elif [[ "${1}" =~ ^-(.*) ]]; then
      echo "error: unrecognized flag \"${1}\"" >&2
      usage
      exit 1
    elif [[ "${1}" =~ (.+)=(.+) ]]; then
      rename_context "${BASH_REMATCH[2]}" "${BASH_REMATCH[1]}"
    else
      set_context "${1}"
    fi
  else
    usage
    exit 1
  fi
}

main "$@"
