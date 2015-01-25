package export

import (
	"encoding/xml"
	"github.com/jochil/scabov/analyzer"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
)

type xmlRoot struct {
	XMLName        xml.Name            `xml:"result"`
	Metrics        xmlMetrics          `xml:"metrics"`
	Files          []xmlFile           `xml:"files>file"`
	Classification []xmlClassification `xml:"classifications>classification"`
}

var root xmlRoot = xmlRoot{}

func SaveFile(file *os.File) {

	//save file
	xmlWriter := io.Writer(file)
	enc := xml.NewEncoder(xmlWriter)

	enc.Indent("  ", "    ")
	if err := enc.Encode(root); err != nil {
		log.Fatalf("Error while creating xml output: %v\n", err)
	}
}

func DumpCfg(function analyzer.Function, workspace string) {

	fileBase := path.Join(workspace, function.Name)
	function.CFG.ToDOTFile(fileBase + ".dot")
	cmd := exec.Command("dot", "-Tpng", fileBase+".dot", "-o", fileBase+".png")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
