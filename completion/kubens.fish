# kubens

function __fish_kubens_arg_number -a number
    set -l cmd (commandline -opc)
    test (count $cmd) -eq $number
end

complete -f -c kubens
complete -f -x -c kubens -n '__fish_kubens_arg_number 1' -a "(kubectl get ns -o=custom-columns=NAME:.metadata.name --no-headers)"
complete -f -x -c kubens -n '__fish_kubens_arg_number 1' -a "-" -d "switch to the previous namespace in this context"
complete -f -x -c kubens -n '__fish_kubens_arg_number 1' -s c -l current -d "show the current namespace"
complete -f -x -c kubens -n '__fish_kubens_arg_number 1' -s h -l help -d "show the help message"
