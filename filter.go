package vcs

import (
	"log"
	"path/filepath"
	"strings"
)

const (
	PHP = "php"
)

type LanguageFilter interface {
	ValidExtension(path string) bool
	Lang() string
}

func NewLanguageFilter(lang string) LanguageFilter {

	var filter LanguageFilter = PassThroughFilter{}

	switch strings.ToLower(lang) {
	case PHP:
		filter = PHPFilter{}
	}

	log.Printf("use %s filter", lang)

	return filter
}

//Filter that filters nothing ;)
type PassThroughFilter struct {
}

func (filter PassThroughFilter) ValidExtension(path string) bool {
	return true
}

func (filter PassThroughFilter) Lang() string {
	return ""
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

func (filter PHPFilter) Lang() string {
	return PHP
}
