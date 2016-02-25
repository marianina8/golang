/**************************************************
* Filename: main.go
*
* Input: Genbank file (extension: .gb)
*
* Outputs: GenBank struct in JSON format
*
* Example:
*	go run main.go /Users/mmontagnino/Documents/sequence.gb
*
* Author: Marian Montagnino
* Date Created: 2/24/16
* https://github.com/marianina8
*
**************************************************/

package gbparser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// GenBank
type GenBank struct {
	Origin     string
	Locus      Lcs
	Source     Src
	References []Reference
	Primaries  []Primary
	Features   []Feature
}

// Origin
var Origin string

// Locus
type Lcs struct {
	Name, SequenceLength, MoleculeType, GenBankDivision, ModDate string
}

// Source
type Src struct {
	Name, Organism string
}

// Reference
var References []Reference

type Reference struct {
	Index, Authors, Title, Journal, PubMed, Remark string
}

// Primary
var Primaries []Primary

type Primary struct {
	RefSeq, PrimaryIdentifier, Primary_Span, Comp string
}

// Feature
var Features []Feature

type Feature struct {
	Name, Location, Sequence string
	Qualifiers               map[string]string
}

// getReference loops through to get all references within genbank file
func getReference(buf int, lines []string) (Reference, int) {

	var Ref Reference

	Ref.Index = strings.TrimSpace(string(lines[buf])[12:14])
	for {
		if strings.Compare(strings.TrimSpace(string(lines[buf])[0:11]), "AUTHORS") == 0 {
			break
		}
		buf++
	}

	Ref.Authors = strings.TrimSpace(string(lines[buf])[12:])
	buf++
	for {
		if string(lines[buf][0:12]) == "            " {
			Ref.Authors += strings.TrimSpace(string(lines[buf])[12:]) + " "
		} else {
			break
		}
		buf++
	}

	for {
		if strings.Compare(strings.TrimSpace(string(lines[buf])[0:11]), "TITLE") == 0 {
			break
		}
		buf++
	}

	Ref.Title = strings.TrimSpace(string(lines[buf])[12:])
	buf++
	for {
		if string(lines[buf][0:12]) == "            " {
			Ref.Title += strings.TrimSpace(string(lines[buf])[12:]) + " "
		} else {
			break
		}
		buf++
	}

	for {
		if strings.Compare(strings.TrimSpace(string(lines[buf])[0:11]), "JOURNAL") == 0 {
			break
		}
		buf++
	}

	Ref.Journal = strings.TrimSpace(string(lines[buf])[12:])
	buf++
	for {
		if string(lines[buf][0:12]) == "            " {
			Ref.Journal += strings.TrimSpace(string(lines[buf])[12:]) + " "
		} else {
			break
		}
		buf++
	}

	for {
		if strings.Compare(strings.TrimSpace(string(lines[buf])[0:11]), "PUBMED") == 0 {
			break
		}
		buf++
	}

	Ref.PubMed = strings.TrimSpace(string(lines[buf])[12:])

	for {
		if strings.Compare(strings.TrimSpace(string(lines[buf])[0:11]), "REMARK") == 0 {
			break
		} else if strings.Compare(strings.TrimSpace(string(lines[buf])[0:11]), "REFERENCE") == 0 || string(lines[buf][0]) != " " {
			buf--
			return Ref, buf
		}
		buf++
	}

	Ref.Remark = strings.TrimSpace(string(lines[buf])[12:])
	buf++
	for {
		if string(lines[buf][0:12]) == "            " {
			Ref.Remark += strings.TrimSpace(string(lines[buf])[12:]) + " "
		} else {
			break
		}
		buf++
	}

	buf--
	return Ref, buf
}

func Parse(filename string) string {

	/*
	     // If missing filename, print usage
	   	if len(os.Args) < 2 {
	   		fmt.Println("Usage:")
	   		fmt.Println("    go run main.go {location/filename}")
	   		return ""
	   	}

	   	// Open file
	   	filename := os.Args[1]
	*/
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	bio := bufio.NewReader(f)

	var lines []string
	i := 0

	// Read all lines of the file into buffer
	for {
		line, err := bio.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		sline := strings.TrimRight(string(line), `\n`)
		lines = append(lines, sline)
		i++
	}
	// End read of file into buffer

	// Create GenBank struct
	var G GenBank

	// Populate Locus
	var L Lcs
	L.Name = strings.TrimSpace(string(lines[0])[12:25])
	L.SequenceLength = strings.TrimSpace(string(lines[0])[36:41])
	L.MoleculeType = strings.TrimSpace(string(lines[0])[41:43])
	L.GenBankDivision = strings.TrimSpace(string(lines[0])[47:51])
	L.ModDate = strings.TrimSpace(string(lines[0])[68:79])

	// Assign locus to GenBank struct
	G.Locus = L

	buf := 0

	// Search for SOURCE
	for {
		if strings.Compare(strings.TrimSpace(string(lines[buf])[0:11]), "SOURCE") == 0 {
			break
		}
		buf++
	}

	// Populate Source
	var S Src
	S.Name = strings.TrimSpace(string(lines[buf])[12:])
	buf++
	S.Organism = strings.TrimSpace(string(lines[buf])[12:]) + " "
	buf++
	for {
		if string(lines[buf][0]) == " " {
			S.Organism += strings.TrimSpace(string(lines[buf])[12:]) + " "
		} else {
			break
		}
		buf++
	}

	// Assign Source to GenBank struct
	G.Source = S

	// Get all References
	for {
		if strings.Compare(strings.TrimSpace(string(lines[buf])[0:11]), "REFERENCE") == 0 {
			Ref, newBuf := getReference(buf, lines)
			buf = newBuf
			// Append Reference to References array in GenBank struct
			G.References = append(G.References, Ref)
		}
		// Once search reaches Comment, break from loop
		if strings.Compare(strings.TrimSpace(string(lines[buf])[0:11]), "COMMENT") == 0 {
			break
		}
		buf++
	}

	buf--

	// Search for PRIMARY
	for {
		if len(lines[buf]) >= 11 {
			if strings.Compare(strings.TrimSpace(string(lines[buf])[0:11]), "PRIMARY") == 0 {
				break
			}
		}

		buf++
	}

	buf++

	// Search for Primaries
	for {
		if string(lines[buf][0:12]) == "            " {
			var P Primary
			P.RefSeq = strings.TrimSpace(string(lines[buf])[12:23])
			P.PrimaryIdentifier = strings.TrimSpace(string(lines[buf])[32:50])
			if len(lines[buf]) > 73 {
				P.Primary_Span = strings.TrimSpace(string(lines[buf])[51:66])
				P.Comp = strings.TrimSpace(string(lines[buf])[72:73])
			} else {
				P.Primary_Span = strings.TrimSpace(string(lines[buf])[51:])
			}
			// Append each Primary to Primaries array in GenBank Struct
			G.Primaries = append(G.Primaries, P)
		}
		if strings.Compare(strings.TrimSpace(string(lines[buf])[0:11]), "FEATURES") == 0 {
			break
		}
		buf++
	}

	// Search for Features
	for {
		if strings.Compare(strings.TrimSpace(string(lines[buf])[0:11]), "FEATURES") == 0 {
			var F Feature
			qualifiers := make(map[string]string)
			if strings.Compare(strings.TrimSpace(string(lines[buf])[0:11]), "FEATURES") == 0 {
				buf++
			}
			for {
				if strings.Compare(strings.TrimSpace(string(lines[buf])[7:8]), "") != 0 {

					// Creating a new feature
					F.Qualifiers = make(map[string]string)
					qualifiers = make(map[string]string)
					F.Name = strings.TrimSpace(string(lines[buf])[5:21])
					F.Location = strings.TrimSpace(string(lines[buf])[21:])
					buf++

					for {
						if strings.Compare(strings.TrimSpace(string(lines[buf])[21:22]), "/") == 0 { //parse qualifier
							// Found qualifier
							q := strings.TrimSpace(string(lines[buf])[22:])
							if strings.Contains(q, "=") {
								quarry := strings.Split(q, "=")
								// Add qualifier
								qualifiers[quarry[0]] = quarry[1]

								// Handle qualifiers that take up multiple lines
								for {
									if len(lines[buf+1]) > 22 {
										if strings.Compare(strings.TrimSpace(string(lines[buf+1])[21:22]), "/") != 0 && strings.Compare(strings.TrimSpace(string(lines[buf+1])[0:21]), "") == 0 {
											if strings.Compare(strings.TrimSpace(string(lines[buf+1])[0:7]), "ORIGIN") != 0 {
												if quarry[0] == "note" || quarry[0] == "experiment" {
													qualifiers[quarry[0]] += " " + strings.TrimSpace(string(lines[buf+1])[21:])
												} else {
													qualifiers[quarry[0]] += strings.TrimSpace(string(lines[buf+1])[21:])
												}
											}
										} else {
											break
										}

									}
									buf++
								}

							}
						} else {
							if qualifiers != nil {
								F.Qualifiers = qualifiers
								buf--
							}
							break
						}
						buf++
					}
					// Assign each feature to features array in GenBank struct.
					G.Features = append(G.Features, F)
					if strings.Compare(strings.TrimSpace(string(lines[buf])[0:7]), "ORIGIN") == 0 {
						break
					}

				}
				buf++
			}

		}
		if strings.Compare(strings.TrimSpace(string(lines[buf])[0:7]), "ORIGIN") == 0 {
			break
		}
		buf++
	}

	// Extract Origin by appending all lines and removing spaces and line information.
	for {
		if strings.Compare(strings.TrimSpace(string(lines[buf])[0:2]), "//") == 0 {
			break
		} else {
			Origin += strings.TrimSpace(string(lines[buf])[10:])
		}
		buf++
	}
	Origin = strings.Replace(Origin, " ", "", -1)

	// Assign Origin to GenBank struct
	G.Origin = Origin

	// Loop back through all features, grabbing the location information and setting the associated sequence.
	for i := range G.Features {
		if strings.Contains(G.Features[i].Location, "JOIN") {
			sublocation := strings.Replace(G.Features[i].Location, "JOIN(", "", -1)
			sublocation = strings.Replace(G.Features[i].Location, ")", "", -1)
			sublocations := strings.Split(sublocation, ",")
			Seq := ""
			for location := range sublocations {
				if strings.Contains(string(location), "..") {
					numbers := strings.Split(G.Features[i].Location, "..")
					start, err := strconv.Atoi(numbers[0])
					start--
					end, err := strconv.Atoi(numbers[1])
					if err != nil {
						fmt.Println(err)
					}
					Seq += Origin[start:end]
				} else {
					Seq += string(Origin[location] - 1)
				}
			}
			// Assign associated sequence to GenBank's feature if a joined value
			G.Features[i].Sequence = Seq
		} else if strings.Contains(G.Features[i].Location, "..") {
			numbers := strings.Split(G.Features[i].Location, "..")
			start, err := strconv.Atoi(numbers[0])
			start--
			end, err := strconv.Atoi(numbers[1])
			// Assign associated sequence to GenBank's feature if a range
			G.Features[i].Sequence = Origin[start:end]
			if err != nil {
				fmt.Println(err)
			}
		} else {
			k, err := strconv.Atoi(G.Features[i].Location)
			k--
			if err != nil {
				fmt.Println(err)
			}
			if k >= 0 {
				// Assign associated sequence to GenBank's feature if a single value
				G.Features[i].Sequence = string(Origin[k])
			}

		}
	}

	a := &G

	// convert GenBank struct into JSON (indented) format
	out, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(out)

}
