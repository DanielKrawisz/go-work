package work

import (
	"encoding/binary"
	"encoding/hex"
	"errors"

	"github.com/libsv/go-bc"
	"github.com/ordishs/go-bitcoin"
)

// Candidate includes the parts of the blockheader that are known
// before the work has been done. It represents the parts of a work
// puzzle that are provided by getminingcandidate.
type Candidate struct {
	Version    uint32
	Digest     []byte
	Bits       []byte
	MerklePath []string
}

func MakeCandidate(mc *bitcoin.MiningCandidate) (*Candidate, error) {
	bits, err := hex.DecodeString(mc.Bits)
	if err != nil {
		return nil, err
	}
	prevHash, err := hex.DecodeString(mc.PreviousHash)
	if err != nil {
		return nil, err
	}
	return &Candidate{
		Version:    mc.Version,
		Digest:     prevHash,
		Bits:       bits,
		MerklePath: mc.MerkleProof}, nil
}

// Puzzle is what we get after the coinbase has been derived.
type Puzzle struct {
	Candidate     Candidate
	CoinbaseBegin []byte
	CoinbaseEnd   []byte
	VersionMask   *uint32
}

func MakePuzzle(candidate Candidate, coinbaseBegin []byte, coinbaseEnd []byte) Puzzle {
	return Puzzle{
		Candidate:     candidate,
		CoinbaseBegin: coinbaseBegin,
		CoinbaseEnd:   coinbaseEnd,
		VersionMask:   nil}
}

func MakePuzzleASICBoost(candidate Candidate, coinbaseBegin []byte, coinbaseEnd []byte, mask uint32) Puzzle {
	bits := new(uint32)
	*bits = mask
	return Puzzle{
		Candidate:     candidate,
		CoinbaseBegin: coinbaseBegin,
		CoinbaseEnd:   coinbaseEnd,
		VersionMask:   bits}
}

// A job is a work puzzle after a stratum id has been assigned to a
// given worker by the mining pool.
type Job struct {
	Puzzle Puzzle

	// ExtraNonce1 is also the user id in Stratum.
	ExtraNonce1 uint32
}

// A share is the data returned by the worker. Job + Share = Proof
type Share struct {
	Time               uint32
	Nonce              uint32
	ExtraNonce2        uint64
	GeneralPurposeBits *uint32
}

func MakeShare(time uint32, nonce uint32, extraNonce2 uint64) Share {
	return Share{
		Time:               time,
		Nonce:              nonce,
		ExtraNonce2:        extraNonce2,
		GeneralPurposeBits: nil}
}

func MakeShareASICBoost(time uint32, nonce uint32, extraNonce2 uint64, gpb uint32) Share {
	bits := new(uint32)
	*bits = gpb
	return Share{
		Time:               time,
		Nonce:              nonce,
		ExtraNonce2:        extraNonce2,
		GeneralPurposeBits: bits}
}

// A Solution solves a Puzzle. Puzzle + Solution = Proof
type Solution struct {
	ExtraNonce1 uint32
	Share       Share
}

type Proof struct {
	Puzzle   Puzzle
	Solution Solution
}

// the metadata corresponds to the coinbase transaction. In a general
// work puzzle, it can contain any information.
func (p *Proof) Metadata() []byte {
	b := len(p.Puzzle.CoinbaseBegin)
	metadata := make([]byte, b+len(p.Puzzle.CoinbaseEnd)+12)

	copy(metadata, p.Puzzle.CoinbaseBegin)

	binary.BigEndian.PutUint32(metadata[b:], p.Solution.ExtraNonce1)
	binary.BigEndian.PutUint64(metadata[b+4:], p.Solution.Share.ExtraNonce2)

	copy(metadata[b+12:], p.Puzzle.CoinbaseEnd)

	return metadata
}

func (p *Proof) MerkleRoot() []byte {
	return bc.BuildMerkleRootFromCoinbase(p.Metadata(), p.Puzzle.Candidate.MerklePath)
}

func (p *Proof) Blockheader() (*bc.BlockHeader, error) {
	var version uint32
	if p.Puzzle.VersionMask == nil && p.Solution.Share.GeneralPurposeBits == nil {
		version = p.Puzzle.Candidate.Version
	} else if p.Puzzle.VersionMask != nil && p.Solution.Share.GeneralPurposeBits != nil {
		version = (p.Puzzle.Candidate.Version & *p.Puzzle.VersionMask) | (*p.Solution.Share.GeneralPurposeBits & ^(*p.Puzzle.VersionMask))
	} else {
		return nil, errors.New("Inconsistent use of BIP 320")
	}

	return &bc.BlockHeader{
		Version:        version,
		Time:           p.Solution.Share.Time,
		Nonce:          p.Solution.Share.Nonce,
		HashPrevBlock:  p.Puzzle.Candidate.Digest,
		HashMerkleRoot: p.MerkleRoot(),
		Bits:           p.Puzzle.Candidate.Bits}, nil
}

func (p *Proof) Valid() bool {
	b, err := p.Blockheader()
	if err != nil {
		return false
	}
	return b.Valid()
}
