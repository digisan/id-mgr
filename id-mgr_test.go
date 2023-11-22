package idmgr

import (
	"fmt"
	"testing"

	lk "github.com/digisan/logkit"
)

// 1111111111111111111111111111111111111111111111111111111111111111

func TestInit64bits(t *testing.T) {
	if err := Init64bits(4, 2, 14, 4, 18, 6, 8, 4, 4); err != nil {
		fmt.Println(err)
		return
	} else {
		// fmt.Printf("%064b\n", masks)
		// fmt.Printf("%064b\n", segs)
	}

	fmt.Println(ID(0).selfStartBitIdx(), ID(0).descAvailableBitIdx(), ID(0).level(), nextLvlDescCap(ID(0).level()))
	fmt.Println(ID(1).selfStartBitIdx(), ID(1).descAvailableBitIdx(), ID(1).level(), nextLvlDescCap(ID(1).level()))
	fmt.Println(ID(2).selfStartBitIdx(), ID(2).descAvailableBitIdx(), ID(2).level(), nextLvlDescCap(ID(2).level()))
	fmt.Println(ID(3).selfStartBitIdx(), ID(3).descAvailableBitIdx(), ID(3).level(), nextLvlDescCap(ID(3).level()))
	fmt.Println(ID(16).selfStartBitIdx(), ID(16).descAvailableBitIdx(), ID(16).level(), nextLvlDescCap(ID(16).level()))
	fmt.Println(ID(63).selfStartBitIdx(), ID(63).descAvailableBitIdx(), ID(63).level(), nextLvlDescCap(ID(63).level()))
	fmt.Println(ID(64).selfStartBitIdx(), ID(64).descAvailableBitIdx(), ID(64).level(), nextLvlDescCap(ID(64).level()))
	fmt.Println(ID(64000).selfStartBitIdx(), ID(64000).descAvailableBitIdx(), ID(64000).level(), nextLvlDescCap(ID(64000).level()))

	fmt.Println()

	fmt.Println(nextLvlDescCap(0))
	fmt.Println(nextLvlDescCap(1))
	fmt.Println(nextLvlDescCap(2))
	fmt.Println(nextLvlDescCap(3))
	fmt.Println(nextLvlDescCap(4))
	fmt.Println(nextLvlDescCap(5))
	fmt.Println(nextLvlDescCap(6))
	fmt.Println(nextLvlDescCap(7))
	fmt.Println(nextLvlDescCap(8))
	fmt.Println(nextLvlDescCap(9))
}

func TestID1(t *testing.T) {
	if err := Init64bits(5, 15, 12, 8, 8, 8, 4, 4); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%016b\n%016b\n", 0x10, trimLowZeroBin(0x10))
	fmt.Printf("%016b\n%016b\n", 0x10, trimLowZeroOct(0x10))
	fmt.Printf("%016b\n%016b\n", 0x10, trimLowZeroHex(0x10))

	fmt.Println("makeID(0, 1):", makeID(0, 1))

	for i := 0; i < 13; i++ {
		id, err := GenID(0)
		if err == nil {
			fmt.Println("id:", id)
		}
	}

	var id ID = 0
	fmt.Println("ID(0) level:", id.level())
	fmt.Println("descAvailableBitIdx:", id.descAvailableBitIdx())

	id = 10
	fmt.Println("ID(10) level:", id.level())
	fmt.Println("descAvailableBitIdx:", id.descAvailableBitIdx())

	id = makeID(12, 1)
	fmt.Printf("makeID(12, 1): %x\n", id)
	fmt.Printf("makeID(12, 1): %d\n", id)
	fmt.Println("id Ancestors:", id.Ancestors())

	fmt.Printf("---> GenID under super class id [%x]:\n", id)
	id, err := GenID(id)
	if err == nil {
		fmt.Println(id)
	} else {
		fmt.Println(err)
	}

	fmt.Printf("---> GenID under super class id [%x]:\n", 12)
	id, err = GenID(12)
	if err == nil {
		fmt.Println(id)
	} else {
		fmt.Println(err)
	}
}

func TestID2(t *testing.T) {

	if err := Init64bits(12, 12, 12, 8, 8, 8, 4); err != nil {
		fmt.Println(err)
		return
	}

	const N = 2

	for i := 0; i < N; i++ {

		id, err := GenID(0)
		lk.FailOnErr("%v", err)
		fmt.Println("Level:", id.level(), id, id.Ancestors())

		// lk.FailOnErr("%d", DelID(id))

		if i == 0 {
			_, err = id.AddAliases([]any{"A", "B"})
			lk.FailOnErr("%v", err)
		}

		_, err = id.RmAliases("AA")
		lk.FailOnErr("%v", err)
		fmt.Println(id.Alias())

		for i := 0; i < N; i++ {
			id, err := GenID(id)
			lk.FailOnErr("%v", err)
			fmt.Println("   Level:", id.level(), id, id.Ancestors())

			for i := 0; i < N; i++ {
				id, err := GenID(id)
				lk.FailOnErr("%v", err)
				fmt.Println("        Level:", id.level(), id, id.Ancestors())

				for i := 0; i < N; i++ {
					id, err := GenID(id)
					lk.FailOnErr("%v", err)
					fmt.Println("            Level:", id.level(), id, id.Ancestors())

					for i := 0; i < N; i++ {
						id, err := GenID(id)
						lk.FailOnErr("%v", err)
						fmt.Println("                Level:", id.level(), id, id.Ancestors())
					}
				}
			}
		}
	}

	fmt.Println("-------------------------------")

}

func TestAlias(t *testing.T) {

	if err := Init64bits(4, 2, 14, 4, 18, 6, 8, 4, 4); err != nil {
		fmt.Println(err)
		return
	} else {
		// fmt.Printf("%064b\n", masks)
		// fmt.Printf("%064b\n", segs)
	}

	for i := 1; i <= 10; i++ {
		id, err := GenID(0)
		lk.FailOnErr("%v", err)
		fmt.Println(id)
		id.AddAliases([]any{i * 10, i * 100})

		for j := 1; j <= 10; j++ {
			id, err := GenID(ID(i))
			lk.FailOnErr("%v", err)
			fmt.Println(id)
			id.AddAliases([]any{i * 10, i * 100})
		}
	}

	descendants := ID(5).Descendants(100)
	fmt.Println(descendants, len(descendants))

	fmt.Println(SearchIDByAlias(2000))
	fmt.Println(SearchIDByAlias(90))

	fmt.Println(WholeIDs())
}

func TestBuildHierarchy(t *testing.T) {

	if err := Init64bits(4, 4, 12, 4, 18, 6, 8, 4, 4); err != nil {
		fmt.Println(err)
		return
	} else {
		// fmt.Printf("%064b\n", masks)
		// fmt.Printf("%064b\n", segs)
	}

	lk.FailOnErr("%v", BuildHierarchy("", "C1", "C 3", "C  3"))
	lk.FailOnErr("%v", BuildHierarchy("C1", "C12"))
	lk.FailOnErr("%v", BuildHierarchy("", "C2"))
	lk.FailOnErr("%v", BuildHierarchy("C2", "C13"))
	lk.FailOnErr("%v", BuildHierarchy("", "C14"))
	lk.FailOnErr("%v", BuildHierarchy("C14", "C141"))
	lk.FailOnErr("%v", BuildHierarchy("C12", "C121"))
	lk.FailOnErr("%v", BuildHierarchy("C12", "C122"))
	lk.WarnOnErr("%v", BuildHierarchy("C121", "C1211", "C1212", "C1213"))
	lk.WarnOnErr("%v", BuildHierarchy("C1213", "C1213-1", "C1213-2"))

	AddAliases("C1213-2", "C1213-2X", "C1213-2Y", "C1213-2Z")

	lk.WarnOnErr("%v", BuildHierarchy("C1213-2Z", "2ZZZ", "2XYZ"))

	AddAliases("C121", "c121", "CC121")

	// RmAliases("C121", "C121")
	// fmt.Println("--->", GetAliases("C121"))

	// RmAliases("C122", "C122")
	// fmt.Println("--->", GetAliases("C122"))

	fmt.Println("-----------------------------")

	GenHierarchy(true)

	fmt.Println("-----------------------------")

	lk.WarnOnErr("%v", DelIDsOnAlias("C1211", "C1212", "C1213"))
	// fmt.Println(WholeIDs())

	GenHierarchy(true)

	fmt.Println("-----------------------------")

	DumpHierarchy("dump.txt")
}

func TestIngestHierarchy(t *testing.T) {

	// if err := Init64bits(4, 4, 12, 4, 18, 6, 8, 4, 4); err != nil {
	// 	fmt.Println(err)
	// 	return
	// } else {
	// 	// fmt.Printf("%016x\n", masks)
	// 	// fmt.Printf("%016x\n", segs)
	// }

	lk.FailOnErr("%v", IngestHierarchy("dump.txt"))

	fmt.Println(ID(0x11).Descendants(100))
	fmt.Println(ID(0x11).Parent())

	fmt.Println("mAlias:", mAlias)
	fmt.Println("mRecord:", mRecord)

	fmt.Println("--------------------------")

	// fmt.Println(ID(0).Descendants(100))
	GenHierarchy(true)
}
