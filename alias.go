package idmgr

import (
	"fmt"

	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
)

func AddAliases(self any, aliases ...any) error {
	id := SearchIDByAlias(self)
	_, err := id.AddAliases(aliases)
	return err
}

func FetchAliases(self any) []any {
	id := SearchIDByAlias(self)
	return id.Alias()
}

func RmAliases(self any, aliases ...any) error {
	id := SearchIDByAlias(self)
	_, err := id.RmAliases(aliases...)
	return err
}

func ChangeAlias(old, new any) error {
	if err := AddAliases(old, new); err != nil {
		return err
	}
	if err := RmAliases(new, old); err != nil {
		return err
	}
	return nil
}

var (
	exclChars = []string{"^", "|", ":", "[", "]"}
)

func validateAlias(alias any) bool {
	return !strs.ContainsAny(fmt.Sprint(alias), exclChars...)
}

func aliasOccupied(alias any, byIDs ...ID) (bool, ID) {
	if len(byIDs) == 0 {
		byIDs = WholeIDs()
	}
	for _, desc := range byIDs {
		if In(alias, desc.Alias()...) {
			return true, desc
		}
	}
	return false, 0
}

// check alias conflict
func CheckAlias(aliases []any, fromIDs ...ID) error {
	if len(fromIDs) == 0 {
		fromIDs = WholeIDs()
	}
	for _, alias := range aliases {
		if !validateAlias(alias) {
			return fmt.Errorf("'%v' contains invalid characters like %+v", alias, exclChars)
		}
		if used, byId := aliasOccupied(alias, fromIDs...); used {
			return fmt.Errorf("'%v' is already used by [%x]", alias, byId)
		}
	}
	return nil
}

func SearchIDByAlias(alias any, fromIDs ...ID) ID {
	if len(fromIDs) == 0 {
		fromIDs = WholeIDs()
	}
	for _, id := range fromIDs {
		if In(alias, id.Alias()...) {
			return id
		}
	}
	return 0
}
