package boxprompt

import "testing"

func TestLookup(t *testing.T) {
	p, err := Lookup("STATE")
	if err != nil || p.Name != "state" {
		t.Fatalf("got %+v %v", p, err)
	}
	if _, err := Lookup("nope"); err == nil {
		t.Fatal("expected unknown")
	}
	if len(All) != 4 {
		t.Fatalf("parked set size %d", len(All))
	}
}
