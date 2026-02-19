package main

import "testing"

func TestGetAriaLabelSelector(t *testing.T) {
	tests := []struct {
		name  string
		match NodeLabelMatch
		want  string
	}{
		{
			name:  "equals",
			match: NodeLabelMatch{MatchType: "equals", MatchValue: "More options"},
			want:  `[aria-label="More options"]`,
		},
		{
			name:  "startsWith",
			match: NodeLabelMatch{MatchType: "startsWith", MatchValue: "Date taken"},
			want:  `[aria-label^="Date taken"]`,
		},
		{
			name:  "contains",
			match: NodeLabelMatch{MatchType: "contains", MatchValue: "Download"},
			want:  `[aria-label*="Download"]`,
		},
		{
			name:  "endsWith",
			match: NodeLabelMatch{MatchType: "endsWith", MatchValue: "photo"},
			want:  `[aria-label$="photo"]`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := getAriaLabelSelector(tc.match)
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
