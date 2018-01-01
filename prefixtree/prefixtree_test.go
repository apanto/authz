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
	value = append(value, 200)
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

	v, _ := tree.Get(url, subject)

	if v != value {
		t.Errorf("Value of %s:%s should be %d", url, subject, value)
	}
}

func TestMatch(t *testing.T) {
	var url, subject string
	var v, val int
	var err error

	//this matches the "www.corpA.com/*" rule and should return the value 200
	url = "www.corpA.com/someresource"
	subject = index[0]
	val = value[0]
	// tree.Add(url, subject, value)

	v, err = tree.Match(url, subject)

	if err != nil {
		t.Errorf("%s:%s should match but it didn't (Error: %s)", subject, url, err)
	} else {
		if v != val {
			t.Errorf("Value of %s:%s should be %d not %d", url, subject, val, v)
		}
	}

	//this is a partial match to "www.corpA.com/admin", but a full match to "www.corpA.com/*"
	//and should return the correct value, 200
	url = "www.corpA.com/add"

	v, err = tree.Match(url, subject)

	if err != nil {
		t.Errorf("%s:%s should match but it didn't (Error: %s)", subject, url, err)
	} else {
		if v != val {
			t.Errorf("Value of %s:%s should be %d not %d", url, subject, val, v)
		}
	}

	//this matches the "www.corpA.com/admin" and "www.corpA.com/*" rules. Since "www.corpA.com/admin"
	//is the longest match it should match that and return the value 100
	url = prefix[2]
	subject = index[2]
	val = value[2]

	v, err = tree.Match(url, subject)

	if err != nil {
		t.Errorf("%s:%s should match but it didn't (Error: %s)", subject, url, err)
	} else {
		if v != val {
			t.Errorf("Value of %s:%s should be %d not %d", url, subject, val, v)
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
