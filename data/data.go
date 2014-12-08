package data

import (
	"fmt"
	"github.com/jochil/analyzer"
	"github.com/jochil/vcs"
)

type DevData struct {
	Dev           string //TODO replace with pointer to vcs.Dev?
	Commits       int
	FileDiff      analyzer.FileSumDiff
	LineDiff      analyzer.LineSumDiff
	LanguageUsage analyzer.LanguageUsage
}

func NewDevData(dev *vcs.Developer) *DevData {
	devData := &DevData{
		Dev:           dev.Id,
		Commits:       len(dev.Commits),
		FileDiff:      analyzer.CalcFileDiff(dev),
		LineDiff:      analyzer.CalcLineDiff(dev),
		LanguageUsage: analyzer.CalcLanguageUsage(dev),
	}
	return devData
}

func (d *DevData) String() string {
	return fmt.Sprintf("%s\n\tCommits: %d \n\tLanguage Usage: %f\n\tFiles: %s \n\tLines: %s",
		d.Dev, d.Commits, d.LanguageUsage.Value(), d.FileDiff, d.LineDiff)
}
