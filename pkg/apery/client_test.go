package apery

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestEvaluate(t *testing.T) {
	if err := os.Chdir("../.."); err != nil {
		t.Fatal(err)
	}

	cli := NewAperyClient("apery")

	var tests = []struct {
		name           string
		expectValue    int
		expectBestmove string
		expectErr      error
		given          string
	}{
		{"", 87, "2g2f", nil, "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL w - 1"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
			defer cancel()

			value, bestmove, _, err := cli.Evaluate(ctx, tt.given, []string{"7g7f", "3c3d", "2g2f"}, 1*time.Second)

			if value != tt.expectValue {
				t.Errorf("(%s): expected %d, actual %d", tt.given, tt.expectValue, value)
			}
			if bestmove != tt.expectBestmove {
				t.Errorf("(%s): expected %s, actual %s", tt.given, tt.expectBestmove, bestmove)
			}
			if err != tt.expectErr {
				t.Errorf("(%s): expected %v, actual %v", tt.given, tt.expectErr, err)
			}
		})
	}
}
