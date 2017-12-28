package basicfunc

import "testing"

func Test(t *testing.T) {

	var tests = []struct {
		s, want string
	}{
		{"Backward", "drawkcaB"},
		{"Hello, АЛЁНКА", "АКНЁЛА ,olleH"},
		{"", ""},
	}

	for _, c := range tests {
		got := Reverse(c.s)
		if got != c.want {
			t.Errorf("basicfunc.Reverse(%q) == %q, want %q", c.s, got, c.want)
		}
	}

	for _, c := range tests {
		got := ReverseMethodSecond(c.s)
		if got != c.want {
			t.Errorf("basicfunc.ReverseMethodSecond(%q) == %q, want %q", c.s, got, c.want)
		}
	}

	var testsString = []struct {
		s1, s2 string
		want   bool
	}{
		{"Backward", "drawkcaB", false},
		{"Hello, АЛЁНКА", "АКНЁЛА ,olleH", false},
		{"", "", true},
		{"abc", "abc", true},
		{"сыр", "сыр", true},
		{"абсде", "абсде", true},
	}

	for _, c := range testsString {
		got := EqualsByChar(c.s1, c.s2)
		if got != c.want {
			t.Errorf("basicfunc.EqualsByChar(%q %q) == %q, want %q", c.s1, c.s2, got, c.want)
		}
	}
}
