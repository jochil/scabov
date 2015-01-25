package export

import (
	"encoding/xml"
	"fmt"
	"github.com/jochil/scabov/analyzer"
)

type xmlFile struct {
	XMLName   xml.Name      `xml:"file"`
	Path      []byte        `xml:",innerxml"`
	Functions []xmlFunction `xml:"functions>function"`
}

type xmlFunction struct {
	XMLName          xml.Name `xml:"function"`
	Name             []byte   `xml:",innerxml"`
	Stability        string   `xml:"stability"`
	SizeGrowth       string   `xml:"growth>size"`
	ComplexityGrowth string   `xml:"growth>complexity"`
}

func SaveFunctions(history map[string]analyzer.FileHistory) {

	//create xml structure
	for filename, fileHistory := range history {
		xmlFile := xmlFile{Path: []byte("<path><![CDATA[" + filename + "]]></path>")}

		for functionName, functionHistory := range fileHistory {

			sizeGrowth, complexityGrowth := functionHistory.Growth()

			xmlFunction := xmlFunction{
				Name:             []byte("<name><![CDATA[" + functionName + "]]></name>"),
				Stability:        fmt.Sprintf("%.4f", functionHistory.Stability()),
				SizeGrowth:       fmt.Sprintf("%.4f", sizeGrowth),
				ComplexityGrowth: fmt.Sprintf("%.4f", complexityGrowth),
			}

			xmlFile.Functions = append(xmlFile.Functions, xmlFunction)
		}

		root.Files = append(root.Files, xmlFile)
	}

}
