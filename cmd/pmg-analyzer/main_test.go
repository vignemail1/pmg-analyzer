package main

import (
	"bytes"
	"testing"
)

func TestRun(t *testing.T) {
	cases := []struct {
		name string
		args []string
		code int
		out  string
		err  string
	}{
		{"version", []string{"version"}, 0, "pmg-analyzer 0.1.0-dev\n", ""},
		{"missing", nil, 2, "", "usage: pmg-analyzer version\n"},
		{"unknown", []string{"ingest"}, 2, "", "usage: pmg-analyzer version\n"},
		{"extra", []string{"version", "extra"}, 2, "", "usage: pmg-analyzer version\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(tc.args, &stdout, &stderr)
			if code != tc.code || stdout.String() != tc.out || stderr.String() != tc.err {
				t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
		})
	}
}
