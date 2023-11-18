package main

import (
	"fmt"
)

// Initialize the game loop and play the game
func main() {
	timer := &Timer{}
	timer.startInit()
	gameMap := buildMap()
	round := 0
	timer.endInit()
	debug("Game map:", gameMap)
	debug("timer", timer)
	for {
		timer.startRound()
		round++
		var bobnetAgentIndex int
		fmt.Scan(&bobnetAgentIndex)
		gameMap.SetBobnetAgentIndex(bobnetAgentIndex)
		debug("Round", round, ":", gameMap.nodes)
		bobnetPath, bobnetPathError := getBobnetPath(gameMap)
		if bobnetPathError != nil {
			debug("Error getting bobnet shortest path:", bobnetPathError)
			continue
		}
		debug("Bobnet path", bobnetPath)
		linkToCut, linkToCutError := getLinkToCutFromPath(gameMap.links, bobnetPath)
		if linkToCutError != nil {
			debug("Error getting link to cut:", linkToCutError)
			debug("Get a link from bobnet agent node")
			linkToCut, linkToCutError = getLinkToCutFromNode(gameMap.nodes[gameMap.bobnetAgentIndex])
			if linkToCutError != nil {
				debug("Error getting link to cut:", linkToCutError)
				continue
			}
		}
		cutLink(gameMap, linkToCut)
		timer.endRound()
		debug("timer", timer)
	}
}
