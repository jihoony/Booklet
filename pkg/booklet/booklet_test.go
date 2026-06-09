package booklet

import (
	"strings"
	"testing"
)

func TestBuildConfigString(t *testing.T) {
	tests := []struct {
		name       string
		opts       Options
		pageCount  int
		wantContains []string
		wantNotContains []string
	}{
		{
			name: "Multifolio manually enabled",
			opts: Options{
				N:          4,
				FormSize:   "A4",
				Margin:     10.0,
				Binding:    "long",
				BType:      "booklet",
				Multifolio: true,
				FolioSize:  8,
			},
			pageCount: 16,
			wantContains: []string{
				"formsize:A4",
				"margin:10.0",
				"binding:long",
				"btype:booklet",
				"multifolio:on",
				"foliosize:8",
			},
			wantNotContains: []string{},
		},
		{
			name: "Multifolio auto-enabled (sheets > 10)",
			opts: Options{
				N:          4, // sheets per sig threshold test
				FormSize:   "A3",
				Margin:     15.5,
				Binding:    "short",
				BType:      "bookletadvanced",
				Multifolio: false,
				FolioSize:  4,
			},
			// N = 4, pagesPerSheet = 8.
			// totalSheets = (90 + 7) / 8 = 12 sheets.
			// 12 sheets > 10 sheets, so multifolio should auto-enable.
			pageCount: 90,
			wantContains: []string{
				"formsize:A3",
				"margin:15.5",
				"binding:short",
				"btype:bookletadvanced",
				"multifolio:on",
				"foliosize:4",
			},
			wantNotContains: []string{},
		},
		{
			name: "Multifolio not enabled (sheets <= 10)",
			opts: Options{
				N:          4,
				FormSize:   "Letter",
				Margin:     5.0,
				Binding:    "long",
				BType:      "perfectbound",
				Multifolio: false,
				FolioSize:  6,
			},
			// N = 4, pagesPerSheet = 8.
			// totalSheets = (80 + 7) / 8 = 10 sheets.
			// 10 sheets <= 10 sheets, so multifolio should NOT auto-enable.
			pageCount: 80,
			wantContains: []string{
				"formsize:Letter",
				"margin:5.0",
				"binding:long",
				"btype:perfectbound",
			},
			wantNotContains: []string{
				"multifolio:on",
				"foliosize",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildConfigString(tt.opts, tt.pageCount)

			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("buildConfigString() got = %q, want to contain %q", got, want)
				}
			}

			for _, notWant := range tt.wantNotContains {
				if strings.Contains(got, notWant) {
					t.Errorf("buildConfigString() got = %q, want NOT to contain %q", got, notWant)
				}
			}
		})
	}
}
