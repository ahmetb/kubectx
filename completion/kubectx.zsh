#compdef kubectx

KUBECTX="${HOME}/.kube/kubectx"
PREV=""
if [ -f "$KUBECTX" ]; then
    # show '-' only if there's a saved previous context
    PREV=$(cat "${KUBECTX}")
    _arguments "1: :((-\:Back\ to\ ${PREV} \
        $(kubectl config get-contexts | awk '{print $2}' | tail -n +2)))"
else
    _arguments "1: :($(kubectl config get-contexts | awk '{print $2}' | tail -n +2))"
fi
