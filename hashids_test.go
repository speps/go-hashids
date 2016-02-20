package hashids

import (
	"math"
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

func TestEncryptDecryptInt64(t *testing.T) {
	hdata := NewData()
	hdata.MinLength = 30
	hdata.Salt = "this is my salt"

	hid := NewWithData(hdata)

	numbers := []int64{45, 434, 1313, 99, math.MaxInt64}
	hash, err := hid.EncodeInt64(numbers)
	if err != nil {
		t.Fatal(err)
	}
	dec := hid.DecodeInt64(hash)

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

func TestEncryptWithKnownHash(t *testing.T) {
	hdata := NewData()
	hdata.MinLength = 0
	hdata.Salt = "this is my salt"

	hid := NewWithData(hdata)

	numbers := []int{45, 434, 1313, 99}
	hash, err := hid.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v -> %v", numbers, hash)

	if hash != "7nnhzEsDkiYa" {
		t.Error("hash does not match expected one")
	}
}

func TestDecryptWithKnownHash(t *testing.T) {
	hdata := NewData()
	hdata.MinLength = 0
	hdata.Salt = "this is my salt"

	hid := NewWithData(hdata)

	hash := "7nnhzEsDkiYa"
	numbers := hid.Decode(hash)

	t.Logf("%v -> %v", hash, numbers)

	expected := []int{45, 434, 1313, 99}
	for i, n := range numbers {
		if n != expected[i] {
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

func TestDecryptWithError(t *testing.T) {
	hdata := NewData()
	hdata.Alphabet = "PleasAkMEFoThStx"
	hdata.Salt = "this is my salt"

	hid := NewWithData(hdata)
	// hash now contains a letter not in the alphabet
	dec, err := hid.DecodeWithError("MAkhkloFAxAoskaZ")

	if dec != nil {
		t.Error("DecryptWithError should have returned nil result")
	}
	if err == nil {
		t.Error("DecryptWithError should have returned error")
	}
}
