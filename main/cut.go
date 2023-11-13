package main

import (
	"errors"
	"fmt"
)

func getLinkToCutFromPath(links []*Link, path *Path) (*Link, error) {
	if path == nil || len(path.indexes) < 2 {
		return nil, errors.New("cannot get a link to cut, the path has less than 2 indexes")
	}
	nbLinks := len(links)
	currentPath := path.indexes
	for len(currentPath) > 1 {
		index1 := currentPath[0]
		index2 := currentPath[1]
		for i := 0; i < nbLinks; i++ {
			link := links[i]
			if (link.node1.index == index1 || link.node1.index == index2) &&
				(link.node2.index == index1 || link.node2.index == index2) &&
				(link.node1.isExit || link.node2.isExit) {
				return link, nil
			}
		}
		currentPath = currentPath[1:]
	}
	return nil, errors.New("cannot get a link to cut from the path")
}

func getLinkToCutFromNode(node *Node) (*Link, error) {
	if node == nil || len(node.links) == 0 {
		return nil, errors.New("cannot get a link to cut, the node has no links")
	}
	return &Link{node1: node, node2: node.links[0]}, nil
}

func cutLink(gameMap *GameMap, link *Link) {
	if gameMap == nil || link == nil {
		return
	}
	gameMap.removeLink(link)
	fmt.Println(fmt.Sprintf("%d %d", link.node1.index, link.node2.index))
}
