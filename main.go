package main

import (
	"fmt"
	"github.com/kpfaulkner/gologmine/pkg/logmine"
)


// hopefully this ends up as a working implementation of LogMine
// See https://www.cs.unm.edu/~mueen/Papers/LogMine.pdf and https://www.youtube.com/watch?v=URQnTPTxBbA for details.
func main() {
	fmt.Printf("so it begins...\n")


  tokenizer := logmine.NewTokenizer()

  res,err := tokenizer.Tokenize("2017/02/04 09:01:00 login 127.0.0.1 user=bear12")
	//res,err := tokenizer.Tokenize("bear12")

  if err != nil {
  	fmt.Printf("error is %s\n", err.Error())
  	return
  }

  fmt.Printf("res is %v\n", res)

}
