package main

import "fmt"

func main() {
	theMine := []string{"rock", "ore", "ore", "rock", "ore"}
	foundOre := finder(theMine)
	minedOre := miner(foundOre)
	smelter(minedOre)
}

func finder(mine []string) []string {
	fo := make([]string, 0)
	for o := 0; o < len(mine); o++ {
		if mine[o] == "ore" {
			fo = append(fo, mine[o])
		}
	}
	fmt.Printf("From Finder: %v\n", fo)
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
