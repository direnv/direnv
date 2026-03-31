# For information about how fish autocomplete works:
# https://github.com/fish-shell/fish-shell/issues/8261
#
# A simpler way to implement this would be to append the autocomplete directories to
# the variable $fish_complete_path when we enter the environment and remove them
# when we leave. However, when we remove a directory from $fish_complete_path,
# fish does not allow autocomplete to be loaded for that command again[1].
# Instead, we keep track of the individual autocomplete entries that get added
# to fish after `eval`ing the autocomplete file and when direnv unloads the
# environment, we remove only the entries that we added.
#
# [1]: https://github.com/fish-shell/fish-shell/issues/8261#:~:text=To%20get%20rid%20of%20the%20lingering%20completions%20we%27d%20need%20to%20allow%20erasing%20completions%20without%20adding%20a%20tombstone

function _complete_fish_unload
    _complete_fish_debug 'In unload'

    if test (count $_complete_fish_files) -gt 0
        _complete_fish_debug 'Removing completions for these files:'\n"$(string join \n $_complete_fish_files)"

        for file_index in (seq (count $_complete_fish_files))
            set -l added_entries (string split --no-empty \n "$_complete_fish_entries[$file_index]")
            set -l file $_complete_fish_files[$file_index]
            set -l command (path change-extension '' (path basename $file))
            set -l current_entries (complete $command)
            for entry in $added_entries
                if set -l entry_index (contains --index -- $entry $current_entries)
                    set --erase current_entries[$entry_index]
                end
            end
            complete --erase $command
            printf %s\n $current_entries | source
        end
    end
    
    set -l our_variables (set --names | string match --regex --groups-only '^(_complete_fish.*)')
    set --erase $our_variables

    set -l our_functions (functions --names --all | string match --regex --groups-only '^(_complete_fish.*)')
    functions --erase $our_functions
end

function _complete_fish_load
    _complete_fish_debug 'In load'

    # The value of XDG_DATA_DIRS before the direnv environment was loaded.
    set -l old_xdg_data_dirs (direnv exec / fish -c 'echo "$XDG_DATA_DIRS"' | string split --no-empty -- ':')

    set -l xdg_files
    for dir in (string split --no-empty -- ':' "$XDG_DATA_DIRS")
        # This was in XDG before direnv was loaded so ignore it
        if contains $dir $old_xdg_data_dirs
            continue
        end

        set -l fish_dir $dir'/fish/vendor_completions.d'
        if not test -d $fish_dir
            continue
        end
        set --append xdg_files $fish_dir/*
    end
    if test (count $xdg_files) -eq 0
        return
    end

    _complete_fish_debug 'Adding completions for these files:'\n"$(string join \n $xdg_files)"
    _complete_fish_add $xdg_files
end

function _complete_fish_debug --argument-names message
    if test -n "$COMPLETE_FISH_DEBUG"
        echo "[completion-sync] $message" >&2
    end
end

function _complete_fish_add
    for file in $argv
        set -l command (path change-extension '' (path basename $file))

        set -l old_entries (complete $command)
        source $file
        set -l new_entries (complete $command)
        set -l added_entries
        for new_entry in $new_entries
            if not contains $new_entry $old_entries
                set --append added_entries $new_entry
            end
        end
        if test (count $added_entries) -eq 0
            _complete_fish_debug "This file didn't add any completion entries, ignoring: $file"
            continue
        end

        set --global --append _complete_fish_files $file

        set -l added_entries_string "$(string join --no-empty \n $added_entries)"
        set --global --append _complete_fish_entries $added_entries_string
    end
end