package rulebase

import (
	"authz/prefixtree"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	// "os"
	"errors"
)

const (
	DENY int = iota
	ALLOW
)

type Config struct {
	Title string `yaml:,optional`
	Rules []struct {
		Subject string            `yaml:"Subject"`
		ACL     map[string]string `yaml:"ACL"`
	}
}

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

func Createrulebase(conf *Config) (*prefixtree.Tree, error) {
	tree := prefixtree.New()

	for _, r := range conf.Rules {
		for url, access := range r.ACL {
			// fmt.Printf("user: %s -> url: %s access: %s\n", r.Subject, url, access)
			a := 0
			if access == "allow" {
				a = ALLOW
			} else if access == "deny" {
				a = DENY
			} else {
				return nil, errors.New(fmt.Sprintf("Unknown access value %s\n", access))
			}
			tree.Add(url, r.Subject, a)
		}
	}

	return tree, nil
}

func TreeLookup(subject string, url string, rb *prefixtree.Tree) (int, error) {
	v, err := rb.Match(url, subject)
	return v, err
}

func Maprulebase(conf *Config) map[string]map[string]int {
	rb := make(map[string]map[string]int)

	for _, r := range conf.Rules {
		m := make(map[string]int)
		for url, access := range r.ACL {
			// fmt.Printf("user: %s -> url: %s access: %s\n", r.Subject, url, access)
			a := 0
			if access == "allow" {
				a = ALLOW
			} else if access == "deny" {
				a = DENY
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
