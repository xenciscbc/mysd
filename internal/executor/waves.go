package executor

import (
	"errors"
	"sort"
)

// ErrCyclicDependency is returned by BuildWaveGroups when a cycle is detected in the task graph.
var ErrCyclicDependency = errors.New("cyclic dependency detected in task graph")

// BuildWaveGroups computes dependency-ordered wave layers from pending tasks.
// Uses Kahn's algorithm (BFS topological sort) to produce layers,
// then splits any same-layer tasks with file overlap into separate waves.
// Returns ErrCyclicDependency if the task graph contains a cycle.
func BuildWaveGroups(tasks []TaskItem) ([][]TaskItem, error) {
	if len(tasks) == 0 {
		return nil, nil
	}

	// Build index and in-degree map
	idToTask := make(map[int]TaskItem, len(tasks))
	inDegree := make(map[int]int, len(tasks))
	adj := make(map[int][]int) // id -> ids that depend on it

	for _, t := range tasks {
		idToTask[t.ID] = t
		if _, ok := inDegree[t.ID]; !ok {
			inDegree[t.ID] = 0
		}
		for _, dep := range t.Depends {
			adj[dep] = append(adj[dep], t.ID)
			inDegree[t.ID]++
		}
	}

	// BFS layer extraction
	var layers [][]TaskItem
	queue := []int{}
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}
	// Sort for determinism
	sort.Ints(queue)

	processed := 0
	for len(queue) > 0 {
		layerIDs := make([]int, len(queue))
		copy(layerIDs, queue)
		queue = queue[:0]

		var layer []TaskItem
		for _, id := range layerIDs {
			layer = append(layer, idToTask[id])
			processed++
			for _, next := range adj[id] {
				inDegree[next]--
				if inDegree[next] == 0 {
					queue = append(queue, next)
				}
			}
		}
		sort.Ints(queue)
		sort.Slice(layer, func(i, j int) bool { return layer[i].ID < layer[j].ID })
		layers = append(layers, layer)
	}

	// Cycle detection: if not all tasks processed, there is a cycle
	if processed != len(tasks) {
		return nil, ErrCyclicDependency
	}

	// File overlap split pass
	return splitByFileOverlap(layers), nil
}

// splitByFileOverlap ensures no two tasks in the same wave touch the same file.
// Uses exact string matching (case-sensitive) — sufficient for file path comparison.
func splitByFileOverlap(layers [][]TaskItem) [][]TaskItem {
	var result [][]TaskItem
	for _, layer := range layers {
		result = append(result, splitLayer(layer)...)
	}
	return result
}

// splitLayer uses greedy first-fit to place tasks into sublayers without file conflicts.
func splitLayer(tasks []TaskItem) [][]TaskItem {
	var sublayers [][]TaskItem
	for _, t := range tasks {
		placed := false
		for i := range sublayers {
			if !hasFileConflict(sublayers[i], t) {
				sublayers[i] = append(sublayers[i], t)
				placed = true
				break
			}
		}
		if !placed {
			sublayers = append(sublayers, []TaskItem{t})
		}
	}
	return sublayers
}

// hasFileConflict returns true if task t shares any file with the existing layer.
func hasFileConflict(layer []TaskItem, t TaskItem) bool {
	fileSet := make(map[string]struct{})
	for _, existing := range layer {
		for _, f := range existing.Files {
			fileSet[f] = struct{}{}
		}
	}
	for _, f := range t.Files {
		if _, ok := fileSet[f]; ok {
			return true
		}
	}
	return false
}

// HasParallelOpportunity returns true if any task has Depends or Files set.
// Used by SKILL.md decision point (D-03): only show wave mode prompt when true.
func HasParallelOpportunity(tasks []TaskItem) bool {
	for _, t := range tasks {
		if len(t.Depends) > 0 || len(t.Files) > 0 {
			return true
		}
	}
	return false
}
