package id

import (
	"fmt"
	"testing"

	lk "github.com/digisan/logkit"
)

func TestDescendants(t *testing.T) {

	lk.FailOnErr("%v", Init(4, 3, 2, 7, 8, 18, 6, 8, 5, 3))

	fmt.Println(ID(0).GenDescID())
	fmt.Println(ID(0).GenDescID())
	fmt.Println(ID(0).GenDescID())
	fmt.Println(ID(0).GenDescID())
	fmt.Println(ID(0).GenDescID())

	fmt.Println(ID(1).GenDescID())
	fmt.Println(ID(17).GenDescID())
	fmt.Println(ID(145).GenDescID())
	fmt.Println(ID(657).GenDescID())

	fmt.Println(ID(1).GenDescID())
	fmt.Println(ID(1).GenDescID())

	fmt.Println(ID(2).GenDescID())
	fmt.Println(ID(3).GenDescID())

	fmt.Println(ID(9).GenDescID())

	fmt.Println(ID(657).AncestorsWithSelf())
	fmt.Println(ID(0).Descendants(100))
	fmt.Println(ID(0).DescendantsWithSelf(100))

	fmt.Println(ID(1).Parent())
	fmt.Println(ID(657).Parent())
	fmt.Println(ID(658).Parent())

	ID(0).PrintDescendants(100, true)
	ID(0).PrintDescendants(100, false)

	fmt.Println("------------------------")

	fmt.Println(MaxID.GenDescID())
	fmt.Println(GenStdalID())
	fmt.Println(GenStdalID())
	fmt.Println(GenStdalID())
	fmt.Println(GenStdalID())
	fmt.Println(GenStdalID())

	fmt.Println("------------------------")

	fmt.Println(HierarchyIDs())
	fmt.Println(StandaloneIDs())
	fmt.Println(WholeIDs())

	fmt.Println(DeleteID(80, false))
	fmt.Println(DeleteID(1, true))

	fmt.Println("------------------------")

	fmt.Println(HierarchyIDs())
	fmt.Println(StandaloneIDs())
	fmt.Println(WholeIDs())
}

func TestSetID(t *testing.T) {
	lk.FailOnErr("%v", Init(4, 3, 2, 7, 8, 18, 6, 8, 5, 3))

	// lk.FailOnErr("%v", SetID(1))
	// lk.FailOnErr("%v", SetID(2))
	// lk.FailOnErr("%v", SetID(3))
	// lk.FailOnErr("%v", SetID(16))
	// lk.FailOnErr("%v", SetID(17))
	// lk.FailOnErr("%v", SetID(18))
	// lk.FailOnErr("%v", SetID(19))
	// lk.FailOnErr("%v", SetID(20))

	for i := 1; i < 50000; i++ {
		// fmt.Printf("inserting...(0x%03x)\n", i)

		// if i == 0x90 {
		// 	fmt.Println("DEBUGGING...")
		// }

		lk.FailOnErr("%v", SetID(ID(i)))
	}

	fmt.Println("------------------------")

	fmt.Println(HierarchyIDs())
	fmt.Println(StandaloneIDs())
	fmt.Println(WholeIDs())
}
