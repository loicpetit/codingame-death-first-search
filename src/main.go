package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*** UTIL ***/

func debug(values ...any) {
	fmt.Fprintln(os.Stderr, values...)
}

/*** MAP ***/

type Node struct {
	index                int
	isBobnetAgentPresent bool
	isExit               bool
	links                []*Node
}

func (node *Node) String() string {
	if node == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(strconv.FormatInt(int64(node.index), 10))
	if node.isExit {
		sb.WriteString(">")
	}
	if node.isBobnetAgentPresent {
		sb.WriteString("*")
	}
	sb.WriteString("(")
	for i := 0; i < len(node.links); i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.FormatInt(int64(node.links[i].index), 10))
	}
	sb.WriteString(")")
	return sb.String()
}

type Link struct {
	node1 *Node
	node2 *Node
}

func (link *Link) String() string {
	if link == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("(")
	if link.node1 != nil {
		sb.WriteString(strconv.FormatInt(int64(link.node1.index), 10))
	}
	if link.node2 != nil {
		if link.node1 != nil {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.FormatInt(int64(link.node2.index), 10))
	}
	sb.WriteString(")")
	return sb.String()
}

type GameMap struct {
	bobnetAgentIndex int
	exits            []*Node
	links            []*Link
	nodes            []*Node
}

func (gameMap *GameMap) String() string {
	if gameMap == nil {
		return ""
	}
	return fmt.Sprintf("{nodes: %v, links: %v, exits: %v}", gameMap.nodes, gameMap.links, gameMap.exits)
}

func (gameMap *GameMap) SetBobnetAgentIndex(index int) {
	if gameMap == nil {
		return
	}
	nbNodes := len(gameMap.nodes)
	if gameMap.bobnetAgentIndex >= 0 && gameMap.bobnetAgentIndex < nbNodes {
		gameMap.nodes[gameMap.bobnetAgentIndex].isBobnetAgentPresent = false
	}
	if index >= 0 && index < nbNodes {
		gameMap.bobnetAgentIndex = index
		gameMap.nodes[index].isBobnetAgentPresent = true
	} else {
		gameMap.bobnetAgentIndex = -1
	}
}

func buildMap() *GameMap {
	var nbNodes, nbLinks, nbExits int
	fmt.Scan(&nbNodes, &nbLinks, &nbExits)
	debug("nb nodes:", nbNodes)
	debug("nb links:", nbLinks)
	debug("nb exits:", nbExits)
	nodes := make([]*Node, nbNodes)
	for i := 0; i < nbNodes; i++ {
		nodes[i] = &Node{index: i, isBobnetAgentPresent: false, isExit: false}
	}
	links := make([]*Link, nbLinks)
	for i := 0; i < nbLinks; i++ {
		var nodeIndex1, nodeIndex2 int
		fmt.Scan(&nodeIndex1, &nodeIndex2)
		node1 := nodes[nodeIndex1]
		node2 := nodes[nodeIndex2]
		links[i] = &Link{node1: node1, node2: node2}
		node1.links = append(node1.links, node2)
		node2.links = append(node2.links, node1)
	}
	exits := make([]*Node, nbExits)
	for i := 0; i < nbExits; i++ {
		var exitNodeIndex int
		fmt.Scan(&exitNodeIndex)
		exits[i] = nodes[exitNodeIndex]
		exits[i].isExit = true
	}
	return &GameMap{bobnetAgentIndex: -1, exits: exits, links: links, nodes: nodes}
}

/*** MAIN ***/

func main() {
	gameMap := buildMap()
	debug("Game map:", gameMap)
	for {
		var bobnetAgentIndex int
		fmt.Scan(&bobnetAgentIndex)
		gameMap.SetBobnetAgentIndex(bobnetAgentIndex)
		debug("Round:", gameMap.nodes)

		// Example: 0 1 are the indices of the nodes you wish to sever the link between
		fmt.Println("0 1")
	}
}
