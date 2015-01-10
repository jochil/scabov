package export

import (
	"encoding/xml"
	"github.com/jochil/analyzer/classifier"
)

type xmlDev struct {
	XMLName        xml.Name `xml:"developer"`
	Id             string   `xml:"id,attr"`
	LinesAdded     int      `xml:"lines>added"`
	LinesRemoved   int      `xml:"lines>removed"`
	FilesAdded     int      `xml:"files>added"`
	FilesRemoved   int      `xml:"files>removed"`
	FilesChanged   int      `xml:"files>changed"`
	CylcoIncreased int      `xml:"cyclo>increased"`
	CylcoDecreased int      `xml:"cyclo>decreased"`
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
				Id:             id,
				LinesAdded:     int(devData["lines_added"]),
				LinesRemoved:   int(devData["lines_removes"]),
				FilesAdded:     int(devData["files_added"]),
				FilesRemoved:   int(devData["files_removed"]),
				FilesChanged:   int(devData["files_changed"]),
				CylcoIncreased: int(devData["cyclo_increased"]),
				CylcoDecreased: int(devData["cyclo_decreased"]),
			}
			xmlGroup.Devs = append(xmlGroup.Devs, dev)
		}
		xmlClassification.Groups = append(xmlClassification.Groups, xmlGroup)
	}

	root.Classification = xmlClassification
}
