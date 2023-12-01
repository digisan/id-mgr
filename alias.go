package idmgr

import (
	"fmt"
	"sync"

	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
	"github.com/digisan/id-mgr/id"
)

var (
	// alias: key id, value: aliases
	mAlias = sync.Map{} // make(map[ID][]any)
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

func aliasOccupied(alias any, byIDs ...id.ID) (bool, id.ID) {
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
func CheckAlias(aliases []any, fromIDs ...id.ID) error {
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

func SearchIDByAlias(alias any, fromIDs ...id.ID) id.ID {
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

func (id ID) Alias() []any {
	if id == MaxID {
		return []any{"standalone"}
	}
	if typ := id.Type(); typ == ID_HRCHY_ALLOC || typ == ID_STDAL_ALLOC {
		if v, ok := mAlias.Load(id); ok {
			return v.([]any)
		}
	}
	return []any{}
}

func (id ID) AddAliases(aliases []any, validRangeIDs ...ID) ([]any, error) {

	if typ := id.Type(); typ == ID_HRCHY_ALLOC || typ == ID_STDAL_ALLOC {

		// check alias conflict
		if err := CheckAlias(aliases, validRangeIDs...); err != nil {
			return id.Alias(), err
		}

		mAlias[id] = append(mAlias[id], aliases...)
		mAlias[id] = Settify(mAlias[id]...)
		return id.Alias(), nil
	}

	return nil, fmt.Errorf("error: %v doesn't exist, cannot do AddAlias", id)
}

// func (id ID) RmAliases(aliases ...any) ([]any, error) {
// 	if !id.Exists() {
// 		return nil, fmt.Errorf("error: %v doesn't exist, cannot do RmAlias", id)
// 	}
// 	mAlias[id] = Filter(id.Alias(), func(i int, e any) bool {
// 		return NotIn(e, aliases...)
// 	})
// 	return id.Alias(), nil
// }
