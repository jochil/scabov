package export

import (
	"encoding/xml"
	"fmt"
)

type xmlMetrics struct {
	XMLName                 xml.Name `xml:"metrics"`
	Stability               string   `xml:"stability"`
	StyleHomogeneity        string   `xml:"homogeneity>style"`
	ContributionHomogeneity string   `xml:"homogeneity>contribution"`
}

func SaveMetricsResult(stability float64, styleHomogeneity float64, contributionHomogeneity float64) {

	xmlMetrics := xmlMetrics{
		Stability:               fmt.Sprintf("%.4f", stability),
		StyleHomogeneity:        fmt.Sprintf("%.4f", styleHomogeneity),
		ContributionHomogeneity: fmt.Sprintf("%.4f", contributionHomogeneity),
	}

	root.Metrics = xmlMetrics
}
