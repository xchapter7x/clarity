# Clarity

## A declaritive test framewark for Terraform
- reason: B/c unit testing terraform needs to be a thing

## Info:
- its gherkin bdd inspired
- provides its own matchers and hcl parser
- must be run from the directory where your terraform files live

### gherkin steps
	| '([^"]*)'                                               | noop to insert context into behavior def                     |
	| 'Terraform'                                             | parses the terraform from your local dir                     |
	| 'a "([^"]*)" of type "([^"]*)"'                         | matches on types such as resource,data and the resource name |
	| 'a "([^"]*)" of type "([^"]*)" named "([^"]*)"'         | matches on types, resource names and instance names          |
	| 'attribute "([^"]*)" equals "([^"]*)"'                  | matches on the value given and the value of the attribute    |
	| 'attribute "([^"]*)" does not equal "([^"]*)"'          | inverse match on attr value and given value                  |
	| 'attribute "([^"]*)" exists'                            | if the given attribute exists in the matching objects        |
	| 'it occurs at least (\d+) times'                        | if the match set contains at least the given number          |
	| 'it occurs at most (\d+) times'                         | if the match set contains at most the given number           |
	| 'it occurs exactly (\d+) times'                         | if the match set continas exactly the given number           |
	| 'attribute "([^"]*)" matches regex "([^"]*)"'           | matches the attributes value on the given regex              |
	| 'attribute "([^"]*)" is greater than (\d+)'             | matches on gt against the given value and attr value         |
	| 'attribute "([^"]*)" is less than (\d+)'                | matches on lt against the given value and attr value         |
	| 'attribute "([^"]*)" is contained in object "([^"]*)"'  | the value in attribute is a subset of the given object       |
	| 'attribute "([^"]*)" contains object "([^"]*)"'         | the value in attribute is a superset of the given object     |
	| 'attribute "([^"]*)" is equivalent to object "([^"]*)"' | the value in attribute matches the given objet               |

## Contributions:
- all issues and PRs welcome.
