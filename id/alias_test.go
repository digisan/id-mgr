package id

import (
	"fmt"
	"testing"

	lk "github.com/digisan/logkit"
)

func TestAddAlias(t *testing.T) {

	lk.FailOnErr("%v", ID(1).AddAlias("abc"))
	lk.FailOnErr("%v", AddAlias("abc", "ABC", "AABBCC"))
	lk.FailOnErr("%v", CreateDescWithAlias("", "abcd"))
	lk.FailOnErr("%v", CreateDescWithAlias("", "abcdef"))
	PrintAlias()

	lk.FailOnErr("%v", ChangeAlias("abcdef", "ABCDEF"))
	PrintAlias()

	lk.FailOnErr("%v", RmAlias("abc", "AABBCC"))
	PrintAlias()

	fmt.Println(ID(4).DefaultAlias())
	fmt.Println(FetchDefaultAlias("ABC"))
}
