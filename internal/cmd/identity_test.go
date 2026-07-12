package cmd

import "testing"

func TestResolveEmail(t *testing.T) {
	cases := map[string]string{
		"99x999":                        "s99x999@ge.osaka-sandai.ac.jp",
		"99X999":                        "s99x999@ge.osaka-sandai.ac.jp",
		"s99x999":                       "s99x999@ge.osaka-sandai.ac.jp",
		"S99X999":                       "s99x999@ge.osaka-sandai.ac.jp",
		" 99x999 ":                      "s99x999@ge.osaka-sandai.ac.jp",
		"foo@example.com":               "foo@example.com",
		"s99x999@ge.osaka-sandai.ac.jp": "s99x999@ge.osaka-sandai.ac.jp",
		"":                              "",
	}
	for in, want := range cases {
		if got := resolveEmail(in); got != want {
			t.Errorf("resolveEmail(%q) = %q, want %q", in, got, want)
		}
	}
}
