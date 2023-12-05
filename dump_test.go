package idmgr

import (
	"fmt"
	"testing"

	. "github.com/digisan/id-mgr/id"
	lk "github.com/digisan/logkit"
)

func init() {
	fmt.Println("init...")
	lk.FailOnErr("%v", Init(4, 5, 7, 8, 18, 6, 8, 5, 3))

	lk.FailOnErr("%v", BuildHierarchy("", "L0_1", "L0_2", "L0_3"))
	lk.FailOnErr("%v", CreateOneDescWithAlias("L0_1", "L01_1", "L01_2"))
	lk.FailOnErr("%v", CreateOneDescWithAlias("L01_1", "L01_1_1", "L01_1_2"))
	lk.FailOnErr("%v", CreateOneDescWithAlias("L01_1_1", "L01_1_3", "L01_1_4"))
	lk.FailOnErr("%v", CreateOneDescWithAlias("L0_3", "L03_1", "L03_2"))
	lk.FailOnErr("%v", BuildStandalone("S1", "S2", "S3"))
	lk.FailOnErr("%v", CreateOneStdalWithAlias("S4", "S5"))
}

func TestDump(t *testing.T) {
	GenIDTree(true)
	DumpIDTree("h2.txt")
}
