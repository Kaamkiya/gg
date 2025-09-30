package shape

import (
	"strconv"
	"testing"
)

func TestNewRandomizerHasSZ(t *testing.T) {
	randomizer := NewRandomizer()

	if randomizer.lastValues[0] != Z ||
		randomizer.lastValues[1] != S ||
		randomizer.lastValues[2] != Z ||
		randomizer.lastValues[3] != S {
		t.Fatal("Initial state of the randomizer should be SZSZ")
	}

}

func TestNewRandomizerUpdatesStateCorrectlyOnNewInt(t *testing.T) {
	randomizer := NewRandomizer()

	firstShape := randomizer.nextInt(7)
	secondShape := randomizer.nextInt(7)
	thirdShape := randomizer.nextInt(7)

	if randomizer.lastValues[0] != S ||
		randomizer.lastValues[1] != firstShape ||
		randomizer.lastValues[2] != secondShape ||
		randomizer.lastValues[3] != thirdShape {
		t.Fatal("The state of the randomizer should be Z" + strconv.Itoa(firstShape) + strconv.Itoa(secondShape) + strconv.Itoa(thirdShape))
	}

}
