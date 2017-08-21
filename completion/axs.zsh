_axs() {
    local arguments
    local targets

    arguments=(
        '(- :)-list[list targets]'
        '(-resolve -list)-resolve[resolve a target]'
        '1: :->target'
    )

    _arguments -S -A "-*" $arguments
    case "$state" in
        target)
                targets=($(axs -list 2>/dev/null))
                if [[ $? -eq 0 ]]; then
                        _values 'target' $targets
                fi
        ;;
    esac
}

compdef _axs axs
