package main

import (
	"fmt"
	"time"
)

func main() {
	theMine := []string{"ore1", "ore2", "ore3"}
	oreChan := make(chan string)
	minedOreChan := make(chan string)

	// Finder
	go func(mine []string) {
		for _, item := range mine {
			oreChan <- item
		}
	}(theMine)

	// Ore Breaker
	go func() {
		for foundOre := range oreChan {
			fmt.Println("Miner Received: " + foundOre + " from finder")
			minedOreChan <- "minedOre"
		}
	}()

	// Smelter
	go func() {
		for minedOre := range minedOreChan {
			fmt.Println("From Miner: ", minedOre)
			fmt.Println("From Smelter: Ore is smelted")
		}
	}()

	<-time.After(time.Second * 5)
}

func finder(mine []string) []string {
	fo := make([]string, 0)
	for o := 0; o < len(mine); o++ {
		if mine[o] == "ore" {
			fo = append(fo, mine[o])
			fmt.Println("From Finder found ore!")
		}
	}
	fmt.Printf("From Finder: %v\n", fo)
	return fo
}

func finder1(mine []string) []string {
	fo := make([]string, 0)
	for o := 0; o < len(mine); o++ {
		if mine[o] == "ore" {
			fo = append(fo, mine[o])
			fmt.Println("From Finder1 found ore!")
		}
	}
	fmt.Printf("From Finder1: %v\n", fo)
	return fo
}

func finder2(mine []string) []string {
	fo := make([]string, 0)
	for o := 0; o < len(mine); o++ {
		if mine[o] == "ore" {
			fo = append(fo, mine[o])
			fmt.Println("From Finder2 found ore!")
		}
	}
	fmt.Printf("From Finder2: %v\n", fo)
	return fo
}

func miner(ore []string) []string {
	mo := make([]string, 0)
	for i := 0; i < len(ore); i++ {
		mo = append(mo, "minedOre")
	}
	fmt.Printf("From Miner: %v\n", mo)
	return mo
}

func smelter(mo []string) []string {
	so := make([]string, 0)
	for i := 0; i < len(mo); i++ {
		so = append(so, "minedOre")
	}
	fmt.Printf("From Smelter: %v\n", so)
	return so
}
