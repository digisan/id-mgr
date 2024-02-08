package id

import (
	"fmt"
	"strings"
	"sync"

	. "github.com/digisan/go-generics"
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

func CheckAlias(aliases []any, fromIDs ...ID) (any, error) {
	if len(fromIDs) == 0 {
		fromIDs = WholeIDs()
	}
	for _, alias := range aliases {
		if !validateAlias(alias) {
			return alias, fmt.Errorf("'%v' contains invalid characters(%+v)", alias, exclChars)
		}
		if used, byId := aliasOccupied(alias, fromIDs...); used {
			return alias, fmt.Errorf("'%v' is already used by [%x]", alias, byId)
		}
	}
	return nil, nil
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
AGAIN:
	if aliasErr, err := CheckAlias(aliases); err != nil {
		if strings.Contains(err.Error(), "used") {
			aliases = FilterMap(aliases, nil, func(i int, e any) any {
				if aliasErr == e {
					return DuplMark(e)
				}
				return e
			})
			goto AGAIN
		} else {
			return err
		}
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
	if len(aliases) > 0 {
		remainder := Filter(id.Alias(), func(i int, e any) bool {
			return NotIn(e, aliases...)
		})
		mAlias.Store(id, remainder)
	} else {
		mAlias.Delete(id)
	}
	return nil
}

func (id ID) ClrAlias() error {
	if NotIn(id.Type(), ID_HRCHY_ALLOC, ID_STDAL_ALLOC) {
		return fmt.Errorf("error: %v doesn't exist, cannot do RmAlias", id)
	}
	mAlias.Delete(id)
	return nil
}

func (id ID) SetAlias(aliases ...any) error {
	if err := id.RmAlias(); err != nil {
		return err
	}
	if err := id.AddAlias(aliases...); err != nil {
		return err
	}
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

func SearchIDByAnyAlias(aliases ...any) (ID, bool) {
	fromIDs := WholeIDs()
	for _, alias := range aliases {
		for _, id := range fromIDs {
			if In(alias, id.Alias()...) {
				return id, true
			}
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

// create ONE descendant of super with multiple input aliases
func CreateOneDescWithAlias(super any, aliases ...any) error {
	if super == "" {
		nid, err := ID(0).GenDescID()
		if err != nil {
			return err
		}
		return nid.AddAlias(aliases...)
	}

	sid, ok := SearchIDByAlias(super)
	if !ok {
		return fmt.Errorf("cannot find id by alias(%v)", super)
	}
	nid, err := sid.GenDescID()
	if err != nil {
		return err
	}
	return nid.AddAlias(aliases...)
}

// create ONE standalone with multiple input aliases
func CreateOneStdalWithAlias(aliases ...any) error {
	nid, err := GenStdalID()
	if err != nil {
		return err
	}
	return nid.AddAlias(aliases...)
}

// Create MULTIPLE descendants, and each descendant with one alias
func BuildHierarchy(super any, aliases ...any) error {
	for _, alias := range aliases {
		if err := CreateOneDescWithAlias(super, alias); err != nil {
			return err
		}
	}
	return nil
}

// Create MULTIPLE standalone id, and each id with one alias
func BuildStandalone(aliases ...any) error {
	for _, alias := range aliases {
		if err := CreateOneStdalWithAlias(alias); err != nil {
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

// [n] is descendant generation, n(1) is children level, n(2) is grandchildren level ...
// return includes self
func FetchDescendantDefaultAlias(n int, self any) (descendants []any) {
	id, ok := SearchIDByAlias(self)
	if !ok {
		return
	}
	for _, desc := range id.Descendants(n, true) {
		if a := desc.Alias(); len(a) > 0 {
			descendants = append(descendants, a[0])
		}
	}
	return
}

// return doesn't include self
func FetchChildrenAlias(self any) []any {
	rt := FetchDescendantDefaultAlias(1, self)
	if len(rt) > 0 {
		return rt[1:]
	}
	return rt
}

// [n] is ancestor generation, n(1) is parent level, n(2) is grandparent level ...
// return includes self
func FetchAncestorDefaultAlias(n int, self any) (ancestors []any) {
	id, ok := SearchIDByAlias(self)
	if !ok {
		return
	}
	for _, anc := range id.Ancestors(true) {
		if a := anc.Alias(); len(a) > 0 {
			ancestors = append(ancestors, a[0])
		}
	}
	rt := Reverse(ancestors)
	if n < len(rt) {
		return rt[:n+1]
	}
	return rt
}

// return doesn't include self
func FetchParentDefaultAlias(self any) any {
	rt := FetchAncestorDefaultAlias(1, self)
	if len(rt) > 0 {
		return rt[1:]
	}
	return rt
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

func ClrAllAlias() error {
	for _, id := range WholeIDs() {
		if err := id.RmAlias(); err != nil {
			return err
		}
	}
	return nil
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

// DeleteIDByAlias incurs updated WholeIDs
func DeleteIDByAlias(alias any, inclDesc bool) error {
	if id, ok := SearchIDByAlias(alias); ok {
		if _, err := DeleteID(id, inclDesc); err != nil {
			return err
		}
	}
	return nil
}

// only delete leaves id
func DeleteIDsByAlias(aliases ...any) error {
	for _, alias := range aliases {
		if err := DeleteIDByAlias(alias, false); err != nil {
			return err
		}
	}
	return nil
}

func cleanupAlias() error {
	list := []ID{}
	mAlias.Range(func(key, value any) bool {
		// fmt.Println("mAlias.Range:", key, value)
		aliases := value.([]any)
		for _, alias := range aliases {
			s := fmt.Sprint(alias)
			io, ic := strings.LastIndex(s, "("), strings.LastIndex(s, ")")
			if io == -1 || ic == -1 {
				continue
			}
			idxstr := s[io+1 : ic]
			if _, ok := AnyTryToType[int](idxstr); ok && strings.HasSuffix(s, ")") {
				list = append(list, key.(ID))
				break
			}
		}
		return true
	})
	// fmt.Println(list)

	for _, id := range list {
		aliases := FilterMap(id.Alias(), nil, func(i int, e any) any {
			s := fmt.Sprint(e)
			io, ic := strings.LastIndex(s, "("), strings.LastIndex(s, ")")
			if io == -1 || ic == -1 {
				return e
			}
			name, idxstr := s[:io], s[io+1:ic]
			if _, ok := AnyTryToType[int](idxstr); ok && strings.HasSuffix(s, ")") {
				if _, ok := mAlias.Load(name); !ok {
					return name
				}
			}
			return e
		})
		if err := id.SetAlias(aliases...); err != nil {
			return err
		}
	}
	return nil
}
