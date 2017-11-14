/**
 * Author:    hashcode55 (Mehul Ahuja)
 * Created:   11.05.2017
 **/

package main

// Remove local imports
import (
	"flag"
	"github.com/HashCode55/GPython"
)

func main() {
	boolPtr := flag.Bool("log", true, "Set it to true to log the details.")
	flag.Parse()
	gpython.ParseEngine("hello = ( 3 * 6 ) - 5 / 2", *boolPtr)
}
