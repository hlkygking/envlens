package graph

import (
	"testing"
)

func TestBuild_ReturnsNodes(t *testing.T) {
	deps := map[string][]string{
		"A": {"B"},
		"B": {},
	}
	nodes := Build(deps)
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
	if nodes[0].Key != "A" {
		t.Errorf("expected first node A, got %s", nodes[0].Key)
	}
}

func TestTopoSort_NoCycle(t *testing.T) {
	deps := map[string][]string{
		"A": {"B"},
		"B": {"C"},
		"C": {},
	}
	order, err := TopoSort(deps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(order))
	}
	pos := map[string]int{}
	for i, k := range order {
		pos[k] = i
	}
	if pos["C"] > pos["B"] || pos["B"] > pos["A"] {
		t.Errorf("wrong order: %v", order)
	}
}

func TestTopoSort_DetectsCycle(t *testing.T) {
	deps := map[string][]string{
		"A": {"B"},
		"B": {"A"},
	}
	_, err := TopoSort(deps)
	if err == nil {
		t.Fatal("expected cycle error")
	}
	cycleErr, ok := err.(*CycleError)
	if !ok {
		t.Fatalf("expected CycleError, got %T", err)
	}
	if len(cycleErr.Cycle) == 0 {
		t.Error("cycle should contain keys")
	}
}

func TestTopoSort_EmptyGraph(t *testing.T) {
	order, err := TopoSort(map[string][]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 0 {
		t.Errorf("expected empty order, got %v", order)
	}
}

func TestTopoSort_Isolated(t *testing.T) {
	deps := map[string][]string{
		"X": {},
		"Y": {},
	}
	order, err := TopoSort(deps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 2 {
		t.Errorf("expected 2 keys, got %d", len(order))
	}
}
