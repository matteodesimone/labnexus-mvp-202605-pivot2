package tokens_test

import (
	"testing"

	"github.com/labnexus/labnexus/internal/tokens"
)

func TestEstimate_CharCountDiv4(t *testing.T) {
	cases := []struct {
		name string
		in   int
		want int
	}{
		{"empty", 0, 0},
		{"100 chars → 25 tokens", 100, 25},
		{"1234 chars → 308 tokens", 1234, 308},
		{"1 char → 0 tokens (integer division)", 1, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tokens.Estimate(tc.in)
			if got != tc.want {
				t.Fatalf("Estimate(%d) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

func TestCheck_ThresholdsFR6(t *testing.T) {
	// Reminder: Estimate(textLen) = textLen / 4, ratio = tokens / context_window.
	// 70% = ratio ≥ 0.70 → Warning; 100% = ratio ≥ 1.00 → Block.
	cases := []struct {
		name      string
		textLen   int
		context   int
		wantWarn  bool
		wantBlock bool
	}{
		{"12.5% no flags", 64000, 128000, false, false},   // tokens=16000  ratio=0.125
		{"70% boundary → warn", 358400, 128000, true, false}, // tokens=89600  ratio=0.70
		{"99% → warn", 506880, 128000, true, false},        // tokens=126720 ratio≈0.99
		{"100% → block", 512000, 128000, true, true},       // tokens=128000 ratio=1.00
		{"over 100% → block", 1024000, 128000, true, true}, // tokens=256000 ratio=2.00
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := tokens.Check(tc.textLen, tc.context)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if r.Warning != tc.wantWarn {
				t.Errorf("Warning=%v, want %v", r.Warning, tc.wantWarn)
			}
			if r.Block != tc.wantBlock {
				t.Errorf("Block=%v, want %v", r.Block, tc.wantBlock)
			}
		})
	}
}
