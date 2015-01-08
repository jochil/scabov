package export

import (
	"fmt"
	"log"
	"sort"
)

func PrintMatrix(matrix map[string]map[string]float64) {

	columnSize := 25

	log.Println("---------------")

	//get first level keys
	keys := []string{}
	for key := range matrix {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	//get second level key
	submap := matrix[keys[0]]
	subkeys := []string{}
	for key := range submap {
		subkeys = append(subkeys, key)
	}
	sort.Strings(subkeys)

	//print header
	out := fillColumn("", columnSize)
	for _, key := range subkeys {
		out += fillColumn(key, columnSize)
	}
	log.Println(out)

	//print rows
	for _, key := range keys {

		submap = matrix[key]

		out = fillColumn(key, columnSize)

		for _, subkey := range subkeys {
			value := submap[subkey]
			out += fillColumn(fmt.Sprintf("%.3f", value), columnSize)
		}

		log.Println(out)
	}
	log.Println("---------------")
}

func fillColumn(content string, columnSize int) string {

	out := content
	offset := columnSize - len(content)
	for i := 0; i <= offset; i++ {
		out += " "
	}

	return out
}
