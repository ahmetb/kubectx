# kubectx

function __fish_kubectx_arg_number -a number
    set -l cmd (commandline -opc)
    test (count $cmd) -eq $number
end

complete -f -c kubectx
complete -f -x -c kubectx -n '__fish_kubectx_arg_number 1' -a "(kubectl config get-contexts --output='name')"
complete -f -x -c kubectx -n '__fish_kubectx_arg_number 1' -a "-" -d "switch to the previous namespace in this context"
