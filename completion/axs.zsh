_axs() {
    local arguments
    local targets

    targets=($(axs --list))

    arguments=(
        '(- :)-list[list targets]'
        '(-resolve -list)-resolve[resolve a target]'
        '1: :->target'
    )

    _arguments -S -A "-*" $arguments
    case "$state" in
        target)
                _values 'target' $targets
        ;;
    esac
}

compdef _axs axs
