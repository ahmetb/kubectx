_kube_contexts()
{
  local curr_arg;
  curr_arg=${COMP_WORDS[COMP_CWORD]}
  COMPREPLY=( $(compgen -W "- $(kubectl config get-contexts | awk '{print $2}' | tail -n +2)" -- $curr_arg ) );
}

complete -F _kube_contexts kubectx