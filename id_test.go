package idmgr

import (
	"fmt"
	"testing"
)

func TestID(t *testing.T) {

	if err := Init64bits(4, 2, 10, 8, 18, 6, 8, 5, 3); err != nil {
		fmt.Println(err)
		return
	} else {
		// fmt.Printf("%064b\n", masks)
		// fmt.Printf("%064b\n", segs)
	}

	//

	// fmt.Println(ID(20))

	fmt.Println(_cap_std)

	fmt.Println(_cap_lvl[0])
	fmt.Println(_cap_lvl[1])
	fmt.Println(_cap_lvl[2])
	fmt.Println(_cap_lvl[3])

	fmt.Println("------------------------")

	fmt.Println(ID(1))
	fmt.Println(ID(15))
	fmt.Println(ID(16))
	fmt.Println(ID(17))
	fmt.Println(ID(18))
	fmt.Println(ID(19))
}
