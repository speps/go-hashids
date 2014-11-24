package hashids

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	hdata := NewData()
	hdata.MinLength = 30
	hdata.Salt = "this is my salt"

	hid := NewWithData(hdata)

	numbers := []int{45, 434, 1313, 99}
	hash, err := hid.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}
	dec := hid.Decode(hash)

	t.Logf("%v -> %v -> %v", numbers, hash, dec)

	if len(numbers) != len(dec) {
		t.Error("lengths do not match")
	}

	for i, n := range numbers {
		if n != dec[i] {
			t.Fail()
		}
	}
}

func TestDefaultLength(t *testing.T) {
	hdata := NewData()
	hdata.Salt = "this is my salt"

	hid := NewWithData(hdata)

	numbers := []int{45, 434, 1313, 99}
	hash, err := hid.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}
	dec := hid.Decode(hash)

	t.Logf("%v -> %v -> %v", numbers, hash, dec)

	if len(numbers) != len(dec) {
		t.Error("lengths do not match")
	}

	for i, n := range numbers {
		if n != dec[i] {
			t.Fail()
		}
	}
}

func TestMinLength(t *testing.T) {
	hdata := NewData()
	hdata.Salt = "salt1"
	hdata.MinLength = 10
	hid := NewWithData(hdata)
	hid.Encode([]int{0})
}

func TestCustomAlphabet(t *testing.T) {
	hdata := NewData()
	hdata.Alphabet = "PleasAkMEFoThStx"
	hdata.Salt = "this is my salt"

	hid := NewWithData(hdata)

	numbers := []int{45, 434, 1313, 99}
	hash, err := hid.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}
	dec := hid.Decode(hash)

	t.Logf("%v -> %v -> %v", numbers, hash, dec)

	if len(numbers) != len(dec) {
		t.Error("lengths do not match")
	}

	for i, n := range numbers {
		if n != dec[i] {
			t.Fail()
		}
	}
}
