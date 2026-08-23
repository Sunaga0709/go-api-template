package location

import (
	"testing"
	"time"
)

func TestJST(t *testing.T) {
	first := JST()
	if first == nil || first.String() != "Asia/Tokyo" {
		t.Fatalf("JST() = %v", first)
	}
	jstLocationMu.Lock()
	jstLocation = nil
	jstLocationMu.Unlock()
	second := JST()
	if second == nil || second.String() != "Asia/Tokyo" {
		t.Errorf("JST() after reset = %v", second)
	}
	if _, offset := time.Now().In(second).Zone(); offset != 9*60*60 {
		t.Errorf("JST offset = %d", offset)
	}
}

func TestJst(t *testing.T) {
	if got := jst(); got == nil || got.String() != "Asia/Tokyo" {
		t.Errorf("jst() = %v", got)
	}
}
