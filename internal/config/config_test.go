package config

import "testing"

func TestNextProfileIndex(t *testing.T) {
	tests := []struct {
		name     string
		profiles []Profile
		want     int
	}{
		{"empty", nil, 1},
		{"sequential", []Profile{{Index: 1}, {Index: 2}}, 3},
		{"gap at start", []Profile{{Index: 2}, {Index: 3}}, 1},
		{"gap in middle", []Profile{{Index: 1}, {Index: 3}}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{Profiles: tt.profiles}
			if got := c.NextProfileIndex(); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFindByID(t *testing.T) {
	c := &Config{Profiles: []Profile{
		{ID: "a", Name: "A"},
		{ID: "b", Name: "B"},
	}}
	if got := c.FindByID("b"); got == nil || got.Name != "B" {
		t.Errorf("FindByID b = %+v", got)
	}
	if got := c.FindByID("missing"); got != nil {
		t.Errorf("FindByID missing should be nil, got %+v", got)
	}
}

func TestNewIDUnique(t *testing.T) {
	seen := make(map[string]bool)
	for range 100 {
		id := NewID()
		if seen[id] {
			t.Fatalf("duplicate id: %s", id)
		}
		seen[id] = true
	}
}
