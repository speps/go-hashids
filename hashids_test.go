package hashids

import (
	"github.com/go-test/deep"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	hdata := NewData()
	hdata.MinLength = 30
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)

	numbers := []int64{45, 434, 1313, 99}
	hash, err := hid.EncodeInt64WithError(numbers)
	if err != nil {
		t.Fatal(err)
	}
	dec, _ := hid.DecodeInt64WithError(hash)

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
	hash, err := hid.EncodeInt64WithError(numbers)
	if err != nil {
		t.Fatal(err)
	}
	dec, _ := hid.DecodeInt64WithError(hash)

	t.Logf("%v -> %v -> %v", numbers, hash, dec)

	if !reflect.DeepEqual(dec, numbers) {
		t.Errorf("Decoded numbers `%v` did not match with original `%v`", dec, numbers)
	}
}

func TestEncodeDecodeUint64(t *testing.T) {
	params := gopter.DefaultTestParameters()
	params.MinSuccessfulTests = 2000
	properties := gopter.NewProperties(params)

	properties.Property("Encoding and decoding are symmetrical for uint64s", prop.ForAll(
		func(minLength int, salt string, numbers []uint64) string {
			hdata := NewData()
			hdata.MinLength = minLength
			hdata.Salt = salt

			hid, err := NewWithData(hdata)
			if err != nil {
				return err.Error()
			}

			s, err := hid.EncodeUint64WithError(numbers)
			if err != nil {
				return err.Error()
			}

			roundtripped, err := hid.DecodeUint64WithError(s)
			if err != nil {
				return err.Error()
			}

			if diff := deep.Equal(numbers, roundtripped); diff != nil {
				return strings.Join(diff, "\n")
			}
			return ""
		}, gen.IntRange(0, math.MaxInt8), gen.AnyString(), gen.SliceOf(gen.UInt64()).SuchThat(
			func(slice []uint64) bool {
				return len(slice) > 0
			})))

	properties.TestingRun(t)
}

func TestEncodeWithKnownHash(t *testing.T) {
	hdata := NewData()
	hdata.MinLength = 0
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)

	numbers := []int64{45, 434, 1313, 99}
	hash, err := hid.EncodeInt64WithError(numbers)
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
	numbers, _ := hid.DecodeInt64WithError(hash)

	t.Logf("%v -> %v", hash, numbers)

	expected := []int64{45, 434, 1313, 99}
	if !reflect.DeepEqual(numbers, expected) {
		t.Errorf("Decoded numbers `%v` did not match with expected `%v`", numbers, expected)
	}
}

func TestDefaultLength(t *testing.T) {
	hdata := NewData()
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)

	numbers := []int64{45, 434, 1313, 99}
	hash, err := hid.EncodeInt64WithError(numbers)
	if err != nil {
		t.Fatal(err)
	}
	dec, _ := hid.DecodeInt64WithError(hash)

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
	hid.EncodeInt64WithError([]int64{0})
}

func TestCustomAlphabet(t *testing.T) {
	hdata := NewData()
	hdata.Alphabet = "PleasAkMEFoThStx"
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)

	numbers := []int64{45, 434, 1313, 99}
	hash, err := hid.EncodeInt64WithError(numbers)

	if err != nil {
		t.Fatal(err)
	}
	dec, _ := hid.DecodeInt64WithError(hash)

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
	dec, err := hid.DecodeInt64WithError("MAkhkloFAxAoskaZ")

	if dec != nil {
		t.Errorf("Expected `nil` but got `%v`", dec)
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

	numbers := []int64{45, 434, 1313, 99}
	hash, _ := hidEncode.EncodeInt64WithError(numbers)

	hdata.Salt = "test"
	hidDecode, _ := NewWithData(hdata)
	dec, err := hidDecode.DecodeInt64WithError(hash)

	t.Logf("%v -> %v -> %v", numbers, hash, dec)

	expected := "mismatch between encode and decode: ePaTMalsPMPlhxMl start MEhloASEPosaE re-encoded. result: [7 199 245 19]"
	if err == nil || err.Error() != expected {
		t.Errorf("Expected error `%s` but got `%s`", expected, err)
	}
}

func checkAllocations(t *testing.T, hid HashID, values []int64, expectedAllocations float64) {
	allocsPerRun := testing.AllocsPerRun(5, func() {
		_, err := hid.EncodeInt64WithError(values)
		if err != nil {
			t.Errorf("Unexpected error encoding test data: %s, %v", err, values)
		}
	})
	if allocsPerRun != expectedAllocations {
		t.Errorf("Expected %v allocations, got %v ", expectedAllocations, allocsPerRun)
	}
}

func TestAllocationsPerEncodeTypical(t *testing.T) {
	hdata := NewData()
	hdata.Salt = "temp"
	hdata.MinLength = 0
	hid, _ := NewWithData(hdata)

	singleNumber := []int64{42}

	maxNumbers := []int64{math.MaxInt64, math.MaxInt64, math.MaxInt64, math.MaxInt64}
	minNumbers := []int64{0, 0, 0, 0}
	mixNubers := []int64{math.MaxInt64, 0, 1024, math.MaxInt64 / 2}

	checkAllocations(t, hid, singleNumber, 5)

	// Same length, same number of allocations
	checkAllocations(t, hid, maxNumbers, 5)
	checkAllocations(t, hid, minNumbers, 5)
	checkAllocations(t, hid, mixNubers, 5)

	// Greater length, same number of allocation
	checkAllocations(t, hid, append(maxNumbers, maxNumbers...), 5)
	checkAllocations(t, hid, append(minNumbers, minNumbers...), 5)
	checkAllocations(t, hid, append(mixNubers, mixNubers...), 5)
}

func TestAllocationsPerEncodeNoSalt(t *testing.T) {
	hdata := NewData()
	hdata.Salt = ""
	hdata.MinLength = 0
	hid, _ := NewWithData(hdata)

	singleNumber := []int64{42}

	maxNumbers := []int64{math.MaxInt64, math.MaxInt64, math.MaxInt64, math.MaxInt64}
	minNumbers := []int64{0, 0, 0, 0}
	mixNubers := []int64{math.MaxInt64, 0, 1024, math.MaxInt64 / 2}

	checkAllocations(t, hid, singleNumber, 5)

	// Same length, same number of allocations
	checkAllocations(t, hid, maxNumbers, 5)
	checkAllocations(t, hid, minNumbers, 5)
	checkAllocations(t, hid, mixNubers, 5)

	// Greater length, same number of allocation
	checkAllocations(t, hid, append(maxNumbers, maxNumbers...), 5)
	checkAllocations(t, hid, append(minNumbers, minNumbers...), 5)
	checkAllocations(t, hid, append(mixNubers, mixNubers...), 5)
}

func TestAllocationsPerEncodeMinLength(t *testing.T) {
	hdata := NewData()
	hdata.Salt = "temp"
	hdata.MinLength = 10
	hid, _ := NewWithData(hdata)

	singleNumber := []int64{42}

	maxNumbers := []int64{math.MaxInt64, math.MaxInt64, math.MaxInt64, math.MaxInt64}
	minNumbers := []int64{0, 0, 0, 0}
	mixNubers := []int64{math.MaxInt64, 0, 1024, math.MaxInt64 / 2}

	checkAllocations(t, hid, singleNumber, 9)

	// Same length, same number of allocations
	checkAllocations(t, hid, maxNumbers, 5)
	checkAllocations(t, hid, minNumbers, 6)
	checkAllocations(t, hid, mixNubers, 5)

	// Greater length, same number of allocation
	checkAllocations(t, hid, append(maxNumbers, maxNumbers...), 5)
	checkAllocations(t, hid, append(minNumbers, minNumbers...), 5)
	checkAllocations(t, hid, append(mixNubers, mixNubers...), 5)
}

func TestAllocationsPerEncodeMinLengthHigh(t *testing.T) {
	hdata := NewData()
	hdata.Salt = "temp"
	hdata.MinLength = 100
	hid, _ := NewWithData(hdata)

	singleNumber := []int64{42}

	maxNumbers := []int64{math.MaxInt64, math.MaxInt64, math.MaxInt64, math.MaxInt64}
	minNumbers := []int64{0, 0, 0, 0}
	mixNubers := []int64{math.MaxInt64, 0, 1024, math.MaxInt64 / 2}

	checkAllocations(t, hid, singleNumber, 15)

	// Same length, same number of allocations
	checkAllocations(t, hid, maxNumbers, 12)
	checkAllocations(t, hid, minNumbers, 15)
	checkAllocations(t, hid, mixNubers, 12)

	// Greater length, same number of allocation
	checkAllocations(t, hid, append(maxNumbers, maxNumbers...), 5)
	checkAllocations(t, hid, append(minNumbers, minNumbers...), 12)
	checkAllocations(t, hid, append(mixNubers, mixNubers...), 9)
}

func checkAllocationsDecode(t *testing.T, hid HashID, values []int64, expectedAllocations float64) {
	encoded, err := hid.EncodeInt64WithError(values)
	if err != nil {
		t.Errorf("Unexpected error encoding test data: %s, %v", err, values)
	}
	allocsPerRun := testing.AllocsPerRun(5, func() {
		_, err := hid.DecodeInt64WithError(encoded)
		if err != nil {
			t.Errorf("Unexpected error decoding test data: %s, %v", err, values)
		}
	})
	if allocsPerRun != expectedAllocations {
		t.Errorf("Expected %v allocations, got %v ", expectedAllocations, allocsPerRun)
	}
}

func TestAllocationsDecodeTypical(t *testing.T) {
	hdata := NewData()
	hdata.Salt = "temp"
	hdata.MinLength = 0
	hid, _ := NewWithData(hdata)

	singleNumber := []int64{42}

	maxNumbers := []int64{math.MaxInt64, math.MaxInt64, math.MaxInt64, math.MaxInt64}
	minNumbers := []int64{0, 0, 0, 0}
	mixNubers := []int64{math.MaxInt64, 0, 1024, math.MaxInt64 / 2}

	checkAllocationsDecode(t, hid, singleNumber, 11)

	// Same length, same number of allocations
	checkAllocationsDecode(t, hid, maxNumbers, 14)
	checkAllocationsDecode(t, hid, minNumbers, 14)
	checkAllocationsDecode(t, hid, mixNubers, 14)

	// Greater length, same number of allocation per case. Length is long enough
	// to not fit inisde the pre-allocated result buffer hence one extra alloc
	checkAllocationsDecode(t, hid, append(maxNumbers, maxNumbers...), 15)
	checkAllocationsDecode(t, hid, append(minNumbers, minNumbers...), 15)
	checkAllocationsDecode(t, hid, append(mixNubers, mixNubers...), 15)
}
