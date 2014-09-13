go-hashids
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
    e := h.Encode([]int{45, 434, 1313, 99})
    fmt.Println(e)
    d := h.Decode(e)
    fmt.Println(d)
}
```

### Test results

```
=== RUN TestEncryptDecrypt
--- PASS: TestEncryptDecrypt (0.00 seconds)
    hashids_test.go:21: [45 434 1313 99] -> woQ2vqjnG7nnhzEsDkiYadKa3O71br -> [45 434 1313 99]
=== RUN TestDefaultLength
--- PASS: TestDefaultLength (0.00 seconds)
    hashids_test.go:47: [45 434 1313 99] -> 7nnhzEsDkiYa -> [45 434 1313 99]
=== RUN TestCustomAlphabet
--- PASS: TestCustomAlphabet (0.00 seconds)
    hashids_test.go:74: [45 434 1313 99] -> MAkhkloFAxAoskax -> [45 434 1313 99]
PASS
```

### Thanks to all the contributors

* [Harm Aarts](https://github.com/haarts)
* [Christoffer G. Thomsen](https://github.com/cgt)
* [Peter Hellberg](https://github.com/peterhellberg)
* [Rémy Oudompheng](https://github.com/remyoudompheng)

Let me know if I forgot anyone of course.

### Changelog

2014/09/13

* Updated to Hashids v1.0.0 (should be compatible with other implementations, let me know if not, was checked against the Javascript version)
* Changed API
    * Encrypt/Decrypt are now Encode/Decode
    * HashID is now constructed from HashIDData containing alphabet, salt and minimum length
