package cmd

import (
	"strings"
	"testing"
)

func TestBashEscape(t *testing.T) {
	assertEqual(t, `''`, BashEscape(""))
	assertEqual(t, `$'escape\'quote'`, BashEscape("escape'quote"))
	assertEqual(t, `$'foo\r\n\tbar'`, BashEscape("foo\r\n\tbar"))
	assertEqual(t, `$'foo bar'`, BashEscape("foo bar"))
	assertEqual(t, `$'\xc3\xa9'`, BashEscape("é"))
}

func TestFishEscape(t *testing.T) {
	assertEqual(t, `''`, (fish{}).escape(""))
	assertEqual(t, `'escape\'quote'`, (fish{}).escape("escape'quote"))
	assertEqual(t, `'foo'\r''\n''\t'bar'`, (fish{}).escape("foo\r\n\tbar"))
	assertEqual(t, `'foo bar'`, (fish{}).escape("foo bar"))
	assertEqual(t, `''\Xc3''\Xa9''`, (fish{}).escape("é"))
}

func TestTcshEscape(t *testing.T) {
	assertEqual(t, `''`, (tcsh{}).escape(""))
	assertEqual(t, `escape\'quote`, (tcsh{}).escape("escape'quote"))
	assertEqual(t, `foo\r\n\tbar`, (tcsh{}).escape("foo\r\n\tbar"))
	assertEqual(t, `foo\ bar`, (tcsh{}).escape("foo bar"))
	assertEqual(t, `\xc3\xa9`, (tcsh{}).escape("é"))
}

func TestPowerShellEscape(t *testing.T) {
	assertEqual(t, "__DiReNv_UnReAcHaBlE__", PowerShellEscapeEnvKey(""))
	assertEqual(t, `a\x2ab`, PowerShellEscapeEnvKey("a*b"))
	assertEqual(t, "a`{b`}", PowerShellEscapeEnvKey("a{b}"))
	assertEqual(t, "__DiReNv_UnReAcHaBlE__", PowerShellEscapeVerbatimEnvKey(""))
	assertEqual(t, "don''t", PowerShellEscapeVerbatimEnvKey("don't"))
	assertEqual(t, "", PowerShellEscapeVerbatimString(""))
	assertEqual(t, "don''t", PowerShellEscapeVerbatimString("don't"))
}

func BenchmarkBashEscape(b *testing.B) {
	var input strings.Builder
	for input.Len() < 60000 {
		input.WriteString("/usr/local/lib/some-package-1.2.3/bin:")
		if input.Len()%7 == 0 {
			input.WriteString("some 'quoted' $value `with` *chars* ")
		}
	}
	str := input.String()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BashEscape(str)
	}
}

func TestShellDetection(t *testing.T) {
	assertNotNil(t, DetectShell("-bash"))
	assertNotNil(t, DetectShell("-/bin/bash"))
	assertNotNil(t, DetectShell("-/usr/local/bin/bash"))
	assertNotNil(t, DetectShell("-zsh"))
	assertNotNil(t, DetectShell("-/bin/zsh"))
	assertNotNil(t, DetectShell("-/usr/local/bin/zsh"))
}

func assertNotNil(t *testing.T, a Shell) {
	if a == nil {
		t.Error("Expected not to be nil")
	}
}

func assertEqual(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected \"%v\" to equal \"%v\"", expected, actual)
	}
}
