// NOT FROM ME
// get it from coding game solutions
// code is a lot shorter!

package main

import (
	"fmt"
	"os"
)

const MaxSize = 100

func calcMoveDisatnces(nodesCount int, edges [100]int, exists [MaxSize]bool, adjMatrix [100][100]int) (distMatrix [100][100]int) {
	distMatrix = adjMatrix

	// this blockk is always same and potentially can be optimized
	for i := 0; i < nodesCount; i++ {
		for j := 0; j < nodesCount; j++ {
			if adjMatrix[i][j] == 0 {
				distMatrix[i][j] = nodesCount // somewhere far-far away
			} else if edges[j] > 0 {
				distMatrix[i][j] = 0
			} else if exists[j] {
				distMatrix[i][j] = nodesCount
			} else {
				distMatrix[i][j] = adjMatrix[i][j]
			}
		}

		distMatrix[i][i] = 0
	}

	for currentEdge := 0; currentEdge < nodesCount; currentEdge++ {
		if edges[currentEdge] == 0 {
			continue // do not calculate distances to non-edge nodes
		}

		for currentPass := 1; currentPass < nodesCount; currentPass++ {
			for firstHop := 0; firstHop < nodesCount; firstHop++ {
				for secondHop := 0; secondHop < nodesCount; secondHop++ {
					currentDistance :=
						distMatrix[currentEdge][firstHop] + // cost of move, positive spend of skynet
							distMatrix[firstHop][secondHop] // cost of staying, zero in case of node, negative in case of edges

					if !exists[firstHop] && !exists[secondHop] && currentDistance < distMatrix[currentEdge][secondHop] {
						distMatrix[currentEdge][secondHop] = currentDistance
					}
				}
			}
		}
	}

	return
}

func calcEdges(adjMatrix [MaxSize][MaxSize]int, nodesCount int, exits [MaxSize]bool) (edges [MaxSize]int) {
	for exit, isExit := range exits {
		if isExit {
			for i := 0; i < nodesCount; i++ {
				if adjMatrix[exit][i] == 1 {
					edges[i]++
				}
			}
		}
	}
	return
}

func mostEndangeredEdge(skynetNode int, edges [MaxSize]int, nodesCount int, distMatrix [MaxSize][MaxSize]int) int {
	bestEdge := 0
	lowestDistance := MaxSize

	for edge := 0; edge < nodesCount; edge++ {
		if edges[edge] > 0 {
			if edge == skynetNode {
				return edge
			}

			curDistance := distMatrix[edge][skynetNode] - edges[edge]

			if curDistance < lowestDistance {
				bestEdge = edge
				lowestDistance = curDistance
			}
		}
	}

	return bestEdge
}

func disconnectEdge(edge int, exists [MaxSize]bool, adjMatrix *[MaxSize][MaxSize]int) int {
	for exit, isExit := range exists {
		if isExit && adjMatrix[edge][exit] == 1 {
			adjMatrix[edge][exit] = 0
			adjMatrix[exit][edge] = 0

			return exit
		}
	}

	return -1
}

func main() {
	// nodesCount: the total number of nodes in the level, including the gateways
	// linksCount: the number of links
	// exitsCount: the number of exit gateways
	var nodesCount, linksCount, exitsCount int
	fmt.Scan(&nodesCount, &linksCount, &exitsCount)

	adjacencyMatrix := [MaxSize][MaxSize]int{} // who cares about the memory...
	exists := [MaxSize]bool{}

	for i := 0; i < linksCount; i++ {
		// N1: N1 and N2 defines a link between these nodes
		var N1, N2 int
		fmt.Scan(&N1, &N2)

		adjacencyMatrix[N1][N2] = 1
		adjacencyMatrix[N2][N1] = 1
	}

	for i := 0; i < exitsCount; i++ {
		// EI: the index of a gateway node
		var EI int
		fmt.Scan(&EI)

		exists[EI] = true
	}

	for {
		var skynetNode int
		fmt.Scan(&skynetNode)

		edges := calcEdges(adjacencyMatrix, nodesCount, exists)
		fmt.Fprintln(os.Stderr, edges)

		distances := calcMoveDisatnces(nodesCount, edges, exists, adjacencyMatrix)
		for edge, linkStrength := range edges {
			if linkStrength > 0 {
				fmt.Fprintln(os.Stderr, edge, "distance = ", distances[edge])
			}
		}

		edge := mostEndangeredEdge(skynetNode, edges, nodesCount, distances)
		exit := disconnectEdge(edge, exists, &adjacencyMatrix)

		fmt.Fprintln(os.Stderr, "mostEndangeredEdge", edges)
		fmt.Printf("%d %d\n", edge, exit)
	}
}
