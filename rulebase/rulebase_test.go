package rulebase

import (
	"authz/prefixtree"
	"encoding/base64"
	// "fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

const test_conf = `---
title: "This is a test rulebase"

rules:
  - Url: www.corpA.com/*
    ACL:
      Jim: [GET,POST]
      John: [GET]
  - Url: www.corpA.com/admin
    ACL:
      Jim: [GET,POST]
      John: [GET]
  - Url: www.public.org/*
    ACL:
      anonymous: [GET]
  - Url: www.public.org/secret*
    ACL:
      anonymous: []`

const invalid_test_conf1 = `---
title: "This is a test rulebase"

rules:
  - Url: www.corp*.com/*
    ACL:
      Jim: [GET,POST]
      John: [GET]
  - Url: www.corpA.com/admin
    ACL:
      Jim: [GET,POST]
      John: [GET]
  - Url: www.public.org/*
    ACL:
      anonymous: [GET]
  - Url: www.public.org/secret*
    ACL:
      anonymous: []`

const invalid_test_conf2 = `---
title: "This is a test rulebase"

rules:
  - Url: www.corpA.com/*
    ACL:
      Jim: [GET,POST]
      John: allow
  - Url: www.corpA.com/admin
    ACL:
      Jim: [GET,POST]
      John: [GET]
  - Url: www.public.org/*
    ACL:
      anonymous: [GET]
  - Url: www.public.org/secret*
    ACL:
      anonymous: []`

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

var test_map_rb map[string]map[string]int
var test_tree_rb *prefixtree.Tree

func readtestconfig() {
	config_filename := "./tmp_conf.yml"
	err := ioutil.WriteFile(config_filename, []byte(test_conf), 0644)
	if err != nil {
		log.Fatalf("Can't read config file %s", config_filename)
	}

	conf, err := Readconfig(config_filename)
	if err != nil {
		log.Fatalf("Cpuldn't read configuration file %s (%s)", config_filename, err)
	}
	test_tree_rb, err = Create(conf)
	if err != nil {
		log.Fatalf("Couldn't create rulebase from configuration (%s)", err)
	}

	test_map_rb, _ = Maprulebase(conf)
}

func TestMain(m *testing.M) {
	//setup()
	//fmt.Printf("m:%#v", m)
	readtestconfig()

	code := m.Run()
	// shutdown()
	os.Exit(code)
}

//TODO: Test with invalid_test_conf2 file seg faults and halts the tests. How can this be avoided?
func TestCreaterulebaseWithInvalidConfig(t *testing.T) {
	var err error
	config_filename := "./tmp_conf.yml"

	err = ioutil.WriteFile(config_filename, []byte(invalid_test_conf1), 0644)
	if err != nil {
		t.Fatalf("Can't write config file %s", config_filename)
	}

	conf, err := Readconfig(config_filename)
	if err != nil {
		t.Error(err)
	}
	test_tree_rb, err = Create(conf)
	if err == nil {
		t.Error("Create() didn't fail with invalid config file 1 (URL containing `*` not at the end)")
	}

	// err = ioutil.WriteFile(config_filename, []byte(invalid_test_conf2), 0644)
	// if err != nil {
	// 	t.Fatalf("Can't write config file %s", config_filename)
	// }

	// conf, err = Readconfig(config_filename)
	// if err != nil {
	// 	t.Error(err)
	// }
	// test_tree_rb, err = Create(conf)
	// if err == nil {
	// 	t.Error("Create() didn't fail with invalid config file 2 (invalid access flags)")
	// }

}

func TestInsertLookupMap(t *testing.T) {
	subject := "John"
	url := "www.corpA.com/*"
	access := MapLookup(subject, url, test_map_rb)
	if access != (GET) {
		t.Errorf("Wrong authorization value %d for %s:%s\n ", access, subject, url)
	}

	subject = "anonymous"
	url = "www.corpA.com/*"
	access = MapLookup(subject, url, test_map_rb)
	if access != 0 {
		t.Errorf("Wrong authorization value %d for %s:%s\n ", access, subject, url)
	}

	subject = "anonymous"
	url = "www.public.org/secret*"
	access = MapLookup(subject, url, test_map_rb)
	if access != 0 {
		t.Errorf("Wrong authorization value %d for %s:%s\n ", access, subject, url)
	}
}

func TestInsertLookupTree(t *testing.T) {
	var err error
	var config_file string = "test_conf.yml"
	conf, _ := Readconfig(config_file)

	test_tree_rb, err = Create(conf)
	if err != nil {
		t.Errorf("Can't create rulebase using config file %s\n", config_file)
	}

	subject := "John"
	url := "www.corpA.com/home"

	val, err := Lookup(subject, url, test_tree_rb)
	if err != nil {
		t.Fatalf("Error %s\n", err)
	}
	if val != (GET + POST) {
		t.Errorf("Wrong authorization value for %s:%s\n", subject, url)
	}

	subject = "anonymous"
	url = "www.corpA.com/"

	val, err = Lookup(subject, url, test_tree_rb)
	if err.Error() != "key does not exist" {
		t.Fatalf("This should fail to find the subject %s, but failed with error: %s\n", subject, err)
	}
	if val != 0 {
		t.Errorf("Wrong authorization value for %s:%s\n", subject, url)
	}

	subject = "anonymous"
	url = "www.public.org/main"

	val, err = Lookup(subject, url, test_tree_rb)
	if err != nil {
		t.Fatalf("Error %s\n", err)
	}
	if val != (GET) {
		t.Errorf("Wrong authorization value for %s:%s\n", subject, url)
	}

	subject = "anonymous"
	url = "www.public.org/secret"

	val, err = Lookup(subject, url, test_tree_rb)
	if err != nil {
		t.Fatalf("Error %s\n", err)
	}
	if val != 0 {
		t.Errorf("Wrong authorization value for %s:%s\n", subject, url)
	}

}

// -----------------------------------
// Benchmarks
// -----------------------------------

var conf *Config
var rb_map map[string]map[string]int
var rb_tree *prefixtree.Tree
var initialized, rb_tree_init, rb_map_init int
var num_subjects int = 200
var num_urls int = 100

func createconfig(num_subjects int, num_urls int) *Config {
	var conf Config
	var acl map[string][]string
	var r int64

	acl = make(map[string][]string)

	conf.Title = "This is a test rulebase"

	for u := 0; u < num_urls; u++ {
		r = src.Int63()
		url, _ := GenerateRandomString(int(r%115) + 1)
		// fmt.Printf("url: %s\n", str)
		for s := 0; s < num_subjects; s++ {
			subject_id := RandString(12)
			if r%1 == 1 {
				acl[subject_id] = []string{"GET", "POST", "PUT"}
			} else {
				acl[subject_id] = []string{"GET"}
			}
		}

		conf.Rules = append(conf.Rules, struct {
			Url string              `yaml:"Url"`
			ACL map[string][]string `yaml:"ACL"`
		}{url, acl})
	}

	return &conf
}

func BenchmarkLookupTree(b *testing.B) {
	var subject, url string
	// var access string

	if initialized == 0 {
		b.Log("Initialising...")
		conf = createconfig(num_subjects, num_urls)
		initialized = 1
	}

	if rb_tree_init == 0 {
		rb_tree, _ = Create(conf)
		rb_tree_init = 1
	}

	//Pick a random rule ...
	// rule := conf.Rules[int(src.Int63())%num_urls]
	// subject = rule.ACL
	//Pick a random ACL rule from that rule
	// i := int(src.Int63()) % nr
	// for k, _ := range rule.ACL {
	// 	i--
	// 	if i == 0 {
	// 		url = k
	// 		// access = v
	// 	}
	// }

	var urls []string
	var l, avg int

	for _, r := range conf.Rules {
		// for u := range r.Url {
		avg = avg + len(r.Url)
		if len(r.Url) > l {
			l = len(r.Url)
			url = r.Url
		}
		// if u[len(u)-1] == '*' {
		// 	x := src.Int63()
		// 	str, _ = GenerateRandomString(int(x % 20))
		// }
		urls = append(urls, url)
		// u = u[:len(u)-1] + str
		// }
	}

	// fmt.Printf("longest url has %d chars, avg url length is %d\n", l, avg/len(urls))

	Lookup("fwgrgerwghwe", url, rb_tree)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Lookup(subject, urls[i%len(urls)], rb_tree)
	}

}

func BenchmarkLookupMap(b *testing.B) {
	var subject string
	// var access string

	if initialized == 0 {
		b.Log("Initialising...")
		conf = createconfig(num_subjects, num_urls)
		initialized = 1
	}
	if rb_map_init == 0 {
		rb_map, _ = Maprulebase(conf)
		rb_map_init = 1
	}

	//Pick a random rule and subject...
	rule := conf.Rules[int(src.Int63())%num_urls]
	for s := range rule.ACL {
		subject = s
	}

	var urls []string
	for _, r := range conf.Rules {
		urls = append(urls, r.Url)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		MapLookup(subject, urls[i%len(urls)], rb_map)
	}

}
