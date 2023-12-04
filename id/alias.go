package id

import (
	"fmt"
	"sync"

	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
)

var (
	// mAlias: key id, value: aliases
	mAlias    = sync.Map{} // make(map[ID][]any)
	exclChars = []string{"^", "|", ":", "[", "]"}
)

func validateAlias(alias any) bool {
	return !strs.ContainsAny(fmt.Sprint(alias), exclChars...)
}

func aliasOccupied(alias any, byIDs ...ID) (bool, ID) {
	if len(byIDs) == 0 {
		byIDs = WholeIDs()
	}
	for _, id := range byIDs {
		if In(alias, id.Alias()...) {
			return true, id
		}
	}
	return false, 0
}

func CheckAlias(aliases []any, fromIDs ...ID) error {
	if len(fromIDs) == 0 {
		fromIDs = WholeIDs()
	}
	for _, alias := range aliases {
		if !validateAlias(alias) {
			return fmt.Errorf("'%v' contains invalid characters(%+v)", alias, exclChars)
		}
		if used, byId := aliasOccupied(alias, fromIDs...); used {
			return fmt.Errorf("'%v' is already used by [%x]", alias, byId)
		}
	}
	return nil
}

func (id ID) Alias() []any {
	if id == MaxID {
		return []any{ID_STDAL_ROOT.String()}
	}
	if id == 0 {
		return []any{ID_HRCHY_ROOT.String()}
	}
	if In(id.Type(), ID_HRCHY_ALLOC, ID_STDAL_ALLOC) {
		if v, ok := mAlias.Load(id); ok {
			return v.([]any)
		}
	}
	return []any{}
}

func (id ID) DefaultAlias() any {
	aliases := id.Alias()
	if len(aliases) > 0 {
		return aliases[0]
	}
	return fmt.Sprintf("%x", id)
}

func (id ID) AddAlias(aliases ...any) error {
	if NotIn(id.Type(), ID_HRCHY_ALLOC, ID_STDAL_ALLOC) {
		return fmt.Errorf("error: %v doesn't exist, cannot do AddAlias", id)
	}
	if err := CheckAlias(aliases); err != nil {
		return err
	}
	if ea, ok := mAlias.Load(id); ok {
		aliases = append(ea.([]any), aliases...)
		mAlias.Store(id, aliases)
	} else {
		mAlias.Store(id, aliases)
	}
	return nil
}

func (id ID) RmAlias(aliases ...any) error {
	if NotIn(id.Type(), ID_HRCHY_ALLOC, ID_STDAL_ALLOC) {
		return fmt.Errorf("error: %v doesn't exist, cannot do RmAlias", id)
	}
	remainder := Filter(id.Alias(), func(i int, e any) bool {
		return NotIn(e, aliases...)
	})
	mAlias.Store(id, remainder)
	return nil
}

func (id ID) ClrAlias() error {
	if NotIn(id.Type(), ID_HRCHY_ALLOC, ID_STDAL_ALLOC) {
		return fmt.Errorf("error: %v doesn't exist, cannot do RmAlias", id)
	}
	mAlias.Delete(id)
	return nil
}

//////////////////////////////////////////////////////////////////

func SearchIDByAlias(alias any, fromIDs ...ID) (ID, bool) {
	if len(fromIDs) == 0 {
		fromIDs = WholeIDs()
	}
	for _, id := range fromIDs {
		if In(alias, id.Alias()...) {
			return id, true
		}
	}
	return 0, false
}

func AddAlias(self any, aliases ...any) error {
	id, ok := SearchIDByAlias(self)
	if !ok {
		return fmt.Errorf("cannot find id by alias(%v)", self)
	}
	return id.AddAlias(aliases...)
}

// each alias for a single descendant
func CreateDescWithAlias(self any, aliases ...any) error {
	if self == "" {
		for _, alias := range aliases {
			nid, err := ID(0).GenDescID()
			if err != nil {
				return err
			}
			if err := nid.AddAlias(alias); err != nil {
				return err
			}
		}
		return nil
	}

	id, ok := SearchIDByAlias(self)
	if !ok {
		return fmt.Errorf("cannot find id by alias(%v)", self)
	}
	for _, alias := range aliases {
		if err := ID(id).AddAlias(alias); err != nil {
			return err
		}
	}
	return nil
}

func FetchAlias(self any) []any {
	id, ok := SearchIDByAlias(self)
	if !ok {
		return nil
	}
	return id.Alias()
}

func FetchDefaultAlias(self any) any {
	id, ok := SearchIDByAlias(self)
	if !ok {
		return nil
	}
	return id.DefaultAlias()
}

func RmAlias(self any, aliases ...any) error {
	id, ok := SearchIDByAlias(self)
	if !ok {
		return fmt.Errorf("cannot find id by alias(%v)", self)
	}
	return id.RmAlias(aliases...)
}

func ClrAlias(self any, aliases ...any) error {
	id, ok := SearchIDByAlias(self)
	if !ok {
		return fmt.Errorf("cannot find id by alias(%v)", self)
	}
	return id.ClrAlias()
}

func ChangeAlias(old, new any) error {
	if err := AddAlias(old, new); err != nil {
		return err
	}
	if err := RmAlias(new, old); err != nil {
		return err
	}
	return nil
}

func PrintAlias() {
	fmt.Println("-------------------------")
	mAlias.Range(func(key, value any) bool {
		fmt.Printf("%v: %v\n", key, value)
		return true
	})
	fmt.Println("-------------------------")
}
