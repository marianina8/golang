# GenBank Parser

This code parses the standard ‘Genbank' biological data format into a usable Go struct

  * Usage:
    import (
   	    "fmt"
	    "github.com/mmontagnino/gbparser"
    )

    func main() {
	fmt.Printf(gbparser.Parse("sequence.gb"))
    }
