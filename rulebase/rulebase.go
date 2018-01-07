package rulebase

import (
	"authz/prefixtree"
	"errors"
	"fmt"
)

const (
	DENY int = iota
	ALLOW
)

const (
	GET    = 1
	PUT    = 2
	POST   = 4
	DELETE = 8
	UPDATE = 16
)

type Rulebase struct {
	tree                 *prefixtree.Tree
	default_access_flags int
}

type Rule struct {
	Url string
	ACL map[string][]string
}

//Creates a new empty rulebase
func New() *Rulebase {
	var rb Rulebase
	rb.tree = prefixtree.New()
	rb.SetDefaultAccess([]string{})
	return &rb
}

//Sets the access flags for the default access policy. Expects and array of HTTP verbs as strings
func (rb Rulebase) SetDefaultAccess(access []string) error {
	access_flags := 0
	for _, v := range access {
		switch v {
		case "GET":
			access_flags += GET
		case "PUT":
			access_flags += PUT
		case "POST":
			access_flags += POST
		case "DELETE":
			access_flags += DELETE
		case "UPDATE":
			access_flags += UPDATE
		default:
			return errors.New(fmt.Sprintf("Unknown HTTP verb %s\n", v))
		}
	}

	rb.default_access_flags = access_flags
	return nil
}

//Creates a new rulebase from an array of rules
func Create(rules *[]Rule) (*Rulebase, error) {
	rb := New()

	for _, r := range *rules {
		err := rb.Add(&r)
		if err != nil {
			return nil, err
		}
	}

	return rb, nil
}

//Adds a rule to a rulebase
func (rb Rulebase) Add(r *Rule) error {
	for subject, access := range r.ACL {
		// fmt.Printf("user: %s -> url: %s access(%T): %s\n", subject, r.Url, access, access)
		access_flags := 0
		for _, v := range access {
			switch v {
			case "GET":
				access_flags += GET
			case "PUT":
				access_flags += PUT
			case "POST":
				access_flags += POST
			case "DELETE":
				access_flags += DELETE
			case "UPDATE":
				access_flags += UPDATE
			default:
				return errors.New(fmt.Sprintf("Unknown HTTP verb %s\n", v))
			}
		}
		// fmt.Printf("user: %s -> url: %s access: %s\n", subject, r.Url, access_flags)
		err := rb.tree.AddKey(r.Url, subject, access_flags)
		if err != nil {
			return err
		}
	}

	// fmt.Println(*rb.tree.Digraph())
	return nil
}

//Looks up a subject and url in the rulebase and returns the access flags as int
func (rb Rulebase) Lookup(subject string, url string) (int, error) {
	v, err := rb.tree.Match(url, subject)
	if err != nil {
		if err.Error() == "index does not exist" {
			//If the subject is not present in the ACL for this prefix return the default access flags of this rb
			return rb.default_access_flags, nil
		} else {
			return 0, err
		}
	}
	return v, nil
}

func Maprulebase(rules *[]Rule) (map[string]map[string]int, error) {
	rb := make(map[string]map[string]int)

	for _, r := range *rules {
		m := make(map[string]int)
		for subject, access := range r.ACL {
			// fmt.Printf("user: %s -> url: %s access: %s\n", subject, r.Url, access_flags)
			access_flags := 0
			for _, v := range access {
				switch v {
				case "GET":
					access_flags += GET
				case "PUT":
					access_flags += PUT
				case "POST":
					access_flags += POST
				case "DELETE":
					access_flags += DELETE
				case "UPDATE":
					access_flags += UPDATE
				default:
					return nil, errors.New(fmt.Sprintf("Unknown HTTP verb %s\n", v))
				}
			}
			m[subject] = access_flags
		}
		rb[r.Url] = m
	}
	// // fmt.Printf("%v\n", rb)
	return rb, nil
}

func MapLookup(subject string, url string, rb map[string]map[string]int) int {
	return rb[url][subject]
}
