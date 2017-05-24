/**
 * Author:    hashcode55 (Mehul Ahuja)
 * Created:   11.05.2017
 **/

package main

// Remove local imports
import (
	".."
	"flag"
)

func main() {
	boolPtr := flag.Bool("log", false, "Set it to true to log the details.")
	flag.Parse()
	gython.ParseEngine("hello = hello", *boolPtr)
}
