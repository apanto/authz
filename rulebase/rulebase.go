package rulebase

import (
	"authz/prefixtree"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	// "encoding/json"
	"io/ioutil"
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

type Rule struct {
	Url string              `yaml:"Url"`
	ACL map[string][]string `yaml:"ACL"`
}

type Config struct {
	Title string `yaml:"Title,omitempty"`
	Rules []Rule
}

//Read a configuration file in YAML format. Example:
// ---
// Title: "This is a test rulebase"
//
// rules:
//   - Url: www.corpA.com/*
//     ACL:
//       Jim: allow
//       John: allow
//   - Url: www.corpA.com/admin
//     ACL:
//       Jim: deny
//       John: allow
func Readconfig(filename string) (*Config, error) {
	var conf Config
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func Create(conf *Config) (*prefixtree.Tree, error) {
	tree := prefixtree.New()

	for _, r := range conf.Rules {
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
					return nil, errors.New(fmt.Sprintf("Unknown HTTP verb %s\n", v))
				}
			}
			// fmt.Printf("user: %s -> url: %s access: %s\n", subject, r.Url, access_flags)
			err := tree.AddKey(r.Url, subject, access_flags)
			if err != nil {
				return nil, err
			}
		}
	}

	// fmt.Println(*tree.Digraph())
	return tree, nil
}

//TODO: make default access policy configurable
func TreeLookup(subject string, url string, rb *prefixtree.Tree) (int, error) {
	v, err := rb.Match(url, subject)
	if err != nil {
		if err.Error() == "index does not exist" {
			//If the subject is not present in the ACL for this prefix the default access policy is DENY
			return DENY, nil
		} else {
			return DENY, err
		}
	}
	return v, nil
}

func Maprulebase(conf *Config) (map[string]map[string]int, error) {
	rb := make(map[string]map[string]int)

	for _, r := range conf.Rules {
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
