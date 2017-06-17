package hashids

import (
	"math"
	"reflect"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	hdata := NewData()
	hdata.MinLength = 30
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)

	numbers := []int{45, 434, 1313, 99}
	hash, err := hid.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}
	dec := hid.Decode(hash)

	t.Logf("%v -> %v -> %v", numbers, hash, dec)

	if !reflect.DeepEqual(dec, numbers) {
		t.Errorf("Decoded numbers `%v` did not match with original `%v`", dec, numbers)
	}
}

func TestEncodeDecodeInt64(t *testing.T) {
	hdata := NewData()
	hdata.MinLength = 30
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)

	numbers := []int64{45, 434, 1313, 99, math.MaxInt64}
	hash, err := hid.EncodeInt64(numbers)
	if err != nil {
		t.Fatal(err)
	}
	dec := hid.DecodeInt64(hash)

	t.Logf("%v -> %v -> %v", numbers, hash, dec)

	if !reflect.DeepEqual(dec, numbers) {
		t.Errorf("Decoded numbers `%v` did not match with original `%v`", dec, numbers)
	}
}

func TestEncodeWithKnownHash(t *testing.T) {
	hdata := NewData()
	hdata.MinLength = 0
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)

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

func TestDecodeWithKnownHash(t *testing.T) {
	hdata := NewData()
	hdata.MinLength = 0
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)

	hash := "7nnhzEsDkiYa"
	numbers := hid.Decode(hash)

	t.Logf("%v -> %v", hash, numbers)

	expected := []int{45, 434, 1313, 99}
	if !reflect.DeepEqual(numbers, expected) {
		t.Errorf("Decoded numbers `%v` did not match with expected `%v`", numbers, expected)
	}
}

func TestDefaultLength(t *testing.T) {
	hdata := NewData()
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)

	numbers := []int{45, 434, 1313, 99}
	hash, err := hid.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}
	dec := hid.Decode(hash)

	t.Logf("%v -> %v -> %v", numbers, hash, dec)

	if !reflect.DeepEqual(dec, numbers) {
		t.Errorf("Decoded numbers `%v` did not match with original `%v`", dec, numbers)
	}
}

func TestMinLength(t *testing.T) {
	hdata := NewData()
	hdata.Salt = "salt1"
	hdata.MinLength = 10
	hid, _ := NewWithData(hdata)
	hid.Encode([]int{0})
}

func TestCustomAlphabet(t *testing.T) {
	hdata := NewData()
	hdata.Alphabet = "PleasAkMEFoThStx"
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)

	numbers := []int{45, 434, 1313, 99}
	hash, err := hid.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}
	dec := hid.Decode(hash)

	t.Logf("%v -> %v -> %v", numbers, hash, dec)

	if !reflect.DeepEqual(dec, numbers) {
		t.Errorf("Decoded numbers `%v` did not match with original `%v`", dec, numbers)
	}
}

func TestDecodeWithError(t *testing.T) {
	hdata := NewData()
	hdata.Alphabet = "PleasAkMEFoThStx"
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)
	// hash now contains a letter not in the alphabet
	dec, err := hid.DecodeWithError("MAkhkloFAxAoskaZ")

	if dec != nil {
		t.Error("Expected `nil` but got `%v`", dec)
	}
	expected := "alphabet used for hash was different"
	if err == nil || err.Error() != expected {
		t.Errorf("Expected error `%s` but got `%s`", expected, err)
	}
}

// tests issue #28
func TestDecodeWithWrongSalt(t *testing.T) {
	hdata := NewData()
	hdata.Alphabet = "PleasAkMEFoThStx"
	hdata.Salt = "temp"

	hidEncode, _ := NewWithData(hdata)

	numbers := []int{45, 434, 1313, 99}
	hash, _ := hidEncode.Encode(numbers)

	hdata.Salt = "test"
	hidDecode, _ := NewWithData(hdata)
	dec, err := hidDecode.DecodeWithError(hash)

	t.Logf("%v -> %v -> %v", numbers, hash, dec)

	expected := "mismatch between encode and decode: ePaTMalsPMPlhxMl start MEhloASEPosaE re-encoded. result: [7 199 245 19]"
	if err == nil || err.Error() != expected {
		t.Errorf("Expected error `%s` but got `%s`", expected, err)
	}
}

func TestAllocationsPerEncode(t *testing.T) {
	hdata := NewData()
	hdata.Salt = "temp"
	hdata.MinLength = 0
	hid, _ := NewWithData(hdata)

	numbers := []int64{math.MaxInt64, 0, 1024, math.MaxInt64 / 2}
	allocsPerRun := testing.AllocsPerRun(100, func() {
		_, err := hid.EncodeInt64(numbers)
		if err != nil {
			t.Errorf("Unexpected error encoding test data: %s, %v", err, numbers)
		}
	})
	if allocsPerRun != 19 {
		t.Errorf("Expected 19 allocations, got %v ", allocsPerRun)
	}
}
