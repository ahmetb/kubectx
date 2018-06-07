_kube_contexts()
{
  local curr_arg
  curr_arg=${COMP_WORDS[COMP_CWORD]}
  if [[ "${COMP_WORDS[$(($COMP_CWORD - 1))]}" == "-i" ]]; then
    compopt -o default
    COMPREPLY=()
  else
    compopt +o default
    COMPREPLY=( $(compgen -W "- $(kubectl config get-contexts --output='name')" -- $curr_arg ) )
  fi
}

complete -F _kube_contexts kubectx kctx
