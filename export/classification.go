package export

import (
	"encoding/xml"
	"github.com/jochil/analyzer/classifier"
)

type xmlDev struct {
	XMLName      xml.Name `xml:"developer"`
	Id           string   `xml:"id,attr"`
	Commits      int      `xml:"commits"`
	LinesAdded   int      `xml:"lines>added"`
	LinesRemoved int      `xml:"lines>removed"`
}

type xmlGroup struct {
	XMLName xml.Name `xml:"group"`
	Devs    []xmlDev `xml:"developers"`
}

type xmlClassification struct {
	XMLName xml.Name   `xml:"classification"`
	Groups  []xmlGroup `xml:"groups"`
}

func SaveClassificationResult(groups []*classifier.Group, rawMatrix map[string]map[string]float64) {

	xmlClassification := xmlClassification{}

	//create xml structure
	for _, group := range groups {
		xmlGroup := xmlGroup{}

		for _, id := range group.Objects {
			devData := rawMatrix[id]
			dev := xmlDev{
				Id:           id,
				Commits:      int(devData["commits"]),
				LinesAdded:   int(devData["lines_added"]),
				LinesRemoved: int(devData["lines_removes"]),
			}
			xmlGroup.Devs = append(xmlGroup.Devs, dev)
		}
		xmlClassification.Groups = append(xmlClassification.Groups, xmlGroup)
	}

	root.Classification = xmlClassification
}
