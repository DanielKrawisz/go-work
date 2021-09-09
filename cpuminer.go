package work

import (
	"encoding/hex"
	"errors"

	"github.com/libsv/go-bc"
)

func CPUSolve(puzzle Puzzle, initial Solution) (*Proof, error) {
	difficulty, err := bc.DifficultyFromBits(hex.EncodeToString(puzzle.Candidate.Bits))
	if err != nil {
		return nil, err
	}

	if difficulty > 1.01 {
		return nil, errors.New("Difficulty is too big. CPU mining is only for difficulty <= 1")
	}

	proof := Proof{
		Puzzle:   puzzle,
		Solution: initial}

	for !proof.Valid() {
		proof.Solution.Share.Nonce++
		if proof.Solution.Share.Nonce == 0 {
			proof.Solution.Share.ExtraNonce2++
		}
	}

	return &proof, nil
}
