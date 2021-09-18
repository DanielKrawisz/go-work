package work

import (
	"encoding/binary"
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

	if len(initial.Share.ExtraNonce2) != 8 {
		return nil, errors.New("Extra nonce 2 must be of size 8. This limitation will be removed in future versions.")
	}

	var extraNonce2 uint64

	proof := Proof{
		Puzzle:   puzzle,
		Solution: initial}

	for !proof.Valid() {
		proof.Solution.Share.Nonce++
		if proof.Solution.Share.Nonce == 0 {
			extraNonce2++
			binary.BigEndian.PutUint64(proof.Solution.Share.ExtraNonce2, extraNonce2)
		}
	}

	return &proof, nil
}
