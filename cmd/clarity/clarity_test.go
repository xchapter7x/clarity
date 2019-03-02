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
	s.Step(`^attribute "([^"]*)" equals (\d+)$`, match.AttributeEqualsInt)
	s.Step(`^attribute "([^"]*)" does not equal (\d+)$`, match.AttributeDoesNotEqualInt)
	s.Step(`^attribute "([^"]*)" equals "([^"]*)"$`, match.AttributeEquals)
	s.Step(`^attribute "([^"]*)" does not equal "([^"]*)"$`, match.AttributeDoesNotEqual)
	s.Step(`^attribute "([^"]*)" exists$`, match.AttributeExists)
	s.Step(`^it occurs at least (\d+) times$`, match.ItOccursAtLeastTimes)
	s.Step(`^it occurs at most (\d+) times$`, match.ItOccursAtMostTimes)
	s.Step(`^it occurs exactly (\d+) times$`, match.ItOccursExactlyTimes)
	s.Step(`^attribute "([^"]*)" matches regex "([^"]*)"$`, match.AttributeRegex)
	s.Step(`^a "([^"]*)" of type "([^"]*)" named "([^"]*)"$`, match.AOfTypeNamed)
	s.Step(`^attribute "([^"]*)" is greater than (\d+)$`, match.AttributeGreaterThan)
	s.Step(`^attribute "([^"]*)" is less than (\d+)$`, match.AttributeLessThan)
	s.BeforeScenario(func(interface{}) {
		match = matchers.NewMatch()
	})
}
