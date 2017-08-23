package main

import (
	"fmt"
	"errors"
	"sort"
)

var (
	A, B []int
	K int
)

type Node struct {
	Id int
	Nearby []*Node
	Children []*Node
	Height int
	Parent *Node
}

type AllowedNodeInfo struct{
	Height int
	Id int
	Child *Node
}
type NodesByHeight []AllowedNodeInfo

// Sortable interface implementation for []AllowedNodeInfo
func (a NodesByHeight) Len() int           { return len(a) }
func (a NodesByHeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NodesByHeight) Less(i, j int) bool { return a[i].Height < a[j].Height }

type Graph []Node

var graph Graph

func main() {
	A = []int{5, 1, 0, 2, 7, 0, 6, 6, 1}
	B = []int{1, 0, 7, 4, 2, 6, 8, 3, 9}
	K = 2
	fmt.Println(Solution(A, B, K))
}

func Solution(A []int, B []int, K int) int {

	// array's length should be equal
	if len(A) != len(B) {
		panic(errors.New("Wrong input data"))
	}

	// when we have speed cameras for each road
	if len(A) == K {
		return 0
	}

	// Nodes definition
	graph = make(Graph, len(A) + 1)
	for i := 0; i < len(A) + 1; i++ {
		graph[i].Id = i
	}
	for i := 0; i < len(A); i++ {
		graph[A[i]].Nearby = append(graph[A[i]].Nearby, &graph[B[i]])
		graph[B[i]].Nearby = append(graph[B[i]].Nearby, &graph[A[i]])
	}

	// Recursive definition parent-child relations between nodes 
	simpleDfs(&graph[0], nil)

	// Search solution via binary search
	return solutionBinarySearch()
}

func solutionBinarySearch() (result int) {
	left := 0
	// Because we know about max length of road
	right := 900

	for left < right {
		mid := (left + right) / 2

		// For each iteration we should calculate minimum value of cameras
		neededCameras := neededCamerasCount(&graph[0], mid)

		// Then change borders of search
		if neededCameras <= K {
			right = mid
		} else {
			left = mid + 1
		}
	}

	result = left
	return
}

// Simple implementation of depth-first search; it's all what we needed:
// for each node define child and parent nodes
func simpleDfs(node *Node, parent *Node) {
	node.Parent = parent
	for _, v := range node.Nearby {
		if v != parent {
			node.Children = append(node.Children, v)
			simpleDfs(v, node)
		}
	}
}

// Calculating minimum cameras for current diameter
func neededCamerasCount(node *Node, diameter int) (result int)  {
	node.Height = 0
	for _, v := range node.Children {
		result += neededCamerasCount(v, diameter)
	}

	// Search allowed nodes when height less then diameter
	var allowedNodes NodesByHeight
	for k, v := range node.Children {
		if v.Height < diameter {
			allowedNodes = append(allowedNodes, AllowedNodeInfo{v.Height, k, v})
		}
	}
	sort.Sort(NodesByHeight(allowedNodes))

	result += len(node.Children) - len(allowedNodes)

	// Remove extra nodes while the sum of penultimate and last heights + 2 (because we have 2 roads to parent node)
	// is more than diameter or while we have more than 1 item in slice. :|
	for len(allowedNodes) > 1 && allowedNodes[len(allowedNodes) - 1].Height + allowedNodes[len(allowedNodes) - 2].Height + 2 > diameter {
		allowedNodes = allowedNodes[:len(allowedNodes)-1]
		result++
	}

	// Node height is maximum of (biggest child node + 1) and current height
	for _, v := range allowedNodes {
		node.Height = max(v.Height + 1, node.Height)
	}

	return
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}