package id

import "testing"

func TestInit(t *testing.T) {
	Init64bitsFromStr("[000000000000000f 00000000000001f0 000000000000fe00 0000000000ff0000 000003ffff000000 0000fc0000000000 00ff000000000000 1f00000000000000 e000000000000000]")
	PrintMasks(true)
	PrintSegs(true)
	PrintCapLevel(true)
	PrintCapStdal(true)
}
