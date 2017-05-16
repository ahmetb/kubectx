#compdef kubectx

KUBECTX="${HOME}/.kube/kubectx"
PREV=""
if [ -f "$KUBECTX" ]; then
    # show '-' only if there's a saved previous context
    PREV=$(cat "${KUBECTX}")
    _arguments "1: :((-\:Back\ to\ ${PREV} \
        $(kubectl config get-contexts --output='name')))"
else
    _arguments "1: :($(kubectl config get-contexts --output='name'))"
fi
