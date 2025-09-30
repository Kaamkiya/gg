package shape

import (
	"math/rand"
	"slices"
)

// Randomizer makes the randrom pick of shapes to fill less 'unfair'. Inspired by info found
// here: https://tetris.fandom.com/wiki/TGM_randomizer
type Randomizer struct {
	lastValues []int
}

func (r *Randomizer) nextInt(maxValue int) int {
	nextShape := rand.Intn(maxValue)

	retries := 0
	for retries < 6 && slices.Contains(r.lastValues, nextShape) {
		nextShape = rand.Intn(maxValue)
		retries++
	}

	r.lastValues = append(r.lastValues, nextShape)
	r.lastValues = r.lastValues[1:]

	return nextShape
}

func NewRandomizer() *Randomizer {
	lastValues := make([]int, 4)

	lastValues[0] = Z
	lastValues[1] = S
	lastValues[2] = Z
	lastValues[3] = S

	return &Randomizer{
		lastValues,
	}
}
