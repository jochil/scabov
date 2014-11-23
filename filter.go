package vcs

import (
	"strings"
	"path/filepath"
)

const (
	PHP = "php"
)

type LanguageFilter interface {
	ValidExtension(path string) bool
}

func NewLanguageFilter(lang string) LanguageFilter {

	var filter LanguageFilter = PassThroughFilter{}

	switch strings.ToLower(lang) {
	case PHP:
		filter = PHPFilter{}
	}

	return filter;
}

//Filter that filters nothing ;)
type PassThroughFilter struct {

}

func (filter PassThroughFilter) ValidExtension(path string) bool {
	return true
}

// PHP filter
type PHPFilter struct {

}

var phpExtensions = [...]string{
	"php",
}

func (filter PHPFilter) ValidExtension(path string) bool {
	ext := filepath.Ext(path)
	ext = strings.Trim(ext, ".")
	for _, phpExt := range phpExtensions {
		if phpExt == ext {
			return true
		}
	}
	return false
}
