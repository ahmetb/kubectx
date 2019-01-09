#!/usr/bin/env zsh
#
# Enables this repo to be autoloaded by ZSH frameworks that support
# oh-my-zsh format plugins.

KUBECTX_DIR="$(dirname $0)"
export PATH=${PATH}:${KUBECTX_DIR}
