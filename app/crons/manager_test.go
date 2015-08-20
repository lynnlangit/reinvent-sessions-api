package crons

import (
	"testing"

	"github.com/supinf/reinvent-sessions-api/app/misc"
)

func TestCrons(t *testing.T) {
	if misc.ZeroOrNil(crons) {
		t.Errorf("Expected %v, but got %v", 1, len(crons))
		return
	}
	actual := len(crons)
	if actual <= 0 {
		t.Errorf("Expected larger than 0 but got %v", actual)
		return
	}
	actual = len(crons[0].Entries())
	if actual <= 0 {
		t.Errorf("Expected larger than 0 but got %v", actual)
		return
	}
}
