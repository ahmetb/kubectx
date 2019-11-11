#!/usr/bin/env bash
#
# kubens(1) is a utility to switch between Kubernetes namespaces.

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

KUBENS_DIR="${XDG_CACHE_HOME:-$HOME/.kube}/kubens"

usage() {
  local SELF
  SELF="kubens"
  if [[ "$(basename "$0")" == kubectl-* ]]; then # invoked as plugin
    SELF="kubectl ns"
  fi

  cat <<EOF
USAGE:
  $SELF                    : list the namespaces in the current context
  $SELF <NAME>             : change the active namespace of current context
  $SELF -                  : switch to the previous namespace in this context
  $SELF -c, --current      : show the current namespace
  $SELF -h,--help          : show this message
EOF
}

exit_err() {
   echo >&2 "${1}"
   exit 1
}

current_namespace() {
  local cur_ctx

  cur_ctx="$(current_context)" || exit_err "error getting current context"
  ns="$($KUBECTL config view -o=jsonpath="{.contexts[?(@.name==\"${cur_ctx}\")].context.namespace}")" \
     || exit_err "error getting current namespace"

  if [[ -z "${ns}" ]]; then
    echo "default"
  else
    echo "${ns}"
  fi
}

current_context() {
  $KUBECTL config current-context
}

get_namespaces() {
  $KUBECTL get namespaces -o=jsonpath='{range .items[*].metadata.name}{@}{"\n"}{end}'
}

escape_context_name() {
  echo "${1//\//-}"
}

namespace_file() {
  local ctx

  ctx="$(escape_context_name "${1}")"
  echo "${KUBENS_DIR}/${ctx}"
}

read_namespace() {
  local f
  f="$(namespace_file "${1}")"
  [[ -f "${f}" ]] && cat "${f}"
  return 0
}

save_namespace() {
  mkdir -p "${KUBENS_DIR}"
  local f saved
  f="$(namespace_file "${1}")"
  saved="$(read_namespace "${1}")"

  if [[ "${saved}" != "${2}" ]]; then
    printf %s "${2}" > "${f}"
  fi
}

switch_namespace() {
  local ctx="${1}"
  $KUBECTL config set-context "${ctx}" --namespace="${2}"
  echo "Active namespace is \"${2}\".">&2
}

choose_namespace_interactive() {
  # directly calling kubens via fzf might fail with a cryptic error like
  # "$FZF_DEFAULT_COMMAND failed", so try to see if we can list namespaces
  # locally first
  if [[ -z "$(list_namespaces)" ]]; then
    echo >&2 "error: could not list namespaces (is the cluster accessible?)"
    exit 1
  fi

  local choice
  choice="$(_KUBECTX_FORCE_COLOR=1 \
    FZF_DEFAULT_COMMAND="${SELF_CMD}" \
    fzf --ansi --no-preview || true)"
  if [[ -z "${choice}" ]]; then
    echo 2>&1 "error: you did not choose any of the options"
    exit 1
  else
    set_namespace "${choice}"
  fi
}

set_namespace() {
  local ctx prev
  ctx="$(current_context)" || exit_err "error getting current context"
  prev="$(current_namespace)" || exit_error "error getting current namespace"

  if grep -q ^"${1}"\$ <(get_namespaces); then
    switch_namespace "${ctx}" "${1}"

    if [[ "${prev}" != "${1}" ]]; then
      save_namespace "${ctx}" "${prev}"
    fi
  else
    echo "error: no namespace exists with name \"${1}\".">&2
    exit 1
  fi
}

list_namespaces() {
  local yellow darkbg normal
  yellow=$(tput setaf 3 || true)
  darkbg=$(tput setab 0 || true)
  normal=$(tput sgr0 || true)

  local cur_ctx_fg cur_ctx_bg
  cur_ctx_fg=${KUBECTX_CURRENT_FGCOLOR:-$yellow}
  cur_ctx_bg=${KUBECTX_CURRENT_BGCOLOR:-$darkbg}

  local cur ns_list
  cur="$(current_namespace)" || exit_err "error getting current namespace"
  ns_list=$(get_namespaces) || exit_err "error getting namespace list"

  for c in $ns_list; do
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

swap_namespace() {
  local ctx ns
  ctx="$(current_context)" || exit_err "error getting current context"
  ns="$(read_namespace "${ctx}")"
  if [[ -z "${ns}" ]]; then
    echo "error: No previous namespace found for current context." >&2
    exit 1
  fi
  set_namespace "${ns}"
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
    if [[ -t 1 && -z ${KUBECTX_IGNORE_FZF:-} && "$(type fzf &>/dev/null; echo $?)" -eq 0 ]]; then
      choose_namespace_interactive
    else
      list_namespaces
    fi
  elif [[ "$#" -eq 1 ]]; then
    if [[ "${1}" == '-h' || "${1}" == '--help' ]]; then
      usage
    elif [[ "${1}" == "-" ]]; then
      swap_namespace
    elif [[ "${1}" == '-c' || "${1}" == '--current' ]]; then
      current_namespace
    elif [[ "${1}" =~ ^-(.*) ]]; then
      echo "error: unrecognized flag \"${1}\"" >&2
      usage
      exit 1
    elif [[ "${1}" =~ (.+)=(.+) ]]; then
      alias_context "${BASH_REMATCH[2]}" "${BASH_REMATCH[1]}"
    else
      set_namespace "${1}"
    fi
  else
    echo "error: too many flags" >&2
    usage
    exit 1
  fi
}

main "$@"
