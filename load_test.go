package idmgr

import (
	"fmt"
	"testing"

	. "github.com/digisan/id-mgr/id"
)

func TestLoad(t *testing.T) {
	fmt.Printf("IngestTree err: %v\n", IngestTree("./h2.txt"))
	GenIDTree(true)
	fmt.Println(ID(0).DescendantsCount(true))
	fmt.Println(ID(0).ChildrenCount())
	fmt.Println(ID(MaxID).ChildrenCount())

	fmt.Println("---------------------------------------------------")
	fmt.Printf("IngestTree err: %v\n", Transplant(ID(0x11), ID(0x13)))
	GenIDTree(true)
}
