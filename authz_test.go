package authz

import (
	// "fmt"
	"authz/prefixtree"
	"encoding/base64"
	"math/rand"
	"os"
	"testing"
	"time"
)

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

func TestMain(m *testing.M) {
	//setup()
	//fmt.Printf("m:%#v", m)
	code := m.Run()
	// shutdown()
	os.Exit(code)
}

func TestInsertLookupMap(t *testing.T) {
	conf := Readconfig("test_conf.yml")

	rb := Maprulebase(conf)
	subject := "John"
	url := "www.corpA.com/*"

	if MapLookup(subject, url, rb) != 1 {
		t.Errorf("Whrong authorization value for %s:%s\n", subject, url)
	}
}

func TestInsertLookupTree(t *testing.T) {
	conf := Readconfig("test_conf.yml")

	rb := Createrulebase(conf)
	subject := "John"
	url := "www.corpA.com/*"

	val, err := TreeLookup(subject, url, rb)
	if err != nil {
		t.Fatalf("Error %s\n", err)
	}
	if val != 1 {
		t.Errorf("Whrong authorization value for %s:%s\n", subject, url)
	}
}

var conf *Config
var rb_map map[string]map[string]int
var rb_tree *prefixtree.Tree
var initialized, rb_tree_init, rb_map_init int
var ns int = 100
var nr int = 100

func createconfig(ns int, nr int) *Config {
	var conf Config
	var sub, str string
	var acl map[string]string
	var r int64

	acl = make(map[string]string)

	conf.Title = "This is a test rulebase"

	for s := 0; s < ns; s++ {
		sub = RandString(12)
		for i := 0; i < nr; i++ {
			r = src.Int63()
			str, _ = GenerateRandomString(int(r % 120))
			if r%1 == 1 {
				// acl[RandString(int(r%120))] = "allow"
				acl[str] = "allow"
			} else {
				// acl[RandString(int(r%120))] = "deny"
				acl[str] = "deny"
			}
		}

		// rule := {Subject: sub, ACL, acl}
		conf.Rules = append(conf.Rules, struct {
			Subject string            "yaml:\"subject\""
			ACL     map[string]string "yaml:\"ACL\""
		}{sub, acl})
	}

	return &conf
}

func BenchmarkLookupTree(b *testing.B) {
	var subject, url string
	// var access string

	if initialized == 0 {
		b.Log("Initialising...")
		conf = createconfig(ns, nr)
		initialized = 1
	}

	if rb_tree_init == 0 {
		rb_tree = Createrulebase(conf)
		rb_tree_init = 1
	}

	//Pick a random rule ...
	rule := conf.Rules[int(src.Int63())%ns]
	subject = rule.Subject
	//Pick a random ACL rule from that rule
	i := int(src.Int63()) % nr
	for k, _ := range rule.ACL {
		i--
		if i == 0 {
			url = k
			// access = v
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		TreeLookup(subject, url, rb_tree)
	}

}

func BenchmarkLookupMap(b *testing.B) {
	var subject, url string
	// var access string

	if initialized == 0 {
		b.Log("Initialising...")
		conf = createconfig(ns, nr)
		initialized = 1
	}
	if rb_map_init == 0 {
		rb_map = Maprulebase(conf)
		rb_map_init = 1
	}

	//Pick a random rule ...
	rule := conf.Rules[int(src.Int63())%ns]
	subject = rule.Subject
	//Pick a random ACL rule from that rule
	i := int(src.Int63()) % nr
	for k, _ := range rule.ACL {
		i--
		if i == 0 {
			url = k
			// access = v
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		MapLookup(subject, url, rb_map)
	}

}
