package hashids

import (
	"reflect"
	"testing"
)

func testMinLength(minLength int, t *testing.T) {
	numbers := []int{1, 2, 3}
	hdata := NewData()
	hdata.MinLength = minLength
	h, _ := NewWithData(hdata)
	e, err := h.Encode(numbers)
	decodedNumbers := h.Decode(e)

	if err != nil {
		t.Errorf("Expected no error but got `%s`", err)
	}
	if len(e) < minLength {
		t.Errorf("Expected hash length to be at least `%d`, was `%d`", minLength, len(e))
	}
	if !reflect.DeepEqual(decodedNumbers, numbers) {
		t.Errorf("Decoded numbers `%v` did not match with original `%v`", decodedNumbers, numbers)
	}
}

func TestShouldWorkWhen0(t *testing.T) {
	testMinLength(0, t)
}

func TestShouldWorkWhen1(t *testing.T) {
	testMinLength(1, t)
}

func TestShouldWorkWhen10(t *testing.T) {
	testMinLength(10, t)
}

func TestShouldWorkWhen999(t *testing.T) {
	testMinLength(999, t)
}

func TestShouldWorkWhen1000(t *testing.T) {
	testMinLength(1000, t)
}
