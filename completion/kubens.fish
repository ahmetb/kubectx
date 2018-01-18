# kubens
complete -f -c kubens -a "(kubectl get ns -o=custom-columns=NAME:.metadata.name --no-headers)"

