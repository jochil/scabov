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
	CycloIncreased int      `xml:"cyclo>increased"`
	CycloDecreased int      `xml:"cyclo>decreased"`
	CycloAvg       int      `xml:"cyclo>avg"`
	CycloMax       int      `xml:"cyclo>max"`
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
			dev := xmlDev{
				Id:             id,
				LinesAdded:     int(devData["lines_added"]),
				LinesRemoved:   int(devData["lines_removes"]),
				FilesAdded:     int(devData["files_added"]),
				FilesRemoved:   int(devData["files_removed"]),
				FilesChanged:   int(devData["files_changed"]),
				CycloIncreased: int(devData["cyclo_increased"]),
				CycloDecreased: int(devData["cyclo_decreased"]),
				CycloAvg:       int(devData["cyclo_avg"]),
				CycloMax:       int(devData["cyclo_max"]),
			}
			xmlGroup.Devs = append(xmlGroup.Devs, dev)
		}
		xmlClassification.Groups = append(xmlClassification.Groups, xmlGroup)
	}

	root.Classification = append(root.Classification, xmlClassification)
}
