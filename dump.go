package idmgr

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	. "github.com/digisan/go-generics"
	. "github.com/digisan/id-mgr/id"
)

func GenIDTree(print bool) string {

	lines := []string{}
	for i, id := range HierarchyIDs() {
		lvl := id.Level()
		indent := strings.Repeat("\t", lvl)
		aliasesStr := ""
		if aliases, ok := AnysTryToTypes[string](id.Alias()); ok {
			aliasesStr = strings.Join(aliases, "^")
		}
		lines = append(lines, fmt.Sprintf("%s%x|%v", indent, id, aliasesStr)) // generated string use hexadecimal id

		if print {
			fmt.Printf("%03d: %s%x|%v\n", i+1, indent, id, aliasesStr) // print with line number, use hexadecimal 0xid
		}
	}

	offset := len(lines)
	for i, id := range StandaloneIDs() {
		aliasesStr := ""
		if aliases, ok := AnysTryToTypes[string](id.Alias()); ok {
			aliasesStr = strings.Join(aliases, "^")
		}
		lines = append(lines, fmt.Sprintf("%x|%v", id, aliasesStr))
		if print {
			fmt.Printf("%03d: %x|%v\n", i+1+offset, id, aliasesStr) // print with line number, use hexadecimal 0xid
		}
	}

	rt := strings.Join(lines, "\n")
	return rt
}

func DumpIDTree(fPath string) error {
	out := PrintSegs(false) + "\n" + GenIDTree(false)
	return os.WriteFile(fPath, []byte(out), os.ModePerm)
}

//////////////////////////////////////////////////////////////////////
// dump as nested json

type Node struct {
	ID   uint64
	Name string
	Desc []*Node
}

func (n *Node) Append(id uint64, name string) {
	n.Desc = append(n.Desc, &Node{ID: id, Name: name, Desc: []*Node{}})
}

func (n *Node) Search(id uint64) *Node {
	if n.ID == id {
		return n
	}
	for _, d := range n.Desc {
		if r := d.Search(id); r != nil {
			return r
		}
	}
	return nil
}

func CvtTree2JSON(fPathIn, fPathOut string) (nodes []*Node, err error) {
	if err = IngestTree(fPathIn); err != nil {
		return nil, err
	}

	ids, _, aliases := ID(0).DescendantsInfo(false)
	for i, id := range ids {
		alias := aliases[i][0]
		if pid, ok := ID(id).Parent(); ok {
			if pid == 0 {
				nodes = append(nodes, &Node{
					ID:   uint64(id),
					Name: alias.(string),
					Desc: []*Node{},
				})
			} else if pid > 0 {
				for _, t := range nodes {
					if n := t.Search(uint64(pid)); n != nil {
						n.Append(uint64(id), alias.(string))
						break
					}
				}
			}
		}
	}

	bytes, err := json.Marshal(nodes)
	if err != nil {
		return nil, err
	}
	return nodes, os.WriteFile(fPathOut, bytes, os.ModePerm)
}
