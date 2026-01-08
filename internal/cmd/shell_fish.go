package cmd

import (
	"fmt"
	"strings"
)

type fish struct{}

// Fish adds support for the fish shell as a host
var Fish Shell = fish{}

const fishHook = `
    function __direnv_export_eval --on-event fish_prompt;
        "{{.SelfPath}}" export fish | source;
        __direnv_update_fish_complete_path;

        if test "$direnv_fish_mode" != "disable_arrow";
            function __direnv_cd_hook --on-variable PWD;
                if test "$direnv_fish_mode" = "eval_after_arrow";
                    set -g __direnv_export_again 0;
                else;
                    "{{.SelfPath}}" export fish | source;
                    __direnv_update_fish_complete_path;
                end;
            end;
        end;
    end;

    function __direnv_export_eval_2 --on-event fish_preexec;
        if set -q __direnv_export_again;
            set -e __direnv_export_again;
            "{{.SelfPath}}" export fish | source;
            __direnv_update_fish_complete_path;
            echo;
        end;

        functions --erase __direnv_cd_hook;
    end;

    function __direnv_update_fish_complete_path;
        # Remove previously added completion paths
        for p in $__direnv_fish_complete_paths;
            set -l idx (contains -i -- $p $fish_complete_path);
            and set -e fish_complete_path[$idx];
        end;
        set -e __direnv_fish_complete_paths;

        # Add completion paths from current XDG_DATA_DIRS
        for dir in (string split ':' -- $XDG_DATA_DIRS);
            set -l completions_dir "$dir/fish/vendor_completions.d";
            if test -d "$completions_dir";
                if not contains -- "$completions_dir" $fish_complete_path;
                    set -ga fish_complete_path $completions_dir;
                    set -ga __direnv_fish_complete_paths $completions_dir;
                end;
            end;
        end;
    end;
`

func (sh fish) Hook() (string, error) {
	return fishHook, nil
}

func (sh fish) Export(e ShellExport) (string, error) {
	var out string
	for key, value := range e {
		if value == nil {
			out += sh.unset(key)
		} else {
			out += sh.export(key, *value)
		}
	}
	return out, nil
}

func (sh fish) Dump(env Env) (string, error) {
	var out string
	for key, value := range env {
		out += sh.export(key, value)
	}
	return out, nil
}

func (sh fish) export(key, value string) string {
	if key == "PATH" {
		command := "set -x -g PATH"
		for _, path := range strings.Split(value, ":") {
			command += " " + sh.escape(path)
		}
		return command + ";"
	}
	return "set -x -g " + sh.escape(key) + " " + sh.escape(value) + ";"
}

func (sh fish) unset(key string) string {
	return "set -e -g " + sh.escape(key) + ";"
}

func (sh fish) escape(str string) string {
	in := []byte(str)
	out := "'"
	i := 0
	l := len(in)

	hex := func(char byte) {
		out += fmt.Sprintf("'\\X%02x'", char)
	}

	backslash := func(char byte) {
		out += string([]byte{BACKSLASH, char})
	}

	escaped := func(str string) {
		out += "'" + str + "'"
	}

	literal := func(char byte) {
		out += string([]byte{char})
	}

	for i < l {
		char := in[i]
		switch {
		case char == TAB:
			escaped(`\t`)
		case char == LF:
			escaped(`\n`)
		case char == CR:
			escaped(`\r`)
		case char <= US:
			hex(char)
		case char == SINGLE_QUOTE:
			backslash(char)
		case char == BACKSLASH:
			backslash(char)
		case char <= TILDE:
			literal(char)
		case char == DEL:
			hex(char)
		default:
			hex(char)
		}
		i++
	}

	out += "'"

	return out
}
