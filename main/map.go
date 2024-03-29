package main

import (
	"fmt"
	"strconv"
	"strings"
)

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

func (node *Node) removeLink(linkedNode *Node) {
	if node == nil || linkedNode == nil {
		return
	}
	linkedNodeIndex := -1
	for i, n := range node.links {
		if n.index == linkedNode.index {
			linkedNodeIndex = i
			break
		}
	}
	if linkedNodeIndex == -1 {
		return
	}
	if len(node.links) == 1 {
		node.links = []*Node{}
		return
	}
	tmp := node.links[0]
	node.links[0] = node.links[linkedNodeIndex]
	node.links[linkedNodeIndex] = tmp
	node.links = node.links[1:]
}

func (node *Node) getNbLinkedExits() int {
	if node == nil {
		return 0
	}
	nbLinkedExits := 0
	for _, linkedNode := range node.links {
		if linkedNode.isExit {
			nbLinkedExits++
		}
	}
	return nbLinkedExits
}

func (node *Node) isLinkedToAnExit() bool {
	if node == nil {
		return false
	}
	for _, linkedNode := range node.links {
		if linkedNode.isExit {
			return true
		}
	}
	return false
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

func (gameMap *GameMap) removeLink(link *Link) {
	if gameMap == nil || link == nil {
		return
	}
	link.node1.removeLink(link.node2)
	link.node2.removeLink(link.node1)
	linkIndex := -1
	for i, l := range gameMap.links {
		if (l.node1.index == link.node1.index || l.node1.index == link.node2.index) &&
			(l.node2.index == link.node1.index || l.node2.index == link.node2.index) {
			linkIndex = i
			break
		}
	}
	if linkIndex == -1 {
		return
	}
	if len(gameMap.links) == 1 {
		gameMap.links = []*Link{}
		return
	}
	tmp := gameMap.links[0]
	gameMap.links[0] = gameMap.links[linkIndex]
	gameMap.links[linkIndex] = tmp
	gameMap.links = gameMap.links[1:]
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
