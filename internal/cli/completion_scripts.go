package cli

// Shell completion scripts. They shell out to arsenal for dynamic candidates,
// so arsenal must be on PATH for tool-name completion to work.

const bashCompletion = `# arsenal bash completion
_arsenal() {
  local cur prev cmds
  cur="${COMP_WORDS[COMP_CWORD]}"
  cmds="install remove switch list search info run op sync doctor bundle version completion"
  if [ "$COMP_CWORD" -eq 1 ]; then
    COMPREPLY=( $(compgen -W "$cmds" -- "$cur") )
    return
  fi
  case "${COMP_WORDS[1]}" in
    install|info)
      COMPREPLY=( $(compgen -W "$(arsenal search "$cur" 2>/dev/null | awk '{print $1}')" -- "$cur") ) ;;
    run|remove|switch)
      COMPREPLY=( $(compgen -W "$(arsenal list 2>/dev/null | sed -n 's/^[^A-Za-z]*\([A-Za-z0-9._-]*\)@.*/\1/p')" -- "$cur") ) ;;
    op)
      COMPREPLY=( $(compgen -W "create use pin list export import" -- "$cur") ) ;;
    completion)
      COMPREPLY=( $(compgen -W "bash zsh fish" -- "$cur") ) ;;
  esac
}
complete -F _arsenal arsenal
`

const zshCompletion = `#compdef arsenal
# arsenal zsh completion
_arsenal() {
  local -a cmds
  cmds=(install remove switch list search info run op sync doctor bundle version completion)
  if (( CURRENT == 2 )); then
    compadd -- $cmds
    return
  fi
  case ${words[2]} in
    install|info)
      compadd -- ${(f)"$(arsenal search ${words[CURRENT]} 2>/dev/null | awk '{print $1}')"} ;;
    run|remove|switch)
      compadd -- ${(f)"$(arsenal list 2>/dev/null | sed -n 's/^[^A-Za-z]*\([A-Za-z0-9._-]*\)@.*/\1/p')"} ;;
    op)
      compadd -- create use pin list export import ;;
    completion)
      compadd -- bash zsh fish ;;
  esac
}
compdef _arsenal arsenal
`

const fishCompletion = `# arsenal fish completion
complete -c arsenal -f
complete -c arsenal -n '__fish_use_subcommand' -a 'install remove switch list search info run op sync doctor bundle version completion'
complete -c arsenal -n '__fish_seen_subcommand_from install info' -a '(arsenal search 2>/dev/null | awk \'{print $1}\')'
complete -c arsenal -n '__fish_seen_subcommand_from run remove switch' -a '(arsenal list 2>/dev/null | sed -n \'s/^[^A-Za-z]*\([A-Za-z0-9._-]*\)@.*/\1/p\')'
complete -c arsenal -n '__fish_seen_subcommand_from op' -a 'create use pin list export import'
complete -c arsenal -n '__fish_seen_subcommand_from completion' -a 'bash zsh fish'
`
