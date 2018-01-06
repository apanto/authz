//Package Tree
//TODO: Adapt the prefixtree to use a custome vocabulary so that if the custom vocabulary has less
//characters than the ASCII set we ca use smaler child arrays and save memory.
package prefixtree

import (
	"errors"
	"fmt"
	"regexp"
)

// type iTree interface {
// 	insert(value string)
// 	lookup(value string) *Node
// 	print()
// }

//Regexp used in Tree.Add to find wildcard characters ('*') in input keys.
var re_star *regexp.Regexp = regexp.MustCompile(`\*`)

type Tree struct {
	root *Node
}

//Each node stores an array of pointers to its childredn (child) and and array of characters
//representing the edge connecting this node to its respective child.
type Node struct {
	wildcard bool
	child    [128]*Node
	value    map[string]int
}

func (n *Node) add(k byte) (*Node, error) {
	// n.index_child = append(n.index_child, k)

	child := new(Node)
	// n.child = append(n.child, child)
	if k > 127 {
		return nil, errors.New(fmt.Sprintf("Key %s is not an ASCII character", string(k)))
	}

	n.child[k] = child
	return child, nil
}

//Return a *string of a .dot notation of this node and all its children
//Node adresses are displayed in each node of the graph and keys are
//displayed on the edges. For bievety the values of the nodes are not displayed
func (n *Node) digraph() (*string, *string) {
	var s, wildcards string

	// if len(n.index_child) == 0 {
	// 	log.Printf("Node %p has 0 children\n", n)
	// }
	// log.Printf("node %p has %d childred\n", &n, len(n.index_child))
	if n.wildcard {
		wildcards += fmt.Sprintf(" \"%p\"", n)
	}
	for i, child := range n.child {
		if child != nil {
			s += fmt.Sprintf("  \"%p\" -> \"%p\" [ label = \"%s\" ]; \n", n, child, string(i))
			s_tmp, wildcards_tmp := child.digraph()
			s += *s_tmp
			wildcards += *wildcards_tmp
		}
	}
	return &s, &wildcards
}

func (n *Node) String() string {
	return fmt.Sprintf("Node(%d): \"%p\" -> (%d)%v", n.value, n, len(n.child), n.child)
}

func New() *Tree {
	tree := new(Tree)
	tree.root = new(Node)
	return tree
}

//Add a prefix and initialize the value map. addprefix is idempotent i.e. if the prefix
//and/or the value map exist nothing will happen the tree t will remain unchanged
//
//TODO: what happens when you encounter an error mid flight during insertion of a key,
//how do you remove the part of the key that was already inserted?
func (t Tree) addprefix(prefix string) (*Node, error) {
	// fmt.Printf("%s, %v\n", prefix, re_star.FindString(prefix[:len(prefix)-1]) != "")
	if len(prefix) == 0 {
		return nil, errors.New("prefix cannot be empty")
	} else if re_star.FindString(prefix[:len(prefix)-1]) != "" {
		return nil, errors.New("prefix cannot contain '*' except at the end")
	}
	n := t.root

	for _, p := range prefix {
		if p == '*' {
			n.wildcard = true
		} else {
			if n.child[p] == nil {
				n.child[p], _ = n.add(byte(p))
			}
			n = n.child[p]
		}

	}

	if n.value == nil {
		n.value = make(map[string]int)
	}

	return n, nil
}

func (t Tree) SetKeys(prefix string, keys map[string]int) error {
	n, err := t.addprefix(prefix)
	if err != nil {
		return err
	}

	n.value = keys

	return nil
}

func (t Tree) AddKeys(prefix string, keys map[string]int) error {
	n, err := t.addprefix(prefix)
	if err != nil {
		return err
	}

	for k, v := range keys {
		n.value[k] = v
	}

	return nil
}

func (t Tree) AddKey(prefix string, key string, value int) error {

	n, err := t.addprefix(prefix)
	if err != nil {
		return err
	}

	n.value[key] = value
	// fmt.Printf("prefix: %s: %v\n", prefix, n.value)
	return nil
}

func (t Tree) Match(prefix string, key string) (int, error) {
	var wildcard *Node
	n := t.root

	// fmt.Printf("prefix: %s, key: %s\n", prefix, key)
	for p := 0; n != nil && p < len(prefix); p++ {
		if n.wildcard {
			wildcard = n // if a wildcard node is encountered along the path note it for checcking on later
		}
		n = n.child[prefix[p]]
		// fmt.Printf("  p: %s\n", string(prefix[p]))
	}

	if n != nil {
		v, exists := n.value[key]
		if exists {
			return v, nil
		} else {
			// fmt.Printf("Subjects: %v\n", n.value)
			return 0, errors.New("key does not exist")
		}
	} else { // n == nil
		if wildcard != nil {
			v, exists := wildcard.value[key]
			if exists {
				return v, nil // return value stored in wildcard node in the partially matching prefix
			} else {
				return 0, errors.New("key does not exist")
			}
		} else { // wildcard == nil
			return 0, errors.New("prefix does not match")
		}
	}

	return 0, nil
}

func (t Tree) Get(prefix string, key string) (int, error) {
	n := t.root

	for i := 0; n != nil && i < len(prefix); i++ {
		n = n.child[prefix[i]]
	}

	if n != nil {
		v, exists := n.value[key]
		if exists {
			return v, nil
		} else {
			return 0, errors.New("key does not exist")
		}
	} else {
		return 0, errors.New("prefix does not exist")
	}

}

func (t Tree) Digraph() *string {
	var s string

	nodes, wildcards := t.root.digraph()
	if len(*wildcards) > 0 {
		s = fmt.Sprintf("digraph G {\n  size=\"8,5\"\n  node [shape = doublecircle];%s;\n  node [shape = circle];\n", *wildcards)
	} else {
		s = fmt.Sprintf("digraph G {\n  size=\"8,5\"\n  node [shape = circle];\n")
	}
	s += *nodes
	s += fmt.Sprintf("}\n")

	return &s
}

// func tree_lookup_r(n *Node, value string) *Node {
// 	index := -1
// 	for i, v := range n.index_child {
// 		if v == value[0] {
// 			index = i
// 			break
// 		}
// 	}
// 	if index == -1 {
// 		return nil
// 	} else {
// 		log.Printf("value: %s index: %d\n", value, index)
// 		if len(value) == 1 {
// 			return n.child[index]
// 		} else {
// 			// log.Printf("%d\n", index)
// 			return tree_lookup_r(n.child[index], value[1:])
// 		}
// 	}
// }

// func tree_insert_r(n *Node, value string) {
// 	if value == "" {
// 		return
// 	}
// 	log.Printf("inserting %s into node %p\n", string(value[0]), n)
// 	c := byte(value[0])
// 	index := -1
// 	for i, v := range n.index_child {
// 		if v == c {
// 			index = i
// 			break
// 		}
// 	}
// 	if index == -1 {
// 		n.index_child = append(n.index_child, c)
// 		// log.Printf("  New node: %p\n", &new_node)
// 		// log.Printf("  %s doesn't exist, inserting and adding new node %p\n", string(value[0]), &new_node)
// 		n.child = append(n.child, new(Node))
// 		// log.Printf("  child: %p\n", n.child[len(n.child)-1])
// 		// tree_insert(n.child[len(n.child)-1], value[1:])
// 		index = len(n.child) - 1
// 	} else {
// 		// log.Printf("  %s exists at %d, moving on...\n", string(value[0]), index)
// 		// tree_insert(n.child[index], value[1:])
// 	}
// 	tree_insert_r(n.child[index], value[1:])
// }
