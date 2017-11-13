/**
 * Author:    hashcode55 (Mehul Ahuja)
 * Created:   11.05.2017
 **/

package main

// Remove local imports
import (
	"flag"
	"github.com/puneets2811/Gpython"
)

func main() {
	boolPtr := flag.Bool("log", true, "Set it to true to log the details.")
	flag.Parse()
	gpython.ParseEngine("a = ab", *boolPtr)
}
