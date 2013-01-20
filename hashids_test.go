package hashids

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	hid := New()
	hid.MinLength = 30
	hid.Salt = "this is my salt"

	numbers := []int{45, 434, 1313, 99}
	hash := hid.Encrypt(numbers)
	dec := hid.Decrypt(hash)

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

func TestZeroMinimumLength(t *testing.T) {
	hid := New()
	hid.Salt = "this is my salt"

	numbers := []int{45, 434, 1313, 99}
	hash := hid.Encrypt(numbers)
	dec := hid.Decrypt(hash)

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
