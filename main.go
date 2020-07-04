package main

import (
	"fmt"
	"github.com/kpfaulkner/gologmine/pkg/logmine"
	"os"
)

// hopefully this ends up as a working implementation of LogMine
// See https://www.cs.unm.edu/~mueen/Papers/LogMine.pdf and https://www.youtube.com/watch?v=URQnTPTxBbA for details.
func testTokenizer() {

	tokenizer := logmine.NewTokenizer()

	res, _ := tokenizer.Tokenize("2017/02/04 09:01:00 login 127.0.0.1 user=bear12")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:02:00 DB Connect 127.0.0.1 user=bear12")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:03:00 DB Disconnect 127.0.0.1 user=bear12")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:04:00 logout 127.0.0.1 user=bear12")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:05:00 login 127.0.0.1 user=bear34")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:06:00 DB Connect 127.0.0.1 user=bear34")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:07:00 DB Disconnect 127.0.0.1 user=bear34")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:08:00 logout 127.0.0.1 user=bear34")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:05:00 login 127.0.0.1 user=bear#1")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:06:00 DB Connect 127.0.0.1 user=bear#1")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:07:00 DB Disconnect 127.0.0.1 user=bear#1")
	fmt.Printf("res is %v\n", res)

	res, _ = tokenizer.Tokenize("2017/02/04 09:08:00 logout 127.0.0.1 user=bear#1")
	fmt.Printf("res is %v\n", res)
}

func main() {
	fmt.Printf("so it begins...\n")

	lm := logmine.NewLogMine()
  f, _ := os.Open("test.log")

  lm.ProcessLogsFromReader(f)

}
