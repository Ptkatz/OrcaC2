package sha3

func Keccak1600State(st *[25]uint64, data []byte) {
	s := &state{rate: 136}
	s.Write(data)
	s.padAndPermute(0x01)
	*st = s.a
}

func Keccak1600Permute(st *[25]uint64) {
	keccakF1600(st)
}
