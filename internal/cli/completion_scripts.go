package cli

// Shell completion scripts. They shell out to arsenal for dynamic candidates,
// so arsenal must be on PATH for tool-name completion to work.

const bashCompletion = `# arsenal bash completion
_arsenal() {
  local cur prev cmds
  cur="${COMP_WORDS[COMP_CWORD]}"
  cmds="install remove switch outdated upgrade list search info versions run fetch op sync doctor bundle version completion"
  if [ "$COMP_CWORD" -eq 1 ]; then
    COMPREPLY=( $(compgen -W "$cmds" -- "$cur") )
    return
  fi
  case "${COMP_WORDS[1]}" in
    install)
      if [[ "$cur" == -* ]]; then
        COMPREPLY=( $(compgen -W "--select --github-select --github-ref" -- "$cur") )
      else
        COMPREPLY=( $(compgen -W "$(arsenal search "$cur" 2>/dev/null | grep -v '(asset)$' | awk '{print $1}')" -- "$cur") )
      fi ;;
    versions)
      if [[ "$cur" == -* ]]; then
        COMPREPLY=( $(compgen -W "--github" -- "$cur") )
      else
        COMPREPLY=( $(compgen -W "$(arsenal search "$cur" 2>/dev/null | grep -v '(asset)$' | awk '{print $1}')" -- "$cur") )
      fi ;;
    info)
      COMPREPLY=( $(compgen -W "$(arsenal search "$cur" 2>/dev/null | grep -v '(asset)$' | awk '{print $1}')" -- "$cur") ) ;;
    fetch)
      COMPREPLY=( $(compgen -W "$(arsenal search "$cur" 2>/dev/null | grep '(asset)$' | awk '{print $1}')" -- "$cur") ) ;;
    switch)
      if [[ "$cur" == -* ]]; then
        COMPREPLY=( $(compgen -W "--select" -- "$cur") )
      else
        COMPREPLY=( $(compgen -W "$(arsenal list 2>/dev/null | sed -n 's/^[^A-Za-z]*\([A-Za-z0-9._-]*\)@.*/\1/p')" -- "$cur") )
      fi ;;
    run|remove|outdated|upgrade)
      COMPREPLY=( $(compgen -W "$(arsenal list 2>/dev/null | sed -n 's/^[^A-Za-z]*\([A-Za-z0-9._-]*\)@.*/\1/p')" -- "$cur") ) ;;
    sync)
      COMPREPLY=( $(compgen -W "--list-refs --select-ref --ref --repo" -- "$cur") ) ;;
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
  cmds=(install remove switch outdated upgrade list search info versions run fetch op sync doctor bundle version completion)
  if (( CURRENT == 2 )); then
    compadd -- $cmds
    return
  fi
  case ${words[2]} in
    install)
      if [[ ${words[CURRENT]} == -* ]]; then
        compadd -- --select --github-select --github-ref
      else
        compadd -- ${(f)"$(arsenal search ${words[CURRENT]} 2>/dev/null | grep -v '(asset)$' | awk '{print $1}')"}
      fi ;;
    versions)
      if [[ ${words[CURRENT]} == -* ]]; then
        compadd -- --github
      else
        compadd -- ${(f)"$(arsenal search ${words[CURRENT]} 2>/dev/null | grep -v '(asset)$' | awk '{print $1}')"}
      fi ;;
    info)
      compadd -- ${(f)"$(arsenal search ${words[CURRENT]} 2>/dev/null | grep -v '(asset)$' | awk '{print $1}')"} ;;
    fetch)
      compadd -- ${(f)"$(arsenal search ${words[CURRENT]} 2>/dev/null | grep '(asset)$' | awk '{print $1}')"} ;;
    switch)
      if [[ ${words[CURRENT]} == -* ]]; then
        compadd -- --select
      else
        compadd -- ${(f)"$(arsenal list 2>/dev/null | sed -n 's/^[^A-Za-z]*\([A-Za-z0-9._-]*\)@.*/\1/p')"}
      fi ;;
    run|remove|outdated|upgrade)
      compadd -- ${(f)"$(arsenal list 2>/dev/null | sed -n 's/^[^A-Za-z]*\([A-Za-z0-9._-]*\)@.*/\1/p')"} ;;
    sync)
      compadd -- --list-refs --select-ref --ref --repo ;;
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
complete -c arsenal -n '__fish_use_subcommand' -a 'install remove switch outdated upgrade list search info versions run fetch op sync doctor bundle version completion'
complete -c arsenal -n '__fish_seen_subcommand_from install info versions' -a '(arsenal search 2>/dev/null | grep -v \'(asset)$\' | awk \'{print $1}\')'
complete -c arsenal -n '__fish_seen_subcommand_from install' -l select -d 'Select a registry tool or curated version'
complete -c arsenal -n '__fish_seen_subcommand_from install' -l github-select -d 'Select an upstream GitHub tag'
complete -c arsenal -n '__fish_seen_subcommand_from install' -l github-ref -r -d 'Install a GitHub tag, branch, or commit'
complete -c arsenal -n '__fish_seen_subcommand_from versions' -l github -d 'List upstream GitHub tags'
complete -c arsenal -n '__fish_seen_subcommand_from fetch' -a '(arsenal search 2>/dev/null | grep \'(asset)$\' | awk \'{print $1}\')'
complete -c arsenal -n '__fish_seen_subcommand_from run remove switch outdated upgrade' -a '(arsenal list 2>/dev/null | sed -n \'s/^[^A-Za-z]*\([A-Za-z0-9._-]*\)@.*/\1/p\')'
complete -c arsenal -n '__fish_seen_subcommand_from switch' -l select -d 'Select an installed version'
complete -c arsenal -n '__fish_seen_subcommand_from sync' -l list-refs -d 'List registry repository tags'
complete -c arsenal -n '__fish_seen_subcommand_from sync' -l select-ref -d 'Select a registry repository tag'
complete -c arsenal -n '__fish_seen_subcommand_from sync' -l ref -r -d 'Pin a tag, branch, or commit'
complete -c arsenal -n '__fish_seen_subcommand_from sync' -l repo -r -d 'Use another owner/repository'
complete -c arsenal -n '__fish_seen_subcommand_from op' -a 'create use pin list export import'
complete -c arsenal -n '__fish_seen_subcommand_from completion' -a 'bash zsh fish'
`
