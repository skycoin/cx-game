package common

//bitset of fixed length, length is specified at creation time
type Bitset []byte

func NewBitSet(len int) Bitset {
	newBitSet := make(Bitset, len)
	return newBitSet
}

//pass enum
func (b Bitset) GetBit(i int) byte {
	// return b.bitfield&i != 0
	return b[i]
}

func (b Bitset) SetBit(i int, value bool) {
	if value && b[i] < 255 {
		b[i]++
	} else {
		b[i] = 0
	}
}
