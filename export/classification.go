package export

import (
	"encoding/xml"
	"fmt"
	"github.com/jochil/scabov/analyzer/classifier"
)

type xmlDev struct {
	XMLName xml.Name `xml:"developer"`
	Id      string   `xml:"id,attr"`
	Data    []byte   `xml:",innerxml"`
}

type xmlGroup struct {
	XMLName xml.Name `xml:"group"`
	Devs    []xmlDev `xml:"developers"`
}

type xmlClassification struct {
	XMLName xml.Name   `xml:"classification"`
	Id      string     `xml:"id,attr"`
	Groups  []xmlGroup `xml:"groups"`
}

func SaveClassificationResult(id string, groups []*classifier.Group, rawMatrix map[string]map[string]float64) {

	xmlClassification := xmlClassification{Id: id}

	//create xml structure
	for _, group := range groups {
		xmlGroup := xmlGroup{}

		for _, id := range group.Objects {
			devData := rawMatrix[id]
			dev := xmlDev{Id: id}

			if xmlClassification.Id == "style" {
				data := fmt.Sprintf("<cyclo><avg>%.4f</avg></cyclo><language><usage>%.4f</usage></language><function><size>%.4f</size></function>",
					devData["cyclo_avg"], devData["language_usage"], devData["function_size"])

				dev.Data = []byte(data)

			} else if xmlClassification.Id == "contribution" {

				fileData := fmt.Sprintf("<files><added>%.0f</added><removed>%.0f</removed><changed>%.0f</changed></files>",
					devData["files_added"], devData["files_removed"], devData["files_changed"])
				lineData := fmt.Sprintf("<lines><added>%.0f</added><removed>%.0f</removed></lines>",
					devData["lines_added"], devData["lines_removed"])
				cycloData := fmt.Sprintf("<cyclo><increased>%.0f</increased><decreased>%.0f</decreased></cyclo>",
					devData["cyclo_increased"], devData["cyclo_decreased"])

				data := fileData + lineData + cycloData
				dev.Data = []byte(data)
			}

			xmlGroup.Devs = append(xmlGroup.Devs, dev)
		}
		xmlClassification.Groups = append(xmlClassification.Groups, xmlGroup)
	}

	root.Classification = append(root.Classification, xmlClassification)
}
