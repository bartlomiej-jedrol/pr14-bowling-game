// main handles Bowling game.
package main

import (
	"fmt"
	"math/rand"
)

const maxPins = 10
const maxFrames = 10

type game struct {
	frames []frame
	score  int
}

type frame struct {
	score    int
	rolls    []int
	isStrike bool
	isSpare  bool
}

// randomRoll stores randomRoll function - returns random integer within range.
var randomRoll = func(pins int) int {
	return rand.Intn(pins) + 1
}

// newGame creates new game.
func newGame() *game {
	return &game{}
}

// newFrame creates new frame.
func newFrame() *frame {
	return &frame{}
}

// addFrame adds frame to game.
func (g *game) addFrame(f *frame) {
	g.frames = append(g.frames, *f)
}

// addRoll adds roll with random score to frame.
func (f *frame) addRoll(pins int) {
	r := randomRoll(pins)
	f.score += r
	f.rolls = append(f.rolls, r)
	if r == maxPins {
		f.isStrike = true
	}
	if !f.isStrike && f.score == maxPins {
		f.isSpare = true
	}
}

// calculateFinalScore calculates final score of the game.
func calculateFinalScore(g *game) {
	for i, f := range g.frames {
		bonus := 0

		if i < maxFrames-1 {
			if i+1 < len(g.frames) { // next frame exists
				nextFrame := g.frames[i+1]

				if f.isSpare {
					bonus += nextFrame.rolls[0]
				}

				if f.isStrike {
					bonus += nextFrame.rolls[0]

					if nextFrame.isStrike {
						if i+2 < len(g.frames) { // next two frames exsist
							bonus += g.frames[i+2].rolls[0]
						} else {
							bonus += nextFrame.rolls[1]
						}
					} else {
						bonus += nextFrame.rolls[1]
					}
				}
			}
		}
		f.score += bonus
		g.frames[i] = f
		g.score += f.score
	}
}

// playGame plays a single Bowling game.
func playGame() {
	g := newGame()

	for i := range maxFrames {
		f := newFrame()

		// First roll
		f.addRoll(maxPins)
		if f.isStrike && i < maxFrames-1 {
			g.addFrame(f)
			continue
		}

		// Second roll
		remainingPins := maxPins - f.rolls[0]
		if remainingPins == 0 {
			remainingPins = maxPins
		}
		f.addRoll(remainingPins)

		// Third roll for the 10th frame
		if i == maxFrames-1 && (f.isStrike || f.isSpare) {
			remainingPins := maxPins - f.rolls[1]
			if remainingPins == 0 {
				remainingPins = maxPins
			}
			f.addRoll(remainingPins)
		}

		g.addFrame(f)
	}
	calculateFinalScore(g)
	for i, f := range g.frames {
		fmt.Printf("i: %d, f: %+v\n", i, f)
	}
	fmt.Printf("game score: %v\n", g.score)
	fmt.Printf("game: %v\n", g)
}

// main handles a single Bowling game.
func main() {
	playGame()
}
