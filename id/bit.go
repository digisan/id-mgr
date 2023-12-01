package id

func lowBits(n uint64, nLow uint) uint64 {
	shift := 64 - nLow
	checker := F16 >> uint64(shift)
	return n & checker
}

func count1(n uint64) uint {
	count := uint(0)
	for i := 0; i < 64; i++ {
		flag1 := uint64(0b01 << i)
		if flag1&n == flag1 {
			count++
		}
	}
	return count
}

func countF(n uint64) uint {
	count := uint(0)
	for i := 0; i < 64; i += 4 {
		flagF := uint64(0x0F << i)
		if flagF&n == flagF {
			count++
		}
	}
	return count
}

func trimLow0(s uint64, bitStep uint8) uint64 {
	var bitChecker = F16 >> (64 - bitStep)
	for {
		if bitChecker&s != 0 {
			return s
		}
		s = s >> bitStep
	}
}

func trimLowB0(s uint64) uint64 {
	return trimLow0(s, 1)
}

func trimLowO0(s uint64) uint64 {
	return trimLow0(s, 3)
}

func trimLowX0(s uint64) uint64 {
	return trimLow0(s, 4)
}
