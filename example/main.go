/**
 * Author:    hashcode55 (Mehul Ahuja)
 * Created:   11.05.2017
 **/

package main

// Remove local imports
import (
	"flag"
	"github.com/hashcode55/gopython"
)

func main() {
	boolPtr := flag.Bool("log", false, "Set it to true to log the details.")
	flag.Parse()
	gopython.ParseEngine("hello = 2 + 3 * 6", *boolPtr)
}
