package cmd

import "testing"

func TestEditorFromEnv(t *testing.T) {
	cases := []struct {
		env  Env
		want string
	}{
		{Env{}, ""},
		{Env{"EDITOR": "nano"}, "nano"},
		{Env{"VISUAL": "vim"}, "vim"},
		{Env{"VISUAL": "vim", "EDITOR": "nano"}, "vim"},
		{Env{"VISUAL": "", "EDITOR": "nano"}, "nano"},
	}
	for _, c := range cases {
		if got := editorFromEnv(c.env); got != c.want {
			t.Errorf("editorFromEnv(%v) = %q, want %q", c.env, got, c.want)
		}
	}
}
