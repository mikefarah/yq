
_yq() {
    local cur prev coms opts car cdr x
    coms="compare delete help merge new prefix read validate write"
    cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ ${cur} == -* ]] ; then
        COMPREPLY=()
    else
        case "$COMP_CWORD" in
            "1")
                COMPREPLY=( $(compgen -W "$coms" -- $cur) )
                ;;
            "2")
                if [[ "${COMP_WORDS[1]}" == read ]] ; then
                    [ "$cur" ] && { p="" ; s="*" ; } || { p="./" ; s="" ; }
                    cur1=$(echo "$cur" | sed -e 's!^~!'$HOME'!g')
                    COMPREPLY=( $(x="$p$cur1$s" ; echo "$x") )
                else
                    COMPREPLY=()
                fi
                ;;
            "3")
                if [[ "${COMP_WORDS[1]}" == read ]] ; then
                    partkey="$(echo $cur| sed -e 's/[^.]*$//g' -e 's/\.$//g')"
                    [ "$partkey" ] && { x="\"${partkey}.\" + keys[]" ; } || { x="keys[]" ; }
                    COMPREPLY=( $(yq read ${COMP_WORDS[2]} ${partkey} -j | 
                                    jq "$x" -r 2>&- | 
                                    grep ${cur}. 
                                    )  )
                else
                    COMPREPLY=()
                fi
                ;;
            *)
                COMPREPLY=()
                ;;
        esac
    fi
    return 0
}
complete -o nospace -F _yq yq
