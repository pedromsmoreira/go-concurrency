package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFinderWhenReceivesArrayWithOneOreReturnsArrayOfSize1(t *testing.T) {
	s := make([]string, 1)
	s = append(s, "ore")

	act := finder(s)

	assert.Equal(t, 1, len(act), "size of array should be 1")
	assert.Equal(t, "ore", act[0], "only item should be called ore")
}

func TestFinderWhenReceivesArrayWithoutOreReturnsEmptyArray(t *testing.T) {
	s := make([]string, 2)
	s = append(s, "rock")
	s = append(s, "rock")

	act := finder(s)

	assert.Equal(t, 0, len(act), "size of array should be 0")
}

func TestMinerWhenReceivesOneOreShouldReturnOneMinedOre(t *testing.T) {
	s := make([]string, 1)
	s = append(s, "ore")

	act := miner(s)

	assert.Equal(t, "minedOre", act[0], "only item should be called minedOre")
}

func TestMinerWhenReceivesTwoOresShouldReturnTwoMinedOres(t *testing.T) {
	s := make([]string, 1)
	s = append(s, "ore")
	s = append(s, "ore")

	act := miner(s)

	for i := 0; i < len(act); i++ {
		assert.Equal(t, "minedOre", act[i], "items should be called minedOre")
	}
}

func TestSmelterWhenReceivesOneMinedOreShouldReturnOneSmeltedOre(t *testing.T) {
	s := make([]string, 1)
	s = append(s, "minedOre")

	act := smelter(s)

	assert.Equal(t, "smeltedOre", act[0], "only item should be called smeltedOre")
}

func TestMinerWhenReceivesTwoMinedOresShouldReturnTwoSmeltedOres(t *testing.T) {
	s := make([]string, 1)
	s = append(s, "minedOre")
	s = append(s, "minedOre")

	act := smelter(s)

	for i := 0; i < len(act); i++ {
		assert.Equal(t, "smeltedOre", act[i], "items should be called smeltedOre")
	}
}
