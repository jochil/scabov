package export

import (
	"encoding/xml"
	"io"
	"log"
	"os"
	"github.com/jochil/scabov/analyzer"
	"path"
	"os/exec"
)

type xmlRoot struct {
	XMLName        xml.Name            `xml:"result"`
	Classification []xmlClassification `xml:"classifications>classification"`
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

func DumpCfg(function analyzer.Function, workspace string) {

	fileBase := path.Join(workspace, function.Name)
	function.CFG.ToDOTFile(fileBase + ".dot")
	cmd := exec.Command("dot", "-Tpng", fileBase+".dot", "-o", fileBase+".png")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
