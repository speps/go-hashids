package hashids

import "testing"

func TestSmallAlphabet(t *testing.T) {
	hdata := NewData()
	hdata.Alphabet = "1234567890"
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected an error as alphabet was too small")
		}
	}()
	NewWithData(hdata)
}

func TestSpacesInAlphabet(t *testing.T) {
	hdata := NewData()
	hdata.Alphabet = "a cdefghijklmnopqrstuvwxyz"
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected an error as spaces in alphabet are not allowed")
		}
	}()
	NewWithData(hdata)
}

func TestNilWithEncode(t *testing.T) {
	h := New()
	v, err := h.Encode(nil)

	if v != "" {
		t.Errorf("Expected empty string and got `%s`", v)
	}
	if err != nil {
		t.Errorf("Expected empty slice and got error `%s`", err)
	}
}

func TestEmptySliceWithEncode(t *testing.T) {
	h := New()
	v, err := h.Encode([]int{})

	if v != "" {
		t.Errorf("Expected empty string and got `%s`", v)
	}
	if err != nil {
		t.Errorf("Expected empty slice and got error `%s`", err)
	}
}

func TestNegativeNumberWithEncode(t *testing.T) {
	h := New()
	v, err := h.Encode([]int{-1})

	if v != "" {
		t.Errorf("Expected empty string and got `%s`", v)
	}
	if err != nil {
		t.Errorf("Expected empty string and got error `%s`", err)
	}
}

func TestEmptySliceWithDecode(t *testing.T) {
	h := New()
	v, err := h.Encode([]int{})

	if v != "" {
		t.Errorf("Expected empty string and got `%s`", v)
	}
	if err != nil {
		t.Errorf("Expected empty slice and got error `%s`", err)
	}
}
