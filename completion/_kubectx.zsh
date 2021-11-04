#compdef kubectx kctx=kubectx

local KUBECTX="${HOME}/.kube/kubectx"
PREV=""

local context_array=("${(@f)$(kubectl config get-contexts --output='name')}")
local all_contexts=(\'${^context_array}\')

if [ -f "$KUBECTX" ]; then
    # show '-' only if there's a saved previous context
    local PREV=$(cat "${KUBECTX}")

    _arguments \
      "-d:*: :(${all_contexts})" \
      "(- *): :(- ${all_contexts})"
else
    _arguments \
      "-d:*: :(${all_contexts})" \
      "(- *): :(${all_contexts})"
fi
