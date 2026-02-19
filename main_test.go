package main

import (
	"testing"
	"time"
)

func setTestLocale() {
	loc = GPhotosLocale{
		ShortMonthNames: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	}
}

func TestImageIDFromURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{name: "photo url", url: "https://photos.google.com/photo/ABC123", want: "ABC123"},
		{name: "album photo url", url: "https://photos.google.com/album/XYZ/photo/DEF456?foo=bar", want: "DEF456"},
		{name: "missing photo segment", url: "https://photos.google.com/album/XYZ", wantErr: true},
		{name: "malformed url", url: "://bad", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := imageIdFromUrl(tc.url)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got id %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	setTestLocale()
	nowYear := time.Now().Year()

	tests := []struct {
		name     string
		dateStr  string
		timeStr  string
		tzStr    string
		want     time.Time
		wantYear int
		wantErr  bool
	}{
		{
			name:    "full timestamp with timezone",
			dateStr: "Feb 12, 2025",
			timeStr: "6:34 PM",
			tzStr:   "GMT-05:00",
			want:    time.Date(2025, time.February, 12, 18, 34, 0, 0, time.FixedZone("", -5*60*60)),
		},
		{
			name:     "defaults year when missing",
			dateStr:  "Feb 12",
			timeStr:  "1:05 AM",
			tzStr:    "",
			wantYear: nowYear,
		},
		{
			name:    "invalid month",
			dateStr: "Foo 12, 2025",
			timeStr: "1:05 AM",
			wantErr: true,
		},
		{
			name:    "invalid timezone",
			dateStr: "Feb 12, 2025",
			timeStr: "1:05 AM",
			tzStr:   "UTC-5",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseDate(tc.dateStr, tc.timeStr, tc.tzStr)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.want.IsZero() && !got.Equal(tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
			if tc.wantYear != 0 && got.Year() != tc.wantYear {
				t.Fatalf("got year %d, want %d", got.Year(), tc.wantYear)
			}
		})
	}
}

func TestCompareMangled(t *testing.T) {
	tests := []struct {
		name string
		s1   string
		s2   string
		want bool
	}{
		{name: "exact match", s1: "IMG_1234.JPG", s2: "IMG_1234.JPG", want: true},
		{name: "underscore wildcard", s1: "IMG-1234.JPG", s2: "IMG_1234.JPG", want: true},
		{name: "url decoded source", s1: "Vacation%20Photo.jpg", s2: "Vacation Photo.jpg", want: true},
		{name: "different names", s1: "a.jpg", s2: "b.jpg", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := compareMangled(tc.s1, tc.s2)
			if got != tc.want {
				t.Fatalf("compareMangled(%q, %q) = %v, want %v", tc.s1, tc.s2, got, tc.want)
			}
		})
	}
}
