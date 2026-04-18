package graph

import "sort"

// Node represents an env key with its dependencies.
type Node struct {
	Key  string
	Deps []string
}

// CycleError is returned when a circular dependency is detected.
type CycleError struct {
	Cycle []string
}

func (e *CycleError) Error() string {
	result := "cycle detected: "
	for i, k := range e.Cycle {
		if i > 0 {
			result += " -> "
		}
		result += k
	}
	return result
}

// Build constructs a dependency graph from a map of key -> referenced keys.
func Build(deps map[string][]string) []Node {
	nodes := make([]Node, 0, len(deps))
	for k, d := range deps {
		nodes = append(nodes, Node{Key: k, Deps: d})
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Key < nodes[j].Key })
	return nodes
}

// TopoSort returns keys in topological order or an error if a cycle exists.
func TopoSort(deps map[string][]string) ([]string, error) {
	visited := map[string]bool{}
	onStack := map[string]bool{}
	result := []string{}

	var dfs func(key string, path []string) error
	dfs = func(key string, path []string) error {
		if onStack[key] {
			for i, k := range path {
				if k == key {
					return &CycleError{Cycle: append(path[i:], key)}
				}
			}
			return &CycleError{Cycle: append(path, key)}
		}
		if visited[key] {
			return nil
		}
		onStack[key] = true
		for _, dep := range deps[key] {
			if err := dfs(dep, append(path, key)); err != nil {
				return err
			}
		}
		onStack[key] = false
		visited[key] = true
		result = append(result, key)
		return nil
	}

	keys := make([]string, 0, len(deps))
	for k := range deps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if err := dfs(k, []string{}); err != nil {
			return nil, err
		}
	}
	return result, nil
}
