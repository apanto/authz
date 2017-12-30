package authz

import (
	"authz/prefixtree"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	// "os"
)

type Config struct {
	Title string `yaml:,optional`
	Rules []struct {
		Subject string            `yaml:"subject"`
		ACL     map[string]string `yaml:"ACL"`
	}
}

func Readconfig(filename string) *Config {
	var conf Config
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		panic(err)
	}

	return &conf
}

func Createrulebase(conf *Config) *prefixtree.Tree {
	tree := prefixtree.New()

	for _, r := range conf.Rules {
		for url, access := range r.ACL {
			// fmt.Printf("user: %s -> url: %s access: %s\n", r.Subject, url, access)
			a := 0
			if access == "allow" {
				a = 1
			} else if access == "deny" {
				a = 0
			} else {
				a = 0
				fmt.Errorf("Unknown access value %s\n", access)
			}
			tree.Add(url, r.Subject, a)
		}
	}

	return tree
}

func TreeLookup(subject string, url string, rb *prefixtree.Tree) int {
	v, _ := rb.Match(url, subject)
	return v
}

func Maprulebase(conf *Config) map[string]map[string]int {
	rb := make(map[string]map[string]int)

	for _, r := range conf.Rules {
		m := make(map[string]int)
		for url, access := range r.ACL {
			// fmt.Printf("user: %s -> url: %s access: %s\n", r.Subject, url, access)
			a := 0
			if access == "allow" {
				a = 1
			} else if access == "deny" {
				a = 0
			} else {
				a = 0
				fmt.Errorf("Unknown access value %s\n", access)
			}
			m[url] = a
		}
		rb[r.Subject] = m
	}

	return rb
}

func MapLookup(subject string, url string, rb map[string]map[string]int) int {
	return rb[subject][url]
}
