package work

// bip320 was designed for ASICBoost and it specifies some general purpose bits
// in the Version field that miners can use for whatever they want.
const Bip320Mask uint32 = 0xe0001fff
const Bip320GeneralPurposeBits uint32 = ^Bip320Mask

func Version(u uint32, gpb uint32) uint32 {
	return (u & Bip320Mask) | (gpb & Bip320GeneralPurposeBits)
}
