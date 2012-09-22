go-hashids
==========

Go (golang) v1 implementation of http://www.hashids.org
under MIT License (same as the original implementations)

Original implementations by ivanakimov

### Setup
<pre>go get github.com/speps/go-hashids</pre>

### Example
```go
package main

import "fmt"
import "github.com/speps/go-hashids"

func main() {
	h := hashids.New()
	h.Salt = "this is my salt"
	h.MinLength = 30
	e := h.Encrypt([]int{45, 434, 1313, 99})
	fmt.Println(e)
	d := h.Decrypt(e)
	fmt.Println(d)
}
```