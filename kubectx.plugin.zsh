#!/usr/bin/env zsh
#
# Enables this repo to be autoloaded by ZSH frameworks that support
# oh-my-zsh compatible plugins.

KUBECTX_DIR="$(dirname $0)"
export PATH=${PATH}:${KUBECTX_DIR}

# Load completions
source "${KUBECTX_DIR}/completion/kubectx.zsh"
source "${KUBECTX_DIR}/completion/kubens.zsh"
