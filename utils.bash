current_context() {
  kubectl config view -o=jsonpath='{.current-context}'
}

get_contexts() {
  kubectl config get-contexts -o=name | sort -n
}

get_namespaces() {
  kubectl get namespaces -o=jsonpath='{range .items[*].metadata.name}{@}{"\n"}{end}'
}
