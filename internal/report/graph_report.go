package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/your-org/envlens/internal/graph"
)

// GraphTextReport renders a topological order or cycle error as plain text.
func GraphTextReport(nodes []graph.Node, order []string, cycleErr error) string {
	var sb strings.Builder
	sb.WriteString("=== Dependency Graph ===\n")
	for _, n := range nodes {
		if len(n.Deps) == 0 {
			fmt.Fprintf(&sb, "  %s (no deps)\n", n.Key)
		} else {
			fmt.Fprintf(&sb, "  %s -> [%s]\n", n.Key, strings.Join(n.Deps, ", "))
		}
	}
	sb.WriteString("\n")
	if cycleErr != nil {
		fmt.Fprintf(&sb, "ERROR: %s\n", cycleErr.Error())
		return sb.String()
	}
	sb.WriteString("=== Topological Order ===\n")
	for i, k := range order {
		fmt.Fprintf(&sb, "  %d. %s\n", i+1, k)
	}
	fmt.Fprintf(&sb, "\nTotal: %d keys\n", len(order))
	return sb.String()
}

type graphJSON struct {
	Nodes  []graph.Node `json:"nodes"`
	Order  []string     `json:"order,omitempty"`
	Cycle  string       `json:"cycle,omitempty"`
	Total  int          `json:"total"`
	HasCycle bool       `json:"has_cycle"`
}

// GraphJSONReport renders graph info as JSON.
func GraphJSONReport(nodes []graph.Node, order []string, cycleErr error) string {
	out := graphJSON{
		Nodes: nodes,
		Order: order,
		Total: len(order),
	}
	if cycleErr != nil {
		out.HasCycle = true
		out.Cycle = cycleErr.Error()
		out.Total = 0
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}
