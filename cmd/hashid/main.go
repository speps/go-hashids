package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/speps/go-hashids/v2"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"usage:\n"+
				"\t%s [options] <intlist> [<intlist>...]\n"+
				"\t%s [options] -d <hashid> [<hashid>...]\n\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	var decode bool
	var params hashids.HashIDData
	var separator string
	flag.StringVar(&params.Salt, `salt`, "", `salt`)
	flag.StringVar(&params.Alphabet, `alphabet`, hashids.DefaultAlphabet, `minimum 16 characters`)
	flag.IntVar(&params.MinLength, `min`, 0, `minimum length (for encoding)`)
	flag.BoolVar(&decode, `d`, false, `decode (instead of encoding)`)
	flag.StringVar(&separator, `sep`, ",", `separator for integers`)
	flag.Parse()

	codec, err := hashids.NewWithData(&params)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	args := os.Args[len(os.Args)-flag.NArg():]
	if decode {
		for _, arg := range args {
			result, err := codec.DecodeInt64WithError(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", arg, err)
			} else {
				var str []byte
				for _, x := range result {
					if len(str) != 0 {
						str = append(str, separator...)
					}
					str = strconv.AppendInt(str, x, 10)
				}
				fmt.Printf("%s: %s\n", arg, str)
			}
		}
	} else {
	ARGS:
		for _, arg := range args {
			spl := strings.Split(arg, separator)
			ints := make([]int64, len(spl))
			var err error
			for i, s := range spl {
				ints[i], err = strconv.ParseInt(s, 0, 64)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s: %s\n", arg, err)
					continue ARGS
				}
			}
			result, err := codec.EncodeInt64(ints)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", arg, err)
			} else {
				fmt.Printf("%s: %v\n", arg, result)
			}
		}
	}
}
