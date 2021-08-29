package work

import (
	"encoding/hex"

	"github.com/libsv/go-bc"
)

func CPUSolve(puzzle Puzzle, initial Solution) (*Proof, error) {
	// check that the difficulty is not too big (1 or less)
	difficulty, err := bc.DifficultyFromBits(hex.EncodeToString(puzzle.Candidate.Bits))
	if err != nil {
		return nil, err
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
