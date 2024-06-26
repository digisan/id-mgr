package id

import (
	"fmt"
	"testing"

	lk "github.com/digisan/logkit"
)

func TestAddAlias(t *testing.T) {

	lk.FailOnErr("%v", BuildHierarchy("", "a", "b"))
	lk.FailOnErr("%v", CreateOneDescWithAlias("", "A", "B"))

	lk.FailOnErr("%v", BuildStandalone("sa", "sb"))
	lk.FailOnErr("%v", CreateOneStdalWithAlias("SA", "SB"))

	// lk.FailOnErr("%v", ID(1).AddAlias("abc"))
	// lk.FailOnErr("%v", AddAlias("abc", "ABC", "AABBCC"))
	// lk.FailOnErr("%v", CreateOneDescWithAlias("", "abcd", "ABCD"))
	// lk.FailOnErr("%v", CreateOneDescWithAlias("", "abcdef"))
	// lk.FailOnErr("%v", CreateOneDescWithAlias("abcd", "abcd1", "abcd2"))
	// lk.FailOnErr("%v", CreateOneStdalWithAlias("SA", "SB", "SC"))
	// lk.FailOnErr("%v", CreateOneStdalWithAlias("SA1", "SB1", "SC1"))
	PrintAlias()

	// lk.FailOnErr("%v", ChangeAlias("abcdef", "ABCDEF"))
	// PrintAlias()

	// lk.FailOnErr("%v", RmAlias("abc", "AABBCC"))
	// PrintAlias()

	// fmt.Println(ID(4).DefaultAlias())
	// fmt.Println(FetchDefaultAlias("ABC"))
}

func TestCleanupAlias(t *testing.T) {

	lk.FailOnErr("%v", BuildHierarchy("", "a", "b"))
	lk.FailOnErr("%v", BuildHierarchy("a", "aa"))
	lk.FailOnErr("%v", BuildHierarchy("aa", "aa0"))
	lk.FailOnErr("%v", BuildHierarchy("aa", "aa1"))
	lk.FailOnErr("%v", CreateOneDescWithAlias("", "DA", "DB"))
	lk.FailOnErr("%v", CreateOneDescWithAlias("", "DA", "DB"))

	fmt.Println("MightBeAncestorOf:", MightBeAncestorOf("a", "aa0"))
	fmt.Println("MightBeParentOf:", MightBeParentOf("a", "aa1"))
	fmt.Println("MightBeSiblingOf:", MightBeSiblingOf("aa0", "aa1"))

	lk.FailOnErr("%v", BuildStandalone("sa", "sb"))
	lk.FailOnErr("%v", CreateOneStdalWithAlias("SA", "SB"))
	lk.FailOnErr("%v", CreateOneStdalWithAlias("SA", "SB"))

	PrintAlias()

	cleanupAlias()

	PrintAlias()
}
