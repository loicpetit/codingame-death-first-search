package main

import (
	"errors"
	"fmt"
)

type Path struct {
	indexes []int
	risk    int
}

func (path *Path) String() string {
	if path == nil {
		return ""
	}
	return fmt.Sprintf("{indexes: %v, risk: %v}", path.indexes, path.risk)
}

func getShortestPath(gameMap *GameMap, startIndex int, endIndex int) []int {
	// debug("Get shortest path between", startIndex, "and", endIndex)
	if startIndex == endIndex {
		return []int{startIndex}
	}
	nbNodes := len(gameMap.nodes)
	parentIndex := make([]int, nbNodes)
	for i := range parentIndex {
		parentIndex[i] = -1
	}
	var queue []int
	queue = append(queue, startIndex)
	for {
		currentIndex := queue[0]
		// debug("current index", currentIndex)
		queue = queue[1:]
		currentNode := gameMap.nodes[currentIndex]
		for _, linkedNode := range currentNode.links {
			if linkedNode.index != startIndex && parentIndex[linkedNode.index] == -1 {
				parentIndex[linkedNode.index] = currentIndex
				queue = append(queue, linkedNode.index)
			}
		}
		if len(queue) == 0 || parentIndex[endIndex] != -1 {
			break
		}
	}
	// debug("Parents", parentIndex)
	if parentIndex[endIndex] == -1 {
		return []int{}
	}
	// debug("Create indexes")
	var indexes []int
	currentParentIndex := endIndex
	for parentIndex[currentParentIndex] != -1 {
		// debug("prepend", currentParentIndex)
		indexes = append([]int{currentParentIndex}, indexes...)
		currentParentIndex = parentIndex[currentParentIndex]
	}
	indexes = append([]int{startIndex}, indexes...)
	return indexes
}

func getTunnelPath(gameMap *GameMap, startIndex int, beforeEndIndex int, endIndex int) []int {
	// debug("Get tunnel path between", startIndex, "and", endIndex)
	if startIndex == endIndex {
		return []int{startIndex}
	}
	nbNodes := len(gameMap.nodes)
	parentIndex := make([]int, nbNodes)
	valueIndex := make([]int, nbNodes)
	visitedIndex := make([]bool, nbNodes)
	for i := 0; i < nbNodes; i++ {
		parentIndex[i] = -1
		valueIndex[i] = 0
		visitedIndex[i] = false
	}
	visitedIndex[beforeEndIndex] = true
	parentIndex[beforeEndIndex] = endIndex
	visitedIndex[endIndex] = true
	valueIndex[endIndex] = getTunnelNodeValue(gameMap, beforeEndIndex, endIndex, endIndex)
	var queue []int
	queue = append(queue, beforeEndIndex)
	for {
		currentIndex := queue[0]
		if currentIndex == startIndex {
			break
		}
		// debug("current index", currentIndex)
		queue = queue[1:]
		currentNode := gameMap.nodes[currentIndex]
		value := getTunnelNodeValue(gameMap, beforeEndIndex, endIndex, currentIndex)
		valueIndex[currentIndex] = value
		for _, linkedNode := range currentNode.links {
			linkedNodeIndex := linkedNode.index
			if linkedNodeIndex == endIndex {
				continue
			}
			if parentIndex[linkedNodeIndex] == -1 ||
				(parentIndex[linkedNodeIndex] != currentIndex &&
					parentIndex[currentIndex] != linkedNodeIndex &&
					valueIndex[parentIndex[linkedNodeIndex]] < value) {
				parentIndex[linkedNodeIndex] = currentIndex
			}
			if !visitedIndex[linkedNodeIndex] {
				visitedIndex[linkedNodeIndex] = true
				queue = append(queue, linkedNodeIndex)
			}
		}
		if len(queue) == 0 {
			break
		}
	}
	// if beforeEndIndex == 6 && endIndex == 0 {
	// 	fmt.Fprint(os.Stderr, "Parents ")
	// 	for i, v := range parentIndex {
	// 		fmt.Fprint(os.Stderr, fmt.Sprintf("[%d:%d]", i, v))
	// 	}
	// 	fmt.Fprintln(os.Stderr)
	// 	fmt.Fprint(os.Stderr, "Values ")
	// 	for i, v := range valueIndex {
	// 		fmt.Fprint(os.Stderr, fmt.Sprintf("[%d:%d]", i, v))
	// 	}
	// 	fmt.Fprintln(os.Stderr)
	// }
	if parentIndex[startIndex] == -1 {
		return []int{}
	}
	// debug("Create indexes")
	var indexes []int
	currentParentIndex := startIndex
	parentCount := 0
	for parentIndex[currentParentIndex] != -1 {
		parentCount++
		// debug("append", currentParentIndex)
		indexes = append(indexes, currentParentIndex)
		currentParentIndex = parentIndex[currentParentIndex]
		if parentCount == nbNodes {
			return []int{}
		}
	}
	indexes = append(indexes, endIndex)
	// debug("Indexes ", indexes)
	return indexes
}

func getTunnelNodeValue(gameMap *GameMap, beforeEndIndex int, endIndex int, index int) int {
	if endIndex == index {
		return 10000
	}
	if beforeEndIndex == index {
		return 100000
	}
	if index < 0 || index >= len(gameMap.nodes) {
		return 0
	}
	return gameMap.nodes[index].getNbLinkedExits()
}

func evaluateRisk(gameMap *GameMap, indexes []int) int {
	if gameMap == nil || len(indexes) < 2 {
		return 0
	}
	lengthRisk := evaluateLengthRisk(gameMap, indexes)
	multiExistsNodeRisk := evaluateMultiExitsNodeRisk(gameMap, indexes)
	exitNextTurnRisk := evaluateExitNextTurnRisk(indexes)
	tunnelRisk := evaluateTunnelRisk(gameMap, indexes)
	//debug(lengthRisk, multiExistsNodeRisk, exitNextTurnRisk, tunnelRisk)
	return lengthRisk + multiExistsNodeRisk + exitNextTurnRisk + tunnelRisk
}

func evaluateLengthRisk(gameMap *GameMap, indexes []int) int {
	if gameMap == nil {
		return 0
	}
	return len(gameMap.nodes) - len(indexes)
}

func evaluateMultiExitsNodeRisk(gameMap *GameMap, indexes []int) int {
	if gameMap == nil || len(indexes) < 2 {
		return 0
	}
	maxNbExits := 0
	nbNodes := len(gameMap.nodes)
	nbIndexes := len(indexes)
	for i := 0; i < nbIndexes; i++ {
		nodeIndex := indexes[i]
		if nodeIndex < 0 || nodeIndex >= nbNodes {
			continue
		}
		node := gameMap.nodes[nodeIndex]
		nbLinkedExits := node.getNbLinkedExits()
		if nbLinkedExits > maxNbExits {
			maxNbExits = nbLinkedExits
		}
	}
	return maxNbExits * 1000
}

func evaluateExitNextTurnRisk(indexes []int) int {
	if len(indexes) == 2 {
		return 10000
	}
	return 0
}

func evaluateTunnelRisk(gameMap *GameMap, indexes []int) int {
	if gameMap == nil || len(indexes) < 2 {
		return 0
	}
	nbNodes := len(gameMap.nodes)
	nbIndexes := len(indexes)
	nbConsecutiveLinkedExits := 0
	for i := nbIndexes - 1; i >= 0; i-- {
		nodeIndex := indexes[i]
		if nodeIndex < 0 || nodeIndex >= nbNodes {
			continue
		}
		node := gameMap.nodes[nodeIndex]
		if node.isExit || node.isLinkedToAnExit() {
			nbConsecutiveLinkedExits++
			continue
		}
		break
	}
	return nbConsecutiveLinkedExits
}

func getShortestPathToExitLink(channel chan *Path, gameMap *GameMap, exitLink *Link) {
	var indexes []int
	if exitLink.node1.isExit {
		indexes = append(getShortestPath(gameMap, gameMap.bobnetAgentIndex, exitLink.node2.index), exitLink.node1.index)
	} else {
		indexes = append(getShortestPath(gameMap, gameMap.bobnetAgentIndex, exitLink.node1.index), exitLink.node2.index)
	}
	risk := evaluateRisk(gameMap, indexes)
	channel <- &Path{indexes: indexes, risk: risk}
}

func getTunnelPathToExitLink(channel chan *Path, gameMap *GameMap, exitLink *Link) {
	var indexes []int
	if exitLink.node1.isExit {
		// debug("get tunnel from", gameMap.bobnetAgentIndex, "by", exitLink.node2.index, "to", exitLink.node1.index)
		indexes = getTunnelPath(gameMap, gameMap.bobnetAgentIndex, exitLink.node2.index, exitLink.node1.index)
	} else {
		// debug("get tunnel from", gameMap.bobnetAgentIndex, "by", exitLink.node1.index, "to", exitLink.node2.index)
		indexes = getTunnelPath(gameMap, gameMap.bobnetAgentIndex, exitLink.node1.index, exitLink.node2.index)
	}
	risk := evaluateRisk(gameMap, indexes)
	channel <- &Path{indexes: indexes, risk: risk}
}

func getAllExitLinks(gameMap *GameMap) (links []*Link) {
	if gameMap == nil || len(gameMap.links) == 0 {
		return
	}
	for _, link := range gameMap.links {
		if link.node1.isExit || link.node2.isExit {
			links = append(links, link)
		}
	}
	return
}

func getBobnetPath(gameMap *GameMap) (*Path, error) {
	if gameMap == nil {
		return nil, errors.New("game map is missing")
	}
	exitLinks := getAllExitLinks(gameMap)
	// debug("Exit links:", exitLinks)
	nbExitLinks := len(exitLinks)
	nbPaths := 2 * nbExitLinks
	debug(nbPaths, "paths to compute")
	pathChannel := make(chan *Path, nbPaths)
	defer close(pathChannel)
	for i := 0; i < nbExitLinks; i++ {
		go getShortestPathToExitLink(pathChannel, gameMap, exitLinks[i])
		go getTunnelPathToExitLink(pathChannel, gameMap, exitLinks[i])
	}
	var path *Path
	for i := 0; i < nbPaths; i++ {
		pathToExit := <-pathChannel
		if pathToExit == nil || len(pathToExit.indexes) < 2 {
			continue
		}
		debug("Possible path:", pathToExit)
		if path == nil || pathToExit.risk > path.risk ||
			(pathToExit.risk == path.risk && len(pathToExit.indexes) < len(path.indexes)) {
			path = pathToExit
		}
	}
	return path, nil
}
