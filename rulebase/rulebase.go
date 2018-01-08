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
	group                map[string][]string
	Groups               []string
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
	rb.group = make(map[string][]string)
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

func (rb Rulebase) AddGroups(groups map[string][]string) {
	//Setup the groups map
	if groups != nil {
		for group, subjects := range groups {
			rb.AddGroup(group, subjects)
		}
	}
	// fmt.Printf("Group map:\n")
	// for s, gs := range rb.group {
	// 	fmt.Printf("subject: %s is memeber of %v\n", s, gs)
	// }
}

func (rb Rulebase) AddGroup(group string, subjects []string) {
	if group != "" {
		rb.Groups = append(rb.Groups, group)
		for _, subject := range subjects {
			rb.group[subject] = append(rb.group[subject], group)
		}
	}
}

//Changes the members of a group. i.e. remover subject X from group Y
func (rb Rulebase) ModGroup(group string, subjects []string) {
}

//Deletes a group and removes all subjects from it.
func (rb Rulebase) DelGroup(group string, subjects []string) {
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
//This only looksup the subject ignoring any groups the subjet is member of
func (rb Rulebase) LookupSubject(subject string, url string) (int, error) {
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

//Looks up a subject and url in the rulebase. This also looks up the groups the subject is member of and
//returns the "combined" access flags.
func (rb Rulebase) Lookup(subject string, url string) (int, error) {
	var access_flags, group_flags, subject_flags int

	// fmt.Printf("Lookup: %s@%s\n", subject, url)
	key_map, err := rb.tree.MatchPrefix(url)
	if err != nil {
		if err.Error() == "prefix does not match" {
			//If the subject is not present in the ACL for this prefix return the default access flags of this rb
			return rb.default_access_flags, nil
		} else {
			return 0, err
		}
	}
	subject_flags = key_map[subject] //if subject doesn't exist subject_flags are 0.
	// fmt.Printf("  subject_flags(%s) %08b\n", subject, subject_flags)

	//Any groups that have an ACL for a prefix matching this URL will exist in the key_map of this prefix.
	//So all we have to do to get all the relevant access flags for this user is to lookup each group the
	//user is a member of in the key map of this prefix. All access flags are then ORed together to calculate
	//the final access flags
	if groups, exists := rb.group[subject]; exists {
		for _, g := range groups {
			if flags, exists := key_map[g]; exists {
				//get this group's access flags and OR it with the previously aggregated access flags
				fmt.Printf("  group_flags(%s): %08b\n", g, flags)
				group_flags |= flags
			}
		}
	}

	access_flags = subject_flags | group_flags | rb.default_access_flags

	return access_flags, nil
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
