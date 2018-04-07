// Go implementation of http://www.hashids.org under MIT license
// Generates hashes from an array of integers, eg. for YouTube like hashes
// Setup: go get github.com/speps/go-hashids
// Original implementations by Ivan Akimov at https://github.com/ivanakimov
// Thanks to Rémy Oudompheng and Peter Hellberg for code review and fixes

package hashids

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

const (
	// Version is the version number of the library
	Version string = "1.0.0"

	// DefaultAlphabet is the default alphabet used by go-hashids
	DefaultAlphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	minAlphabetLength int     = 16
	sepDiv            float64 = 3.5
	guardDiv          float64 = 12.0
)

var sepsOriginal = []rune("cfhistuCFHISTU")

// HashID interface contains only functions needed for encoding/decoding int64s. May be expanded as needed later.
type HashID interface {
	EncodeInt64WithError(numbers []int64) (string, error)
	DecodeInt64WithError(hash string) ([]int64, error)

	EncodeUint64WithError(numbers []uint64) (string, error)
	DecodeUint64WithError(hash string) ([]uint64, error)
}

// hashIDImpl contains everything needed to encode/decode hashids
type hashIDImpl struct {
	alphabet           []rune
	minLength          int
	maxLengthPerNumber int
	salt               []rune
	seps               []rune
	guards             []rune
}

// HashIDData contains the information needed to generate hashids
type HashIDData struct {
	// Alphabet is the alphabet used to generate new ids
	Alphabet string

	// MinLength is the minimum length of a generated id
	MinLength int

	// Salt is the secret used to make the generated id harder to guess
	Salt string
}

// NewData creates a new HashIDData with the DefaultAlphabet already set.
func NewData() *HashIDData {
	return &HashIDData{Alphabet: DefaultAlphabet}
}

// New creates a new HashID
func New() (HashID, error) {
	return NewWithData(NewData())
}

// NewWithData creates a new HashIDImpl with the provided HashIDData
func NewWithData(data *HashIDData) (HashID, error) {
	if len(data.Alphabet) < minAlphabetLength {
		return nil, fmt.Errorf("alphabet must contain at least %d characters", minAlphabetLength)
	}
	if strings.Contains(data.Alphabet, " ") {
		return nil, fmt.Errorf("alphabet may not contain spaces")
	}
	// Check if all characters are unique in Alphabet
	uniqueCheck := make(map[rune]bool, len(data.Alphabet))
	for _, a := range data.Alphabet {
		if _, found := uniqueCheck[a]; found {
			return nil, fmt.Errorf("duplicate character in alphabet: %s", string([]rune{a}))
		}
		uniqueCheck[a] = true
	}

	alphabet := []rune(data.Alphabet)
	salt := []rune(data.Salt)

	seps := duplicateRuneSlice(sepsOriginal)

	// seps should contain only characters present in alphabet; alphabet should not contains seps
	for i := 0; i < len(seps); i++ {
		foundIndex := -1
		for j, a := range alphabet {
			if a == seps[i] {
				foundIndex = j
				break
			}
		}
		if foundIndex == -1 {
			seps = append(seps[:i], seps[i+1:]...)
			i--
		} else {
			alphabet = append(alphabet[:foundIndex], alphabet[foundIndex+1:]...)
		}
	}
	consistentShuffleInPlace(seps, salt)

	if len(seps) == 0 || float64(len(alphabet))/float64(len(seps)) > sepDiv {
		sepsLength := int(math.Ceil(float64(len(alphabet)) / sepDiv))
		if sepsLength == 1 {
			sepsLength++
		}
		if sepsLength > len(seps) {
			diff := sepsLength - len(seps)
			seps = append(seps, alphabet[:diff]...)
			alphabet = alphabet[diff:]
		} else {
			seps = seps[:sepsLength]
		}
	}
	consistentShuffleInPlace(alphabet, salt)

	guardCount := int(math.Ceil(float64(len(alphabet)) / guardDiv))
	var guards []rune
	if len(alphabet) < 3 {
		guards = seps[:guardCount]
		seps = seps[guardCount:]
	} else {
		guards = alphabet[:guardCount]
		alphabet = alphabet[guardCount:]
	}

	hid := &hashIDImpl{
		alphabet:  alphabet,
		minLength: data.MinLength,
		salt:      salt,
		seps:      seps,
		guards:    guards,
	}

	// Calculate the maximum possible string length by hashing the maximum possible id
	encoded, err := hid.EncodeInt64WithError([]int64{math.MaxInt64})
	if err != nil {
		return nil, fmt.Errorf("Unable to encode maximum int64 to find max encoded value length: %s", err)
	}
	hid.maxLengthPerNumber = len(encoded)

	return hid, nil
}

// Encode hashes an array of int to a string containing at least MinLength characters taken from the Alphabet.
// Use Decode using the same Alphabet and Salt to get back the array of int.
func (h *hashIDImpl) Encode(numbers []int) (string, error) {
	numbers64 := make([]int64, 0, len(numbers))
	for _, id := range numbers {
		numbers64 = append(numbers64, int64(id))
	}
	return h.EncodeInt64WithError(numbers64)
}

// EncodeInt64WithError hashes an array of int64 to a string containing at least MinLength characters taken from the Alphabet.
// Use DecodeInt64 using the same Alphabet and Salt to get back the array of int64.
func (h *hashIDImpl) EncodeInt64WithError(numbers []int64) (string, error) {
	if len(numbers) == 0 {
		return "", errors.New("encoding empty array of numbers makes no sense")
	}
	for _, n := range numbers {
		if n < 0 {
			return "", errors.New("negative number not supported")
		}
	}

	alphabet := duplicateRuneSlice(h.alphabet)

	numbersHash := int64(0)
	for i, n := range numbers {
		numbersHash += (n % int64(i+100))
	}

	maxRuneLength := h.maxLengthPerNumber * len(numbers)
	if maxRuneLength < h.minLength {
		maxRuneLength = h.minLength
	}

	result := make([]rune, 0, maxRuneLength)
	lottery := alphabet[numbersHash%int64(len(alphabet))]
	result = append(result, lottery)
	hashBuf := make([]rune, maxRuneLength)
	buffer := make([]rune, len(alphabet)+len(h.salt)+1)

	for i, n := range numbers {
		buffer = buffer[:1]
		buffer[0] = lottery
		buffer = append(buffer, h.salt...)
		buffer = append(buffer, alphabet...)
		consistentShuffleInPlace(alphabet, buffer[:len(alphabet)])
		hashBuf = hash(n, alphabet, hashBuf)
		result = append(result, hashBuf...)

		if i+1 < len(numbers) {
			n %= int64(hashBuf[0]) + int64(i)
			result = append(result, h.seps[n%int64(len(h.seps))])
		}
	}

	if len(result) < h.minLength {
		guardIndex := (numbersHash + int64(result[0])) % int64(len(h.guards))
		result = append([]rune{h.guards[guardIndex]}, result...)

		if len(result) < h.minLength {
			guardIndex = (numbersHash + int64(result[2])) % int64(len(h.guards))
			result = append(result, h.guards[guardIndex])
		}
	}

	halfLength := len(alphabet) / 2
	for len(result) < h.minLength {
		consistentShuffleInPlace(alphabet, duplicateRuneSlice(alphabet))
		result = append(alphabet[halfLength:], append(result, alphabet[:halfLength]...)...)
		excess := len(result) - h.minLength
		if excess > 0 {
			result = result[excess/2 : excess/2+h.minLength]
		}
	}

	return string(result), nil
}

func (h *hashIDImpl) EncodeUint64WithError(numbers []uint64) (string, error) {
	if len(numbers) == 0 {
		return "", errors.New("encoding empty array of numbers makes no sense")
	}

	alphabet := duplicateRuneSlice(h.alphabet)

	numbersHash := uint64(0)
	for i, n := range numbers {
		numbersHash += n % uint64(i + 100)
	}

	maxRuneLength := h.maxLengthPerNumber * len(numbers)
	if maxRuneLength < h.minLength {
		maxRuneLength = h.minLength
	}

	result := make([]rune, 0, maxRuneLength)
	lottery := alphabet[numbersHash % uint64(len(alphabet))]
	result = append(result, lottery)
	hashBuf := make([]rune, maxRuneLength)
	buffer := make([]rune, len(alphabet)+len(h.salt) + 1)

	for i, n := range numbers {
		buffer = buffer[:1]
		buffer[0] = lottery
		buffer = append(buffer, h.salt...)
		buffer = append(buffer, alphabet...)
		consistentShuffleInPlace(alphabet, buffer[:len(alphabet)])
		hashBuf = hashUnsigned(n, alphabet, hashBuf)
		result = append(result, hashBuf...)

		if i+1 < len(numbers) {
			n %= uint64(hashBuf[0]) + uint64(i)
			result = append(result, h.seps[n % uint64(len(h.seps))])
		}
	}

	if len(result) < h.minLength {
		guardIndex := (numbersHash + uint64(result[0])) % uint64(len(h.guards))
		result = append([]rune{h.guards[guardIndex]}, result...)

		if len(result) < h.minLength {
			guardIndex = (numbersHash + uint64(result[2])) % uint64(len(h.guards))
			result = append(result, h.guards[guardIndex])
		}
	}

	halfLength := len(alphabet) / 2
	for len(result) < h.minLength {
		consistentShuffleInPlace(alphabet, duplicateRuneSlice(alphabet))
		result = append(alphabet[halfLength:], append(result, alphabet[:halfLength]...)...)
		excess := len(result) - h.minLength
		if excess > 0 {
			result = result[excess/2 : excess/2 + h.minLength]
		}
	}

	return string(result), nil
}

// DEPRECATED: Use DecodeWithError instead
// Decode unhashes the string passed to an array of int.
// It is symmetric with Encode if the Alphabet and Salt are the same ones which were used to hash.
// MinLength has no effect on Decode.
func (h *hashIDImpl) Decode(hash string) []int {
	result, err := h.DecodeWithError(hash)
	if err != nil {
		panic(err)
	}
	return result
}

// Decode unhashes the string passed to an array of int.
// It is symmetric with Encode if the Alphabet and Salt are the same ones which were used to hash.
// MinLength has no effect on Decode.
func (h *hashIDImpl) DecodeWithError(hash string) ([]int, error) {
	result64, err := h.DecodeInt64WithError(hash)
	if err != nil {
		return nil, err
	}
	result := make([]int, 0, len(result64))
	for _, id := range result64 {
		result = append(result, int(id))
	}
	return result, nil
}

// DEPRECATED: Use DecodeInt64WithError instead
// DecodeInt64 unhashes the string passed to an array of int64.
// It is symmetric with EncodeInt64WithError if the Alphabet and Salt are the same ones which were used to hash.
// MinLength has no effect on DecodeInt64.
func (h *hashIDImpl) DecodeInt64(hash string) []int64 {
	result, err := h.DecodeInt64WithError(hash)
	if err != nil {
		panic(err)
	}
	return result
}

// DecodeInt64 unhashes the string passed to an array of int64.
// It is symmetric with EncodeInt64WithError if the Alphabet and Salt are the same ones which were used to hash.
// MinLength has no effect on DecodeInt64.
func (h *hashIDImpl) DecodeInt64WithError(hash string) ([]int64, error) {
	hashes := splitRunes([]rune(hash), h.guards)
	hashIndex := 0
	if len(hashes) == 2 || len(hashes) == 3 {
		hashIndex = 1
	}

	result := make([]int64, 0, 10)

	hashBreakdown := hashes[hashIndex]
	if len(hashBreakdown) > 0 {
		lottery := hashBreakdown[0]
		hashBreakdown = hashBreakdown[1:]
		hashes = splitRunes(hashBreakdown, h.seps)
		alphabet := duplicateRuneSlice(h.alphabet)
		buffer := make([]rune, len(alphabet)+len(h.salt)+1)
		for _, subHash := range hashes {
			buffer = buffer[:1]
			buffer[0] = lottery
			buffer = append(buffer, h.salt...)
			buffer = append(buffer, alphabet...)
			consistentShuffleInPlace(alphabet, buffer[:len(alphabet)])
			number, err := unhash(subHash, alphabet)
			if err != nil {
				return nil, err
			}
			result = append(result, number)
		}
	}

	sanityCheck, _ := h.EncodeInt64WithError(result)
	if sanityCheck != hash {
		return result, fmt.Errorf("mismatch between encode and decode: %s start %s"+
			" re-encoded. result: %v", hash, sanityCheck, result)
	}

	return result, nil
}

func (h *hashIDImpl) DecodeUint64WithError(hash string) ([]uint64, error) {
	hashes := splitRunes([]rune(hash), h.guards)
	hashIndex := 0
	if len(hashes) == 2 || len(hashes) == 3 {
		hashIndex = 1
	}

	result := make([]uint64, 0, 10)

	hashBreakdown := hashes[hashIndex]
	if len(hashBreakdown) > 0 {
		lottery := hashBreakdown[0]
		hashBreakdown = hashBreakdown[1:]
		hashes = splitRunes(hashBreakdown, h.seps)
		alphabet := duplicateRuneSlice(h.alphabet)
		buffer := make([]rune, len(alphabet)+len(h.salt)+1)
		for _, subHash := range hashes {
			buffer = buffer[:1]
			buffer[0] = lottery
			buffer = append(buffer, h.salt...)
			buffer = append(buffer, alphabet...)
			consistentShuffleInPlace(alphabet, buffer[:len(alphabet)])
			number, err := unhashUnsigned(subHash, alphabet)
			if err != nil {
				return nil, err
			}
			result = append(result, number)
		}
	}

	sanityCheck, _ := h.EncodeUint64WithError(result)
	if sanityCheck != hash {
		return result, fmt.Errorf("mismatch between encode and decode: %s start %s"+
			" re-encoded. result: %v", hash, sanityCheck, result)
	}

	return result, nil
}

func splitRunes(input, seps []rune) [][]rune {
	splitIndices := make([]int, 0)
	for i, inputRune := range input {
		for _, sepsRune := range seps {
			if inputRune == sepsRune {
				splitIndices = append(splitIndices, i)
			}
		}
	}

	result := make([][]rune, 0, len(splitIndices)+1)
	inputLeft := input[:]
	for _, splitIndex := range splitIndices {
		splitIndex -= len(input) - len(inputLeft)
		result = append(result, inputLeft[:splitIndex])
		inputLeft = inputLeft[splitIndex+1:]
	}
	result = append(result, inputLeft)

	return result
}

func hash(input int64, alphabet []rune, result []rune) []rune {
	result = result[:0]
	for {
		r := alphabet[input%int64(len(alphabet))]
		result = append(result, r)
		input /= int64(len(alphabet))
		if input == 0 {
			break
		}
	}
	for i := len(result)/2 - 1; i >= 0; i-- {
		opp := len(result) - 1 - i
		result[i], result[opp] = result[opp], result[i]
	}
	return result
}

func hashUnsigned(input uint64, alphabet []rune, result []rune) []rune {
	result = result[:0]
	for {
		r := alphabet[input % uint64(len(alphabet))]
		result = append(result, r)
		input /= uint64(len(alphabet))
		if input == 0 {
			break
		}
	}
	for i := len(result)/2 - 1; i >= 0; i-- {
		opp := len(result) - 1 - i
		result[i], result[opp] = result[opp], result[i]
	}
	return result
}

func unhash(input, alphabet []rune) (int64, error) {
	result := int64(0)
	for _, inputRune := range input {
		alphabetPos := -1
		for pos, alphabetRune := range alphabet {
			if inputRune == alphabetRune {
				alphabetPos = pos
				break
			}
		}
		if alphabetPos == -1 {
			return 0, errors.New("alphabet used for hash was different")
		}

		result = result*int64(len(alphabet)) + int64(alphabetPos)
	}
	return result, nil
}

func unhashUnsigned(input, alphabet []rune) (uint64, error) {
	result := uint64(0)
	for _, inputRune := range input {
		alphabetPos := -1
		for pos, alphabetRune := range alphabet {
			if inputRune == alphabetRune {
				alphabetPos = pos
				break
			}
		}
		if alphabetPos == -1 {
			return 0, errors.New("alphabet used for hash was different")
		}

		result = result * uint64(len(alphabet)) + uint64(alphabetPos)
	}
	return result, nil
}

func consistentShuffle(alphabet, salt []rune) []rune {
	if len(salt) == 0 {
		return alphabet
	}

	result := duplicateRuneSlice(alphabet)
	for i, v, p := len(result)-1, 0, 0; i > 0; i-- {
		p += int(salt[v])
		j := (int(salt[v]) + v + p) % i
		result[i], result[j] = result[j], result[i]
		v = (v + 1) % len(salt)
	}

	return result
}

func consistentShuffleInPlace(alphabet []rune, salt []rune) {
	if len(salt) == 0 {
		return
	}

	for i, v, p := len(alphabet)-1, 0, 0; i > 0; i-- {
		p += int(salt[v])
		j := (int(salt[v]) + v + p) % i
		alphabet[i], alphabet[j] = alphabet[j], alphabet[i]
		v = (v + 1) % len(salt)
	}
	return
}

func duplicateRuneSlice(data []rune) []rune {
	result := make([]rune, len(data))
	copy(result, data)
	return result
}
