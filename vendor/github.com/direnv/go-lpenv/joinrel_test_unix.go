// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package lpenv

import "testing"

func TestJoinRel(t *testing.T) {
	ret := JoinRel("/foo", "/bar")
	if ret != "/bar" {
		t.Errorf("expected %s to be /bar", ret)
	}
}
