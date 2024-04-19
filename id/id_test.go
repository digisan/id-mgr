package id

import (
	"fmt"
	"testing"

	lk "github.com/digisan/logkit"
)

func init() {
	fmt.Println("init...")
	lk.FailOnErr("%v", Init(4, 5, 7, 8, 18, 6, 8, 5, 3))

	id, err := ID(0).GenDescID()
	lk.FailOnErr("%v", err)
	id, err = ID(id).GenDescID()
	lk.FailOnErr("%v", err)
	id, err = ID(id).GenDescID()
	lk.FailOnErr("%v", err)
	id, err = ID(id).GenDescID()
	lk.FailOnErr("%v", err)
	_, err = ID(id).GenDescID()
	lk.FailOnErr("%v", err)
	id, err = ID(17).GenDescID()
	lk.FailOnErr("%v", err)
	_, err = ID(id).GenDescID()
	lk.FailOnErr("%v", err)
	_, err = ID(0).GenDescID()
	lk.FailOnErr("%v", err)
	_, err = ID(0).GenDescID()
	lk.FailOnErr("%v", err)
	_, err = ID(0).GenDescID()
	lk.FailOnErr("%v", err)
	_, err = ID(2).GenDescID()
	lk.FailOnErr("%v", err)

	_, err = GenStdalID()
	lk.FailOnErr("%v", err)
	_, err = GenStdalID()
	lk.FailOnErr("%v", err)
	_, err = GenStdalID()
	lk.FailOnErr("%v", err)

	fmt.Println("------------------------")
	fmt.Println(HierarchyIDs())
	fmt.Println(StandaloneIDs())
	fmt.Println(WholeIDs())
	fmt.Println("------------------------")
}

func TestClearAllID(t *testing.T) {

	// _, err := DeleteID(MaxID, false)
	// lk.FailOnErr("%v", err)

	// lk.FailOnErr("%v", ClrAllID())

	fmt.Println(HierarchyIDs())
	fmt.Println(StandaloneIDs())
	fmt.Println(WholeIDs())

	fmt.Println(ID(17).IsAncestorOf(16843281))
	fmt.Println(ID(17).IsParentOf(16843281))
	fmt.Println(ID(17).IsParentOf(529))
	fmt.Println(ID(16843281).IsDescendantOf(17))
	fmt.Println(ID(16843281).IsDescendantOf(18))
	fmt.Println(ID(1).IsSiblingOf(ID(2)))
}

func TestDescendants(t *testing.T) {

	fmt.Println(ID(16843281).Ancestors(true))
	fmt.Println(ID(16843281).Ancestors(false))
	fmt.Println(ID(0).Descendants(100, false))
	fmt.Println(ID(0).Descendants(100, true))

	return

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

	// lk.FailOnErr("%v", SetID(1))
	// lk.FailOnErr("%v", SetID(2))
	// lk.FailOnErr("%v", SetID(3))
	// lk.FailOnErr("%v", SetID(16))
	// lk.FailOnErr("%v", SetID(17))
	// lk.FailOnErr("%v", SetID(18))
	// lk.FailOnErr("%v", SetID(19))
	// lk.FailOnErr("%v", SetID(20))

	for i := 1; i < 20000; i++ {
		// fmt.Printf("inserting...(0x%03x)\n", i)

		// if i == 0x90 {
		// 	fmt.Println("DEBUGGING...")
		// }

		id, err := SetID(ID(i))
		lk.FailOnErr("%x, %v", id, err)
	}

	// fmt.Println("------------------------")
	// fmt.Println(HierarchyIDs())
	// fmt.Println(StandaloneIDs())
	// fmt.Println(WholeIDs())
}

func TestLeftShift(t *testing.T) {

	ids, err := DeleteID(1, false)
	fmt.Printf("%x, %v\n", ids, err)

	ids, err = DeleteID(3, false)
	fmt.Printf("%x, %v\n", ids, err)

	ids, err = DeleteID(5, false)
	fmt.Printf("%x, %v\n", ids, err)

	mRecord.Range(func(key, value any) bool {
		fmt.Println("-->", key, value)
		return true
	})

	// fmt.Println(id.leftShift(1, true))

	fmt.Println("------------------------")
	fmt.Println(HierarchyIDs())
	fmt.Println(StandaloneIDs())
	fmt.Println(WholeIDs())
	fmt.Println("------------------------")
}

func TestCanBeBranch(t *testing.T) {

	// fmt.Println(CopyBranchAt(2, 0x2, 0x12))
	// fmt.Println(TransplantAt(1, 0x12))

	fmt.Println("------------------------")
	fmt.Println(HierarchyIDs())
	fmt.Println(StandaloneIDs())
	fmt.Println(WholeIDs())
	fmt.Println("------------------------")
}

func TestDescTree(t *testing.T) {

	aid, err := ID(0).AvailableDescID()
	fmt.Printf("available: %x %v\n", aid, err)

	// lk.FailOnErr("%v", CopyBranch(17, 2))
	lk.FailOnErr("%v", Transplant(17, 2))

	fmt.Println("------------------------")
	fmt.Println(HierarchyIDs())
	fmt.Println("------------------------")
}
