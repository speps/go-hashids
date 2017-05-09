package hashids

import (
	"reflect"
	"testing"
)

func testAlphabet(alphabet string, t *testing.T) {
	numbers := []int{1, 2, 3}
	hdata := NewData()
	hdata.Alphabet = alphabet
	h, _ := NewWithData(hdata)
	e, err := h.Encode(numbers)
	decodedNumbers := h.Decode(e)

	if err != nil {
		t.Errorf("Expected no error but got `%s`", err)
	}
	if !reflect.DeepEqual(decodedNumbers, numbers) {
		t.Errorf("Decoded numbers `%v` did not match with original `%v`", decodedNumbers, numbers)
	}
}

func TestShouldWorkWithWorst(t *testing.T) {
	testAlphabet("cCsSfFhHuUiItT01", t)
}

func TestShouldWorkWithHalfSeparators(t *testing.T) {
	testAlphabet("abdegjklCFHISTUc", t)
}

func TestShouldWorkWithTwoSeparators(t *testing.T) {
	testAlphabet("abdegjklmnopqrSF", t)
}

func TestShouldWorkWithNoSeparators(t *testing.T) {
	testAlphabet("abdegjklmnopqrvwxyzABDEGJKLMNOPQRVWXYZ1234567890", t)
}

func TestShouldWorkWithSuperlong(t *testing.T) {
	testAlphabet("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890`~!@#$%^&*()-_=+\\|'\";:/?.>,<{[}]", t)
}

func TestShouldWorkWithWeird(t *testing.T) {
	testAlphabet("`~!@#$%^&*()-_=+\\|'\";:/?.>,<{[}]", t)
}
