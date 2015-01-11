package analyzer

import (
	"github.com/jochil/scabov/analyzer/php"
	"github.com/jochil/scabov/vcs"
	"github.com/stephens2424/php/token"
)

func CalcLanguageUsage(dev *vcs.Developer) LanguageUsage {

	//TODO add the added lines to calculation

	parser := NewParser()
	langUsage := NewLanguageUsage()

	for _, commit := range dev.Commits {
		for _, file := range commit.AddedFiles {

			parser.UpdateLanguageUsage(langUsage, file)

		}
	}

	return langUsage
}

type LanguageUsage interface {
	NumUsedElements() uint
	NumTotalElements() uint
	Value() float64
}

type PHPLanguageUsage struct {
	usedTokens            map[token.Token]uint
	usedInternalFunctions map[string]uint
}

func NewLanguageUsage() LanguageUsage {

	if Filter.Lang() == vcs.PHP {
		return &PHPLanguageUsage{
			usedTokens:            make(map[token.Token]uint),
			usedInternalFunctions: make(map[string]uint),
		}
	}
	return nil
}

func (langUsage *PHPLanguageUsage) AddItem(item token.Item) {

	//TODO this an invalid approach, as identifiers could also be variable names
	if item.Typ == token.Identifier {
		if php.IsInternalFunction(item.Val) {
			langUsage.usedInternalFunctions[item.Val]++
		}
	} else {
		langUsage.usedTokens[item.Typ]++
	}
}

func (langUsage *PHPLanguageUsage) NumUsedElements() uint {
	return uint(len(langUsage.usedTokens) + len(langUsage.usedInternalFunctions))

}

func (langUsage *PHPLanguageUsage) NumTotalElements() uint {
	return uint(len(token.TokenList)) + php.NumInternalFunctions()
}

func (langUsage *PHPLanguageUsage) Value() float64 {
	return float64(langUsage.NumUsedElements()) * 100 / float64(langUsage.NumTotalElements())
}
