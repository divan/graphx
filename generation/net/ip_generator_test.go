package net

import "testing"

func TestIPGenerator(t *testing.T) {
	g := NewIPGenerator("10.0.0.1")
	ip1 := g.NextAddress()
	if ip1 != "10.0.0.2" {
		t.Fatal("Expected other value")
	}
	for i := 0; i < 255; i++ {
		g.NextAddress()
	}
	ip2 := g.NextAddress()
	if ip2 != "10.0.1.2" {
		t.Fatal("Expected other value")
	}
}
