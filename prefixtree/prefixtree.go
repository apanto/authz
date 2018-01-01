//Package Tree
//TODO: Adapt the prefixtree to use a custome vocabulary so that if the custom vocabulary has less
//characters than the ASCII set we ca use smaler child arrays and save memory.
package prefixtree

import (
	"errors"
	"fmt"
	// "log"
)

// type iTree interface {
// 	insert(value string)
// 	lookup(value string) *Node
// 	print()
// }

type Tree struct {
	root *Node
}

//Each node stores an array of pointers to its childredn (child) and and array of characters
//representing the edge connecting this node to its respective child.
type Node struct {
	// index_child []byte
	child [128]*Node
	value map[string]int
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
func (n *Node) digraph() *string {
	var s string
	// if len(n.index_child) == 0 {
	// 	log.Printf("Node %p has 0 children\n", n)
	// }
	// log.Printf("node %p has %d childred\n", &n, len(n.index_child))
	for i, child := range n.child {
		if child != nil {
			s += fmt.Sprintf("  \"%p\" -> \"%p\" [ label = \"%s\" ]; \n", n, child, string(i))
			s += *child.digraph()
		}
	}
	return &s
}

func (n *Node) String() string {
	return fmt.Sprintf("Node(%d): \"%p\" -> (%d)%v", n.value, n, len(n.child), n.child)
}

func New() *Tree {
	tree := new(Tree)
	tree.root = new(Node)
	return tree
}

func (t Tree) Add(key string, index string, value int) {
	var next *Node
	if key == "" {
		return
	}
	n := t.root

	for i := 0; i < len(key); i++ {

		next = n.child[key[i]]
		if next == nil {
			next, _ = n.add(key[i])
			// n.index_child = append(n.index_child, key[i])
			// n.child = append(n.child, new(Node))
			// next = n.child[len(n.child)-1]
		}
		n = next
	}

	// //if the last character of the key we are adding exists then the key
	// //exiists and we should not add it
	// next = n.next(key[len(key)-1])
	// if next == nil {
	// 	n = n.add(key[len(key)-1])
	// } else {
	// 	fmt.Printf("ERROR: key %s already exists with value %d\n", key, next.value)
	// }

	// for _, c := range key[:len(key)-1] {
	// 	log.Printf("inserting %s into node %p (%T)\n", string(c), n, c)
	// 	index := -1
	// 	for i, k := range n.index_child {
	// 		if k == byte(c) {
	// 			index = i
	// 			break
	// 		}
	// 	}
	// 	if index == -1 {
	// 		n.index_child = append(n.index_child, byte(c))
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
	// 	n = n.child[index]
	// }

	if n.value == nil {
		n.value = make(map[string]int)
	}
	n.value[index] = value
	// fmt.Printf("key: %s: %v\n", key, n.value)
}

func (t Tree) Match(key string, index string) (int, error) {
	var wildcard *Node
	n := t.root

	for k := 0; n != nil && k < len(key); k++ {
		if n.child['*'] != nil {
			wildcard = n.child['*']
		}
		n = n.child[key[k]]
		// next := n.child[key[k]]
		// if next != nil {
		// 	n = next
		// } else {
		// 	n = n.child['*']
		// 	break
		// }

		// // if next = n.next(key[i]); next == nil {
		// index := -1
		// for i, v := range n.index_child {
		// 	if v == key[k] {
		// 		index = i
		// 		break
		// 	}
		// }

		// if index == -1 {
		// 	n = n.next('*')
		// } else {
		// 	n = n.child[index]
		// }

		// // if next == nil {
		// // 	n = n.next('*')
		// // 	break
		// // }
		// // n = next
	}

	if n != nil {
		v, exists := n.value[index]
		if exists {
			return v, nil
		} else {
			return 0, errors.New("index does not exist")
		}
	} else { // n == nil
		if wildcard != nil {
			v, exists := wildcard.value[index]
			if exists {
				return v, nil // return value stored in wildcard node in the partially matching prefix
			} else {
				return 0, errors.New("index does not exist")
			}
		} else { // wildcard == nil
			return 0, errors.New("Key does not match")
		}
	}
}

func (t Tree) Get(key string, index string) (int, error) {
	n := t.root

	for i := 0; n != nil && i < len(key); i++ {
		n = n.child[key[i]]
	}

	if n != nil {
		v, exists := n.value[index]
		if exists {
			return v, nil
		} else {
			return 0, errors.New("index does not exist")
		}
	} else {
		return 0, errors.New("Key does not exist")
	}

}

func (t Tree) Digraph() *string {
	var s string

	s = fmt.Sprintf("digraph G {\n")
	s += *t.root.digraph()
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
