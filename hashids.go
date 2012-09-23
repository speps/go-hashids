// Go implementation of http://www.hashids.org under MIT license
// Generates hashes from an array of integers, eg. for YouTube like hashes
// Setup: go get github.com/speps/go-hashids
// Original implementations by Ivan Akimov at https://github.com/ivanakimov
// Thanks to Rémy Oudompheng for some code review

package hashids

import (
	"bytes"
	"math"
	"strconv"
)

const DefaultAlphabet string = "xcS4F6h89aUbidefI7jkyunopqrsgCYE5GHTKLMtARXz"

var primes []int = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43}
var sepsIndices []int = []int{0, 4, 8, 12}

type HashID struct {
	Alphabet  string
	MinLength int
	Salt      string
}

// New creates a new HashID with the DefaultAlphabet already set.
func New() *HashID {
	return &HashID{Alphabet: DefaultAlphabet}
}

// Encrypt hashes an array of int to a string containing at least MinLength characters taken from the Alphabet.
// Use Decrypt using the same Alphabet and Salt to get back the array of int.
func (h *HashID) Encrypt(numbers []int) string {
	if len(numbers) == 0 {
		panic("encrypting empty array of numbers makes no sense")
	}
	for _, n := range numbers {
		if n < 0 {
			panic("negative number not supported")
		}
	}
	if len(h.Alphabet) < 4 {
		panic("alphabet must contain at least 4 characters")
	}

	alphabetRunes := []rune(h.Alphabet)
	saltRunes := []rune(h.Salt)

	alphabetRunes, seps, guards := getSepsAndGuards(alphabetRunes)

	alphabetRunes = consistentShuffle(alphabetRunes, saltRunes)

	return string(encode(numbers, alphabetRunes, saltRunes, seps, guards, h.MinLength))
}

func encode(numbers []int, alphabetOriginal, salt, sepsOriginal, guards []rune, minLength int) []rune {
	numbersBytes := make([]byte, 0)
	for _, n := range numbers {
		numbersBytes = strconv.AppendInt(numbersBytes, int64(n), 10)
	}
	seps := consistentShuffle(sepsOriginal, bytes.Runes(numbersBytes))

	alphabet := make([]rune, len(alphabetOriginal))
	copy(alphabet, alphabetOriginal)

	lotterySalt := make([]byte, 0, 2*len(numbers))
	for i, n := range numbers {
		if i > 0 {
			lotterySalt = append(lotterySalt, '-')
		}
		lotterySalt = strconv.AppendInt(lotterySalt, int64(n), 10)
	}
	for _, n := range numbers {
		lotterySalt = append(lotterySalt, '-')
		lotterySalt = strconv.AppendInt(lotterySalt, int64((n+1)*2), 10)
	}
	lottery := consistentShuffle(alphabet, bytes.Runes(lotterySalt))
	lotteryRune := lottery[0]

	for i, r := range alphabet {
		if r == lotteryRune {
			alphabet = append([]rune{lotteryRune}, append(alphabet[:i], alphabet[i+1:]...)...)
			break
		}
	}
	saltL := append([]rune(strconv.FormatInt(int64(lotteryRune&12345), 10)), salt...)

	result := make([]rune, 0, minLength)
	result = append(result, lotteryRune)
	for i, n := range numbers {
		alphabet = consistentShuffle(alphabet, saltL)
		hash := hash(n, alphabet)
		result = append(result, hash...)
		if (i + 1) < len(numbers) {
			sepsIndex := (n + i) % len(seps)
			result = append(result, seps[sepsIndex])
		}
	}

	if len(result) < minLength {
		guardIndex := 0
		for i, n := range numbers {
			guardIndex += (i + 1) * n
		}

		guardIndex %= len(guards)
		guard := guards[guardIndex]

		result = append([]rune{guard}, result...)
		if len(result) < minLength {
			guardIndex = (guardIndex + len(result)) % len(guards)
			guard = guards[guardIndex]
			result = append(result, guard)
		}
	}

	for len(result) < minLength {
		padArray := []int{int(alphabet[1]), int(alphabet[0])}
		padLeft := encode(padArray, alphabet, salt, sepsOriginal, guards, 0)
		padArrayRunes := append([]rune(strconv.FormatInt(int64(padArray[0]), 10)), []rune(strconv.FormatInt(int64(padArray[1]), 10))...)
		padRight := encode(padArray, alphabet, padArrayRunes, sepsOriginal, guards, 0)

		result = append(padLeft, append(result, padRight...)...)
		excess := len(result) - minLength
		if excess > 0 {
			result = result[excess/2 : excess/2+minLength]
		}

		alphabet = consistentShuffle(alphabet, append(salt, result...))
	}

	return result
}

// Decrypt unhashes the string passed to an array of int.
// It is symmetric with Encrypt if the Alphabet and Salt are the same ones which were used to hash.
// MinLength has no effect on Decrypt.
func (h *HashID) Decrypt(hash string) []int {
	alphabetRunes := []rune(h.Alphabet)
	saltRunes := []rune(h.Salt)
	hashRunes := []rune(hash)

	alphabetRunes, seps, guards := getSepsAndGuards(alphabetRunes)

	alphabetRunes = consistentShuffle(alphabetRunes, saltRunes)

	return decode(hashRunes, alphabetRunes, saltRunes, seps, guards)
}

func decode(hash, alphabetOriginal, salt, seps, guards []rune) []int {
	hashes := splitRunes(hash, guards)
	hashIndex := 0
	if len(hashes) == 2 || len(hashes) == 3 {
		hashIndex = 1
	} else {
		panic("malformed hash input")
	}

	hashes = splitRunes(hashes[hashIndex], seps)
	lotteryRune := hashes[0][0]
	hashes[0] = hashes[0][1:]

	alphabet := make([]rune, len(alphabetOriginal))
	copy(alphabet, alphabetOriginal)
	for i, r := range alphabet {
		if r == lotteryRune {
			alphabet = append([]rune{lotteryRune}, append(alphabet[:i], alphabet[i+1:]...)...)
			break
		}
	}

	saltL := append([]rune(strconv.FormatInt(int64(lotteryRune&12345), 10)), salt...)

	result := make([]int, len(hashes))
	for i, subHash := range hashes {
		alphabet = consistentShuffle(alphabet, saltL)
		result[i] = unhash(subHash, alphabet)
	}

	return result
}

func getSepsAndGuards(alphabet []rune) ([]rune, []rune, []rune) {
	guards := make([]rune, 0, len(sepsIndices))
	seps := make([]rune, 0, len(alphabet))
	for _, prime := range primes {
		index := prime - 1 - len(seps)
		if index < len(alphabet) {
			seps = append(seps, alphabet[index])
			alphabet = append(alphabet[:index], alphabet[index+1:]...)
		} else {
			break
		}
	}
	for _, index := range sepsIndices {
		if index < len(seps) {
			guards = append(guards, seps[index])
			seps = append(seps[:index], seps[index+1:]...)
		}
	}
	return alphabet, seps, guards
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
		subInput := make([]rune, splitIndex)
		copy(subInput, inputLeft[:splitIndex])
		result = append(result, subInput)
		inputLeft = inputLeft[splitIndex+1:]
	}
	result = append(result, inputLeft)

	return result
}

func hash(input int, alphabet []rune) []rune {
	result := make([]rune, 0)
	for ; input > 0; input /= len(alphabet) {
		r := alphabet[input%len(alphabet)]
		result = append(result, r)
	}
	reversed := make([]rune, len(result))
	for i, r := range result {
		reversed[len(result)-i-1] = r
	}
	return reversed
}

func unhash(input, alphabet []rune) int {
	result := 0
	for i, inputRune := range input {
		alphabetPos := -1
		for pos, alphabetRune := range alphabet {
			if inputRune == alphabetRune {
				alphabetPos = pos
				break
			}
		}
		if alphabetPos == -1 {
			panic("should not happen, alphabet used for hash was different")
		}

		result += alphabetPos * int(math.Pow(float64(len(alphabet)), float64(len(input)-i-1)))
	}
	return result
}

func consistentShuffle(alphabet, salt []rune) []rune {
	sortingArray := make([]int, len(salt))
	for i, saltRune := range salt {
		sortingArray[i] = int(saltRune)
	}
	for i, _ := range sortingArray {
		add := true
		for k, j := i, len(sortingArray)+i-1; k != j; k++ {
			nextIndex := (k + 1) % len(sortingArray)
			if add {
				sortingArray[i] += sortingArray[nextIndex] + (k * i)
			} else {
				sortingArray[i] -= sortingArray[nextIndex]
			}
			add = !add
		}
		if sortingArray[i] < 0 {
			sortingArray[i] = -sortingArray[i]
		}
	}

	alphabetCopy := make([]rune, len(alphabet))
	copy(alphabetCopy, alphabet)
	result := make([]rune, 0, len(alphabet))
	for i := 0; len(alphabetCopy) > 0; i++ {
		pos := sortingArray[i%len(sortingArray)] % len(alphabetCopy)
		result = append(result, alphabetCopy[pos])
		alphabetCopy = append(alphabetCopy[:pos], alphabetCopy[pos+1:]...)
	}
	return result
}
