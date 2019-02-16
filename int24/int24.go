package int24

type Int24 [3]byte

func (i24 Int24) ToInt() int {
	return int(i24[2]) | int(i24[1])<<8 | int(i24[0])<<16
}

func New(i int) Int24 {
	var out Int24
	out[0], out[1], out[2] = byte(i), byte(i>>8), byte(i>>16)
	return out
}
