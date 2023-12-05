package idmgr

import (
	"fmt"
	"testing"

	. "github.com/digisan/id-mgr/id"
)

func TestLoad(t *testing.T) {
	fmt.Printf("IngestTree err: %v\n", IngestTree("./h2.txt"))
	GenIDTree(true)
	fmt.Println(ID(0).DescendantsInfo(true))
	// fmt.Println(ID(0).ChildrenCount())
	// fmt.Println(ID(MaxID).ChildrenCount())

	fmt.Println("---------------------------------------------------")
	fmt.Printf("IngestTree err: %v\n", Transplant(ID(0x1), ID(0x2)))
	// fmt.Printf("IngestTree err: %v\n", Transplant(ID(0x11), ID(0x2)))
	// fmt.Printf("IngestTree err: %v\n", Transplant(ID(0x11), ID(0x21)))
	GenIDTree(true)
}

func TestAddAliasEx(t *testing.T) {
	ID(1).AddAlias("A")
	ID(1).AddAlias("A")
	ID(1).AddAlias("A")
	GenIDTree(true)
}
