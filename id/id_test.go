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
	lk.FailOnErr("%v", Init(4, 3, 2, 7, 8, 18, 6, 8, 5, 3))

	id, err := SetID(ID(1))
	fmt.Printf("%x, %v\n", id, err)
	id, err = SetID(ID(2))
	fmt.Printf("%x, %v\n", id, err)
	id, err = SetID(ID(3))
	fmt.Printf("%x, %v\n", id, err)

	id, err = ID(0).GenDescID()
	fmt.Printf("%x, %v\n", id, err)
	id, err = ID(0).GenDescID()
	fmt.Printf("%x, %v\n", id, err)
	id, err = ID(0).GenDescID()
	fmt.Printf("%x, %v\n", id, err)

	fmt.Println("------------------------")
	fmt.Println(HierarchyIDs())
	fmt.Println(StandaloneIDs())
	fmt.Println(WholeIDs())
	fmt.Println("------------------------")

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
	lk.FailOnErr("%v", Init(4, 3, 2, 7, 8, 18, 6, 8, 5, 3))

	id, err := ID(0).GenDescID()
	fmt.Printf("%x, %v\n", id, err)
	id, err = ID(id).GenDescID()
	fmt.Printf("%x, %v\n", id, err)
	id, err = ID(id).GenDescID()
	fmt.Printf("%x, %v\n", id, err)

	id, err = ID(0).GenDescID()
	fmt.Printf("%x, %v\n", id, err)
	id, err = ID(id).GenDescID()
	fmt.Printf("%x, %v\n", id, err)

	fmt.Println("------------------------")
	fmt.Println(HierarchyIDs())
	fmt.Println(StandaloneIDs())
	fmt.Println(WholeIDs())
	fmt.Println("------------------------")

	// fmt.Println(CopyBranchAt(2, 0x2, 0x12))
	// fmt.Println(TransplantAt(1, 0x12))

	fmt.Println("------------------------")
	fmt.Println(HierarchyIDs())
	fmt.Println(StandaloneIDs())
	fmt.Println(WholeIDs())
	fmt.Println("------------------------")
}

func TestDescTree(t *testing.T) {
	lk.FailOnErr("%v", Init(4, 3, 2, 7, 8, 18, 6, 8, 5, 3))

	fmt.Println(ID(131746).Part(0))
	fmt.Println(ID(131746).Part(1))
	fmt.Println(ID(131746).Part(2))
	fmt.Println(ID(131746).Part(3))
	fmt.Println(ID(131746).Part(4))
	return

	id, err := ID(0).GenDescID()
	fmt.Printf("%x, %v\n", id, err)
	id, err = ID(id).GenDescID()
	fmt.Printf("%x, %v\n", id, err)
	id, err = ID(id).GenDescID()
	fmt.Printf("%x, %v\n", id, err)
	id, err = ID(id).GenDescID()
	fmt.Printf("%x, %v\n", id, err)
	id, err = ID(id).GenDescID()
	fmt.Printf("%x, %v\n", id, err)

	id, err = ID(0).GenDescID()
	fmt.Printf("%x, %v\n", id, err)

	aid, err := ID(0).AvailableDescID()
	fmt.Printf("available: %x %v\n", aid, err)

	id, err = ID(0).GenDescID()
	fmt.Printf("%x, %v\n", id, err)

	id, err = ID(2).GenDescID()
	fmt.Printf("%x, %v\n", id, err)

	fmt.Println("------------------------")
	fmt.Println(HierarchyIDs())
	fmt.Println(StandaloneIDs())
	fmt.Println(WholeIDs())
	fmt.Println("------------------------")

	// fmt.Println(CopyBranchAt(2, 0x2, 0x12))
	// fmt.Println(TransplantAt(1, 0x12))

	// if id, nShiftedSeg, err := ID(17).descAsTree(); err == nil {
	// 	fmt.Println(id, nShiftedSeg)
	// }

	// id, err = ID(17).HookDesc(1)
	// fmt.Printf("%x, %v\n", id, err)

	nid, err := ID(2).TransplantBranch(ID(17))
	fmt.Printf("%v, --- %v\n", nid, err)

	fmt.Println("------------------------")
	fmt.Println(HierarchyIDs())
	fmt.Println(StandaloneIDs())
	fmt.Println(WholeIDs())
	fmt.Println("------------------------")
}
