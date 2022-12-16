#compdef kubens kns=kubens
_arguments >/dev/null 2>&1 "1: :(- $(kubectl get namespaces -o=jsonpath='{range .items[*].metadata.name}{@}{"\n"}{end}'))"
