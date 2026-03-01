package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Mainnet Genesis Time: December 1, 2020, 12:00:23 UTC
const GenesisTime = 1606824023
const SecondsPerSlot = 12

// Beacon API structures
type BeaconHeaderResponse struct {
	Data struct {
		Header struct {
			Message struct {
				Slot       string `json:"slot"`
				StateRoot  string `json:"state_root"`
				ParentRoot string `json:"parent_root"`
			} `json:"message"`
		} `json:"header"`
	} `json:"data"`
}

func main() {
	// Public Beacon Node Endpoint (provided by ethereumpow.org or similar community nodes)
	// You can also use an Infura/Alchemy URL here.
	beaconNodeURL := "https://ethereum-beacon-api.publicnode.com/eth/v1/beacon/headers/head"

	fmt.Println("--- Ethereum Consensus Layer Heartbeat Audit ---")

	// 1. Calculate Expected Slot based on Genesis
	now := time.Now().Unix()
	elapsed := now - GenesisTime
	expectedSlot := elapsed / SecondsPerSlot

	fmt.Printf("Current Unix Time:    %d\n", now)
	fmt.Printf("Expected Global Slot: %d\n", expectedSlot)

	// 2. Fetch Actual Head from the Beacon Node
	resp, err := http.Get(beaconNodeURL)
	if err != nil {
		fmt.Printf("Error fetching from node: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var header BeaconHeaderResponse
	if err := json.Unmarshal(body, &header); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	actualSlot := header.Data.Header.Message.Slot
	stateRoot := header.Data.Header.Message.StateRoot

	// 3. Output Audit Results
	fmt.Printf("Node Reported Slot:   %s\n", actualSlot)
	fmt.Printf("Latest State Root:    %s\n", stateRoot)

	// Logic Check: Is the node lagging?
	fmt.Println("\n--- Audit Summary ---")
	if fmt.Sprintf("%d", expectedSlot) == actualSlot {
		fmt.Println("[PASS] Node is perfectly synchronized with the global slot clock.")
	} else {
		fmt.Println("[WARN] Node slot mismatch. This could indicate latency or a synchronization lag.")
	}
}
