for hook in unload load; do
    add_to_hook "$hook" fish "
        source $(shell_escape "$PWD/demo/autocomplete-plugin/hooks.fish")
        _complete_fish_$hook
    "

    add_to_hook "$hook" bash "
        echo This is the $hook hook for bash. Currently, autocomplete is unimplemented.
    "
done