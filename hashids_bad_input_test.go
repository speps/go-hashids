package hashids

import (
	"testing"
)

func TestSmallAlphabet(t *testing.T) {
	hdata := NewData()
	hdata.Alphabet = "1234567890"
	_, err := NewWithData(hdata)
	expected := "alphabet must contain at least 16 characters"
	if err == nil || err.Error() != expected {
		t.Errorf("Expected error `%s` but got `%s`", expected, err)
	}
}

func TestSpacesInAlphabet(t *testing.T) {
	hdata := NewData()
	hdata.Alphabet = "a cdefghijklmnopqrstuvwxyz"
	_, err := NewWithData(hdata)
	expected := "alphabet may not contain spaces"
	if err == nil || err.Error() != expected {
		t.Errorf("Expected error `%s` but got `%s`", expected, err)
	}
}

func TestNilWithEncode(t *testing.T) {
	h, _ := New()
	_, err := h.Encode(nil)
	expected := "encoding empty array of numbers makes no sense"
	if err == nil || err.Error() != expected {
		t.Errorf("Expected error `%s` but got `%s`", expected, err)
	}
}

func TestEmptySliceWithEncode(t *testing.T) {
	h, _ := New()
	_, err := h.Encode([]int{})
	expected := "encoding empty array of numbers makes no sense"
	if err == nil || err.Error() != expected {
		t.Errorf("Expected error `%s` but got `%s`", expected, err)
	}
}

func TestNegativeNumberWithEncode(t *testing.T) {
	h, _ := New()
	_, err := h.Encode([]int{-1})
	expected := "negative number not supported"
	if err == nil || err.Error() != expected {
		t.Errorf("Expected error `%s` but got `%s`", expected, err)
	}
}
