package gb 

func boolToUint8(b bool) uint8 {
	if b == true {
		return 1
	}
	return 0 
}

func bytesToWord(b1 byte, b2 byte) uint16 {
	return (uint16(b1) << 8) ^ uint16(b2)
}

func wordToBytes(n uint16) (byte, byte) {
	return byte(n >> 8), byte(n & 0x00ff)
}

func sumBytes(a, b byte) (sum byte, overflow bool) {
	sum = a + b
	return sum, a > 0xff - b 
}

func sumUint16(a, b uint16) (sum uint16, overflow bool) {
	sum = a + b
	return sum, a > 0xffff - b 
}

func subBytes(a, b byte) (res byte, underflow bool) {
	underflow = (a < b)
	res = a - b
	return res, underflow
}

func getUpperByte(target uint16) byte {
	return byte((target & 0xff00) >> 8)
}

func getLowerByte(target uint16) byte {
	return byte(target & 0xff)
}

func setLowerByte(target uint16, lo byte) uint16 {
	return (target & 0xff00) | uint16(lo)
}

func setUpperByte(target uint16, hi byte) uint16 {
	return (target & 0x00ff) | (uint16(hi) << 8)
}

func bitEnabled(target *byte, bit int) bool {
	return (*target >> bit) & 1 == 1
}

func setBit(target *byte, bit byte) {
	*target |= (1 << bit)
}
