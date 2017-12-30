package prefixtree

import (
	"os"
	"testing"
)

var prefix, index []string
var value []int
var tree *Tree

func TestMain(m *testing.M) {

	prefix = append(prefix, "www.corpA.com/*")
	index = append(index, "John")
	value = append(value, 100)
	prefix = append(prefix, "www.corpA.com/*")
	index = append(index, "Jim")
	value = append(value, 100)
	prefix = append(prefix, "www.corpA.com/admin")
	index = append(index, "John")
	value = append(value, 100)

	tree = New()
	for i := 0; i < len(prefix); i++ {
		tree.Add(prefix[i], index[i], value[i])
	}

	code := m.Run()
	os.Exit(code)
}

func TestGet(t *testing.T) {
	url := prefix[0]
	subject := index[0]
	value := value[0]
	// tree.Add(url, subject, value)

	v, _ := tree.Get(url, subject)

	if v != value {
		t.Errorf("Value of %s:%s should be %d", url, subject, value)
	}
}

func TestMatch(t *testing.T) {
	url := "www.corpA.com/someresource"
	subject := index[0]
	value := value[0]
	// tree.Add(url, subject, value)

	v, err := tree.Match(url, subject)

	if err != nil {
		t.Errorf("Error: %s", err)
	} else {
		if v != value {
			t.Errorf("Value of %s:%s should be %d", url, subject, value)
		}
	}
}

func TestGetNonexistent(t *testing.T) {
	var err error

	url := prefix[0]
	nox_url := "www.corpC.com/*"
	subject := index[0]
	nox_subject := "nox_subject"
	// value := 100
	// tree.Add(url, subject, value)

	_, err = tree.Get(nox_url, subject)
	if err == nil {
		t.Errorf("No error returned although key doesn't exist")
	}

	_, err = tree.Get(url, nox_subject)
	if err == nil {
		t.Errorf("No error returned although index doesn't exist")
	}

}

//TODO: Add a test where a number of random strings of random
//length are added and then one of those is looked up
