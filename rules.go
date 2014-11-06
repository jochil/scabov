package analyzer

import (
	"strings"
	"path/filepath"
)

type RuleSet interface {
	ValidExtension(ext string) bool
}

type PHPRuleSet struct {
}

func (r *PHPRuleSet) ValidExtension(path string) bool {
	ext := filepath.Ext(path)
	ext = strings.Trim(ext, ".")
	return (ext == "php")
}
