package idmgr

import (
	"fmt"
	"testing"

	lk "github.com/digisan/logkit"
)

func TestInit64bits(t *testing.T) {
	if err := Init64bits(1, 7, 12, 4, 18, 6, 8, 4, 4); err != nil {
		fmt.Println(err)
	}

	fmt.Println(maxDescCap(0))
	fmt.Println(maxDescCap(1))
	fmt.Println(maxDescCap(2))
	fmt.Println(maxDescCap(3))
	fmt.Println(maxDescCap(4))
	fmt.Println(maxDescCap(5))
	fmt.Println(maxDescCap(6))
	fmt.Println(maxDescCap(7))
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
	fmt.Println("availableSegBitIdx:", id.availableSegBitIdx())

	id = 10
	fmt.Println("ID(10) level:", id.level())
	fmt.Println("availableSegBitIdx:", id.availableSegBitIdx())

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

		if i == 1 {
			_, err = id.AddAlias("A", "B")
			lk.FailOnErr("%v", err)
		}
		id.RmAlias("AA")
		fmt.Println(id.Alias())

		for i := 0; i < N; i++ {
			id, err := GenID(id)
			lk.FailOnErr("%d", err)
			fmt.Println("   L1:", id, id.Ancestors())

			for i := 0; i < N; i++ {
				id, err := GenID(id)
				lk.FailOnErr("%d", err)
				fmt.Println("        L2:", id, id.Ancestors())

				for i := 0; i < N; i++ {
					id, err := GenID(id)
					lk.FailOnErr("%d", err)
					fmt.Println("            L3:", id, id.Ancestors())

					for i := 0; i < N; i++ {
						id, err := GenID(id)
						lk.FailOnErr("%d", err)
						fmt.Println("                L4:", id, id.Ancestors())
					}
				}
			}
		}
	}

	fmt.Println("-------------------------------")

	// id = ID(0)
	// descendants := id.Descendants(1)
	// fmt.Println(len(descendants), descendants)

	// fmt.Println(SearchIDByAlias("A"))
}
