package idmgr

import (
	"fmt"
	"testing"
)

func TestID(t *testing.T) {

	if err := Init64bits(4, 3, 2, 7, 8, 18, 6, 8, 5, 3); err != nil {
		fmt.Println(err)
		return
	} else {
		// fmt.Printf("%064b\n", masks)
		// fmt.Printf("%064b\n", segs)
	}

	//

	// fmt.Println(ID(20))

	fmt.Println(_cap_std)
	fmt.Println(_cap_lvl)

	fmt.Println("------------------------")

	ID(0).Print()
	ID(F16).Print()
	ID(1).Print()
	ID(15).Print()
	ID(16).Print()
	ID(63).Print()

	// fmt.Print(ID(63), "  ")
	// fmt.Println(ID(63).Ancestors())

	// fmt.Print(ID(64), "  ")
	// fmt.Println(ID(64).Ancestors())

	// fmt.Print(ID(255), "  ")
	// fmt.Println(ID(255).Ancestors())

	// fmt.Print(ID(256), "  ")
	// fmt.Println(ID(256).Ancestors())

	// fmt.Println(ID(0))
	// fmt.Println(ID(F16))

}
