package main

import (
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/xchapter7x/clarity/pkg/matchers"
)

func FeatureContext(s *godog.Suite) {
	match := matchers.NewMatch()
	s.BeforeScenario(func(*messages.Pickle) {
		match = matchers.NewMatch()
	})
	s.Step(`^Pivotal OpsManager .ptc$`, markPending)

	s.Step(`^Terraform$`, match.Terraform)
	s.Step(`^a "([^"]*)" of type "([^"]*)"$`, match.AOfType)
	s.Step(`^a "([^"]*)" of type "([^"]*)" named "([^"]*)"$`, match.AOfTypeNamed)

	s.Step(`^"([^"]*)"$`, noopComment)
	s.Step(`^pending "([^"]*)"$`, markPending)
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

		s.Step(`^`+phrasePrefix+` "([^"]*)" always equals (\d+)$`, match.AlwaysAttributeEqualsInt)
		s.Step(`^`+phrasePrefix+` "([^"]*)" never equals (\d+)$`, match.AlwaysAttributeDoesNotEqualInt)
		s.Step(`^`+phrasePrefix+` "([^"]*)" always equals "([^"]*)"$`, match.AlwaysAttributeEquals)
		s.Step(`^`+phrasePrefix+` "([^"]*)" never equals "([^"]*)"$`, match.AlwaysAttributeDoesNotEqual)
		s.Step(`^`+phrasePrefix+` "([^"]*)" always matches regex "([^"]*)"$`, match.AlwaysAttributeRegex)
		s.Step(`^`+phrasePrefix+` "([^"]*)" is always greater than (\d+)$`, match.AlwaysAttributeGreaterThan)
		s.Step(`^`+phrasePrefix+` "([^"]*)" is always less than (\d+)$`, match.AlwaysAttributeLessThan)
	}

	s.Step(`^attribute "([^"]*)" exists$`, match.AttributeExists)
	s.Step(`^it occurs at least (\d+) times$`, match.ItOccursAtLeastTimes)
	s.Step(`^it occurs at most (\d+) times$`, match.ItOccursAtMostTimes)
	s.Step(`^it occurs exactly (\d+) times$`, match.ItOccursExactlyTimes)
}

func noopComment(arg1 string) error {
	return nil
}

func markPending(arg1 string) error {
	return godog.ErrPending
}
