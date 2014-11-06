package analyzer

import "log"

var ruleSet RuleSet

func Init(language string) {
	switch language{
	case "php":
		ruleSet = &PHPRuleSet{}
	}

	log.Println("startet analyzer with PHP rule set")
}
