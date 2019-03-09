package main

import (
	"github.com/DATA-DOG/godog"
	"github.com/xchapter7x/clarity/pkg/matchers"
)

func noopComment(arg1 string) error {
	return nil
}

func markPending(arg1 string) error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	match := matchers.NewMatch()
	s.Step(`^"([^"]*)"$`, noopComment)
	s.Step(`^pending "([^"]*)"$`, markPending)
	s.Step(`^Terraform$`, match.Terraform)
	s.Step(`^a "([^"]*)" of type "([^"]*)"$`, match.AOfType)
	s.Step(`^a "([^"]*)" of type "([^"]*)" named "([^"]*)"$`, match.AOfTypeNamed)

	for _, phrasePrefix := range []string{
		"the value of",
		"attribute",
	} {
		s.Step(`^`+phrasePrefix+` "([^"]*)" equals (\d+)$`, match.AttributeEqualsInt)
		s.Step(`^`+phrasePrefix+` "([^"]*)" does not equal (\d+)$`, match.AttributeDoesNotEqualInt)
		s.Step(`^`+phrasePrefix+` "([^"]*)" equals "([^"]*)"$`, match.AttributeEquals)
		s.Step(`^`+phrasePrefix+` "([^"]*)" does not equal "([^"]*)"$`, match.AttributeDoesNotEqual)
		s.Step(`^`+phrasePrefix+` "([^"]*)" matches regex "([^"]*)"$`, match.AttributeRegex)
		s.Step(`^`+phrasePrefix+` "([^"]*)" is greater than (\d+)$`, match.AttributeGreaterThan)
		s.Step(`^`+phrasePrefix+` "([^"]*)" is less than (\d+)$`, match.AttributeLessThan)
	}

	s.Step(`^attribute "([^"]*)" exists$`, match.AttributeExists)
	s.Step(`^it occurs at least (\d+) times$`, match.ItOccursAtLeastTimes)
	s.Step(`^it occurs at most (\d+) times$`, match.ItOccursAtMostTimes)
	s.Step(`^it occurs exactly (\d+) times$`, match.ItOccursExactlyTimes)
	s.BeforeScenario(func(interface{}) {
		match = matchers.NewMatch()
	})
}
