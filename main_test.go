package main

import (
	"testing"
)

// contains verifies if collection contains scalar value.
func contains(scalar int, collection []int) bool {
	for i := 0; i < len(collection); i++ {
		if scalar == collection[i] {
			return true
		}
	}
	return false
}

// TestRandomRoll verifies if random roll is within range of provided pins.
func TestRandomRoll(t *testing.T) {
	tests := []struct {
		name         string
		pins         int
		allowedRange []int
	}{
		{
			name:         "Random roll within range of 10 pins",
			pins:         10,
			allowedRange: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:         "Random roll within range of 9 pins",
			pins:         9,
			allowedRange: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:         "Random roll with 1 pin",
			pins:         1,
			allowedRange: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := randomRoll(tt.pins)
			withinRange := contains(r, tt.allowedRange)
			if !withinRange {
				t.Errorf("Random roll is not within allowed range: %v", r)
			}
		})
	}
}

// mockRandomRoll replaces random roll with provided roll.
func mockRandomRoll(rolls []int) func(r int) int {
	i := 0
	return func(r int) int {
		roll := rolls[i]
		i++
		return roll
	}
}

// TestAddRoll verifies if roll is correctly added to frame along with required fields.
func TestAddRoll(t *testing.T) {
	tests := []struct {
		name          string
		rolls         []int
		expectedFrame frame
	}{
		{
			name:  "One roll, no strike",
			rolls: []int{7},
			expectedFrame: frame{
				score:    7,
				rolls:    []int{7},
				isStrike: false,
				isSpare:  false,
			},
		},
		{
			name:  "One roll, strike",
			rolls: []int{10},
			expectedFrame: frame{
				score:    10,
				rolls:    []int{10},
				isStrike: true,
				isSpare:  false,
			},
		},
		{
			name:  "Two rolls, no spare",
			rolls: []int{7, 1},
			expectedFrame: frame{
				score:    8,
				rolls:    []int{7, 1},
				isStrike: false,
				isSpare:  false,
			},
		},
		{
			name:  "Two rolls, strike",
			rolls: []int{0, 10},
			expectedFrame: frame{
				score:    10,
				rolls:    []int{0, 10},
				isStrike: true,
				isSpare:  false,
			},
		},
		{
			name:  "Two rolls, spare",
			rolls: []int{7, 3},
			expectedFrame: frame{
				score:    10,
				rolls:    []int{7, 3},
				isStrike: false,
				isSpare:  true,
			},
		},
		{
			name:  "Three rolls, spare",
			rolls: []int{9, 1, 2},
			expectedFrame: frame{
				score:    12,
				rolls:    []int{9, 1, 2},
				isStrike: false,
				isSpare:  true,
			},
		},
		{
			name:  "Three rolls, strike",
			rolls: []int{0, 10, 5},
			expectedFrame: frame{
				score:    15,
				rolls:    []int{0, 10, 5},
				isStrike: true,
				isSpare:  false,
			},
		},
	}

	// Store the orignial randomRoll function
	originalRandomRoll := randomRoll

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Replace randomRoll with mock
			randomRoll = mockRandomRoll(tt.rolls)

			f := newFrame()
			for _, roll := range tt.rolls {
				f.addRoll(roll)
			}

			// Reset randomRoll after the test
			defer func() {
				randomRoll = originalRandomRoll
			}()

			if len(f.rolls) != len(tt.expectedFrame.rolls) {
				t.Errorf("Length of rolls does not match between the frames. Expected: %v, Current: %v", len(tt.expectedFrame.rolls), len(f.rolls))
			}
			if f.score != tt.expectedFrame.score || f.isStrike != tt.expectedFrame.isStrike || f.isSpare != tt.expectedFrame.isSpare {
				t.Errorf("Rolls' fields do not match between frames. Expected: %v, Current: %v", f, tt.expectedFrame)
			}
			for i, roll := range f.rolls {
				if roll != tt.expectedFrame.rolls[i] {
					t.Errorf("Rolls do not match between frames. Expected: %v, Current: %v", tt.expectedFrame.rolls[i], roll)
				}
			}
		})
	}
}

// TestCalculateFinalScore verifies final score calculation of the Bowling game.
func TestCalculateFinalScore(t *testing.T) {
	tests := []struct {
		name          string
		game          game
		expectedScore int
	}{
		{
			name: "Random game",
			game: game{
				frames: []frame{
					{score: 10, rolls: []int{8, 2}, isStrike: false, isSpare: true},
					{score: 8, rolls: []int{3, 5}, isStrike: false, isSpare: false},
					{score: 8, rolls: []int{2, 6}, isStrike: false, isSpare: false},
					{score: 9, rolls: []int{8, 1}, isStrike: false, isSpare: false},
					{score: 10, rolls: []int{7, 3}, isStrike: false, isSpare: true},
					{score: 10, rolls: []int{8, 2}, isStrike: false, isSpare: true},
					{score: 10, rolls: []int{10}, isStrike: true, isSpare: false},
					{score: 9, rolls: []int{1, 8}, isStrike: false, isSpare: false},
					{score: 6, rolls: []int{3, 3}, isStrike: false, isSpare: false},
					{score: 24, rolls: []int{10, 6, 8}, isStrike: true, isSpare: false},
				},
				score: 0,
			},
			expectedScore: 134,
		},
		{
			name: "Perfect game",
			game: game{
				frames: []frame{
					{score: 10, rolls: []int{10}, isStrike: true, isSpare: false},
					{score: 10, rolls: []int{10}, isStrike: true, isSpare: false},
					{score: 10, rolls: []int{10}, isStrike: true, isSpare: false},
					{score: 10, rolls: []int{10}, isStrike: true, isSpare: false},
					{score: 10, rolls: []int{10}, isStrike: true, isSpare: false},
					{score: 10, rolls: []int{10}, isStrike: true, isSpare: false},
					{score: 10, rolls: []int{10}, isStrike: true, isSpare: false},
					{score: 10, rolls: []int{10}, isStrike: true, isSpare: false},
					{score: 10, rolls: []int{10}, isStrike: true, isSpare: false},
					{score: 30, rolls: []int{10, 10, 10}, isStrike: true, isSpare: false},
				},
				score: 0,
			},
			expectedScore: 300,
		},
		{
			name: "All spares with 5s",
			game: game{
				frames: []frame{
					{score: 10, rolls: []int{5, 5}, isStrike: false, isSpare: true},
					{score: 10, rolls: []int{5, 5}, isStrike: false, isSpare: true},
					{score: 10, rolls: []int{5, 5}, isStrike: false, isSpare: true},
					{score: 10, rolls: []int{5, 5}, isStrike: false, isSpare: true},
					{score: 10, rolls: []int{5, 5}, isStrike: false, isSpare: true},
					{score: 10, rolls: []int{5, 5}, isStrike: false, isSpare: true},
					{score: 10, rolls: []int{5, 5}, isStrike: false, isSpare: true},
					{score: 10, rolls: []int{5, 5}, isStrike: false, isSpare: true},
					{score: 10, rolls: []int{5, 5}, isStrike: false, isSpare: true},
					{score: 15, rolls: []int{5, 5, 5}, isStrike: false, isSpare: true},
				},
				score: 0,
			},
			expectedScore: 150,
		},
		{
			name: "Open frames",
			game: game{
				frames: []frame{
					{score: 7, rolls: []int{3, 4}, isStrike: false, isSpare: false},
					{score: 9, rolls: []int{5, 4}, isStrike: false, isSpare: false},
					{score: 6, rolls: []int{3, 3}, isStrike: false, isSpare: false},
					{score: 9, rolls: []int{4, 5}, isStrike: false, isSpare: false},
					{score: 6, rolls: []int{3, 3}, isStrike: false, isSpare: false},
					{score: 6, rolls: []int{1, 5}, isStrike: false, isSpare: false},
					{score: 7, rolls: []int{3, 4}, isStrike: false, isSpare: false},
					{score: 3, rolls: []int{2, 1}, isStrike: false, isSpare: false},
					{score: 9, rolls: []int{4, 5}, isStrike: false, isSpare: false},
					{score: 3, rolls: []int{1, 2}, isStrike: false, isSpare: false},
				},
				score: 0,
			},
			expectedScore: 65,
		},
		{
			name: "All nines",
			game: game{
				frames: []frame{
					{score: 9, rolls: []int{9, 0}, isStrike: false, isSpare: false},
					{score: 9, rolls: []int{9, 0}, isStrike: false, isSpare: false},
					{score: 9, rolls: []int{9, 0}, isStrike: false, isSpare: false},
					{score: 9, rolls: []int{9, 0}, isStrike: false, isSpare: false},
					{score: 9, rolls: []int{9, 0}, isStrike: false, isSpare: false},
					{score: 9, rolls: []int{9, 0}, isStrike: false, isSpare: false},
					{score: 9, rolls: []int{9, 0}, isStrike: false, isSpare: false},
					{score: 9, rolls: []int{9, 0}, isStrike: false, isSpare: false},
					{score: 9, rolls: []int{9, 0}, isStrike: false, isSpare: false},
					{score: 9, rolls: []int{9, 0}, isStrike: false, isSpare: false},
				},
				score: 0,
			},
			expectedScore: 90,
		},
		{
			name: "Last frame spare",
			game: game{
				frames: []frame{
					{score: 3, rolls: []int{1, 2}, isStrike: false, isSpare: false},
					{score: 7, rolls: []int{3, 4}, isStrike: false, isSpare: false},
					{score: 5, rolls: []int{2, 3}, isStrike: false, isSpare: false},
					{score: 6, rolls: []int{4, 2}, isStrike: false, isSpare: false},
					{score: 6, rolls: []int{5, 1}, isStrike: false, isSpare: false},
					{score: 5, rolls: []int{2, 3}, isStrike: false, isSpare: false},
					{score: 7, rolls: []int{3, 4}, isStrike: false, isSpare: false},
					{score: 3, rolls: []int{2, 1}, isStrike: false, isSpare: false},
					{score: 6, rolls: []int{4, 2}, isStrike: false, isSpare: false},
					{score: 15, rolls: []int{7, 3, 5}, isStrike: false, isSpare: true},
				},
				score: 0,
			},
			expectedScore: 63,
		},
		{
			name: "Last frame strike",
			game: game{
				frames: []frame{
					{score: 3, rolls: []int{1, 2}, isStrike: false, isSpare: false},
					{score: 7, rolls: []int{3, 4}, isStrike: false, isSpare: false},
					{score: 5, rolls: []int{2, 3}, isStrike: false, isSpare: false},
					{score: 6, rolls: []int{4, 2}, isStrike: false, isSpare: false},
					{score: 6, rolls: []int{5, 1}, isStrike: false, isSpare: false},
					{score: 5, rolls: []int{2, 3}, isStrike: false, isSpare: false},
					{score: 7, rolls: []int{3, 4}, isStrike: false, isSpare: false},
					{score: 3, rolls: []int{2, 1}, isStrike: false, isSpare: false},
					{score: 6, rolls: []int{4, 2}, isStrike: false, isSpare: false},
					{score: 18, rolls: []int{10, 3, 5}, isStrike: true, isSpare: false},
				},
				score: 0,
			},
			expectedScore: 66,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculateFinalScore(&tt.game)
			if tt.game.score != tt.expectedScore {
				t.Errorf("Calculated score does not match the exepcted score. Expected: %v, Current: %v", tt.expectedScore, tt.game.score)
			}
		})
	}
}
