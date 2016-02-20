go-hashids [![Build Status](https://ci.appveyor.com/api/projects/status/1s8yeafycpa2vdaq?svg=true)](https://ci.appveyor.com/project/speps/go-hashids) [![GoDoc](https://godoc.org/github.com/speps/go-hashids?status.svg)](https://godoc.org/github.com/speps/go-hashids)
==========

Go (golang) v1 implementation of http://www.hashids.org
under MIT License (same as the original implementations)

Original implementations by [Ivan Akimov](https://github.com/ivanakimov)

### Setup
<pre>go get github.com/speps/go-hashids</pre>

### Example
```go
package main

import "fmt"
import "github.com/speps/go-hashids"

func main() {
    hd := hashids.NewData()
    hd.Salt = "this is my salt"
    hd.MinLength = 30
    h := hashids.NewWithData(hd)
    e, _ := h.Encode([]int{45, 434, 1313, 99})
    fmt.Println(e)
    d, _ := h.DecodeWithError(e)
    fmt.Println(d)
}
```

### Test results

```
=== RUN   TestEncryptDecrypt
--- PASS: TestEncryptDecrypt (0.00s)
        hashids_test.go:22: [45 434 1313 99] -> woQ2vqjnG7nnhzEsDkiYadKa3O71br -> [45 434 1313 99]
=== RUN   TestEncryptDecryptInt64
--- PASS: TestEncryptDecryptInt64 (0.00s)
        hashids_test.go:49: [45 434 1313 99 9223372036854775807] -> ZvGlaahBptQNfPOuPjJ51zO3wVzP01 -> [45 434 1313 99 9223372036854775807]
=== RUN   TestEncryptWithKnownHash
--- PASS: TestEncryptWithKnownHash (0.00s)
        hashids_test.go:75: [45 434 1313 99] -> 7nnhzEsDkiYa
=== RUN   TestDecryptWithKnownHash
--- PASS: TestDecryptWithKnownHash (0.00s)
        hashids_test.go:92: 7nnhzEsDkiYa -> [45 434 1313 99]
=== RUN   TestDefaultLength
--- PASS: TestDefaultLength (0.00s)
        hashids_test.go:115: [45 434 1313 99] -> 7nnhzEsDkiYa -> [45 434 1313 99]
=== RUN   TestMinLength
--- PASS: TestMinLength (0.00s)
=== RUN   TestCustomAlphabet
--- PASS: TestCustomAlphabet (0.00s)
        hashids_test.go:150: [45 434 1313 99] -> MAkhkloFAxAoskax -> [45 434 1313 99]
=== RUN   TestDecryptWithError
--- PASS: TestDecryptWithError (0.00s)
PASS
```

### Thanks to all the contributors

* [Harm Aarts](https://github.com/haarts)
* [Christoffer G. Thomsen](https://github.com/cgt)
* [Peter Hellberg](https://github.com/peterhellberg)
* [Rémy Oudompheng](https://github.com/remyoudompheng)
* [Mart Roosmaa](https://github.com/roosmaa)

Let me know if I forgot anyone of course.

### Changelog

2014/09/13

* Updated to Hashids v1.0.0 (should be compatible with other implementations, let me know if not, was checked against the Javascript version)
* Changed API
    * Encrypt/Decrypt are now Encode/Decode
    * HashID is now constructed from HashIDData containing alphabet, salt and minimum length
