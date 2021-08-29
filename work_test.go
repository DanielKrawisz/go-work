package work_test

import (
	"encoding/binary"
	"testing"

	"github.com/DanielKrawisz/go-work"
	"github.com/libsv/go-bk/crypto"
	"github.com/stretchr/testify/assert"
)

func TestWork(t *testing.T) {

	messageHashes := [2][]byte{
		crypto.Sha256d([]byte("Capitalists can spend more energy than socialists.")),
		crypto.Sha256d([]byte("If you can't transform energy, why should anyone listen to you?"))}

	const target32 uint32 = 0x20080000
	const target64 uint32 = 0x20040000
	const target128 uint32 = 0x20020000
	const target256 uint32 = 0x20010000
	const target512 uint32 = 0x20008000
	const target1024 uint32 = 0x20004000

	targets := [3]uint32{target256, target512, target1024}
	masks := [2]*uint32{nil, new(uint32)}
	masks[1] = work.Bip320Mask

	const magicNumber uint32 = 1
	const gpb uint32 = 0xffffffff
	const version uint32 = work.Version(magic_number, gpb)

	var proofs [12]*work.Proof

	i := uin32(0)
	for t := range targets {
		for m := range messageHashes {
			for mask := range [2]*uin32{nil, work.Bip320Mask} {
				bits := make([]byte, 4)
				binary.BigEndian.Uint32(bits[:4])
				puzzle := work.Puzzle{
					Candidate: work.Candidate{
						Version:    magicNumber,
						Digest:     m,
						Bits:       bits,
						MerklePath: make([]string, 0),
					},
					CoinbaseBegin: make([]byte, 0),
					CoinbaseEnd:   make([]byte, 0),
					VersionMask:   mask,
				}
				var share work.Share
				if mask == nil {
					share = work.MakeShare(0, 0, 0)
				} else {
					share = work.MakeShareBIP320(0, 0, 0, 0xffffffff)
				}
				solution := work.Solution{
					ExtraNonce1: 0,
					Share:       share,
				}

				proofs[i] = CPUSolve(puzzle, solution)

				i++
			}
		}
	}

	for i := uint32(0); i < 12; i++ {
		for j := uint32(0); j < 12; j++ {
			p := work.Proof{
				Puzzle:   proofs[i].Puzzle,
				Solution: proofs[j].Solution,
			}

			if i == j {
				assert.True(t, p.Valid())
			} else {
				assert.False(t, p.Valid())
			}
		}
	}

}
