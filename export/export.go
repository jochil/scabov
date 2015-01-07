package export

import (
	"encoding/xml"
	"io"
	"log"
	"os"
)

type xmlRoot struct {
	XMLName        xml.Name          `xml:"result"`
	Classification xmlClassification `xml:"classification"`
}

var root xmlRoot = xmlRoot{}

func SaveFile(file *os.File) {

	//save file
	xmlWriter := io.Writer(file)
	enc := xml.NewEncoder(xmlWriter)

	enc.Indent("  ", "    ")
	if err := enc.Encode(root); err != nil {
		log.Fatalf("Error while create classification xml output: %v\n", err)
	}
}
