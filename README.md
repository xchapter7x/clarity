# Clarity
[![CircleCI](https://circleci.com/gh/xchapter7x/clarity.svg?style=svg)](https://circleci.com/gh/xchapter7x/clarity)

## A declaritive test framewark for Terraform
- reason: B/c unit testing terraform needs to be a thing

## Info:
- its gherkin bdd inspired
- provides its own matchers and hcl parser
- must be run from the directory where your terraform files live

## Download Binaries
[HERE](https://github.com/xchapter7x/clarity/releases/latest)

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

## Run the Tests
```bash
$ make unit
```

### Build the binary
```bash
$ make build
```

### Feature setup
```bash
$ nvim terraform/modules/control_plane.feature

Feature: We should have a LB for our control plane and its components and as
  such we should configure the proper security groups and listeners

  Scenario: we are using a single LB to route to all control plane components
    Given Terraform
    And a "aws_lb" of type "resource"
    Then attribute "load_balancer_type" equals "network"

  Scenario Outline: Every component of the control plane which needs a LB
    should be properly configured to have one

    Given Terraform
    And a "aws_security_group" of type "resource"
    And "our component is <component>"
    When attribute "ingress" exists
    Then attribute "from_port" equals "<port>"
    And attribute "to_port" equals "<port>"

    Given a "aws_lb_listener" of type "resource"
    And "our component is <component>"
    Then attribute "port" equals "<port>"

    Given a "aws_lb_target_group" of type "resource"
    And "our component is <component>"
    Then attribute "port" equals "<port>"

    Examples:
    | port | component |
    | 443  | ATC       |
    | 80   | ATC       |
    | 8443 | UAA       |
    | 2222 | TSA       |
    | 8844 | CredHub   |
```

### Sample usage
```
-> % clarity control_plane.feature
Feature: We should have a LB for our control plane and its components and as
  such we should configure the proper security groups and listeners

  Scenario: we are using a single LB to route to all control plane components # control_plane.feature:4
    Given Terraform                                                           # clarity_test.go:15 -> *Match
    And a "aws_lb" of type "resource"                                         # clarity_test.go:16 -> *Match
    Then attribute "load_balancer_type" equals "network"                      # clarity_test.go:17 -> *Match

  Scenario Outline: Every component of the control plane which needs a LB # control_plane.feature:9
    Given Terraform                                                       # clarity_test.go:15 -> *Match
    And a "aws_security_group" of type "resource"                         # clarity_test.go:16 -> *Match
    And "our component is <component>"                                    # clarity_test.go:9 -> noopComment
    When attribute "ingress" exists                                       # clarity_test.go:19 -> *Match
    Then attribute "from_port" equals "<port>"                            # clarity_test.go:17 -> *Match
    And attribute "to_port" equals "<port>"                               # clarity_test.go:17 -> *Match
    Given a "aws_lb_listener" of type "resource"                          # clarity_test.go:16 -> *Match
    And "our component is <component>"                                    # clarity_test.go:9 -> noopComment
    Then attribute "port" equals "<port>"                                 # clarity_test.go:17 -> *Match
    Given a "aws_lb_target_group" of type "resource"                      # clarity_test.go:16 -> *Match
    And "our component is <component>"                                    # clarity_test.go:9 -> noopComment
    Then attribute "port" equals "<port>"                                 # clarity_test.go:17 -> *Match

    Examples:
      | port | component |
      | 443  | ATC       |
        no matches found for attribute from_port
      | 80   | ATC       |
        no matches found for attribute from_port
      | 8443 | UAA       |
        no matches found for attribute from_port
      | 2222 | TSA       |
        no matches found for attribute from_port
      | 8844 | CredHub   |
        no matches found for attribute from_port

--- Failed steps:

  Scenario Outline: Every component of the control plane which needs a LB # control_plane.feature:9
    Then attribute "from_port" equals "443" # control_plane.feature:16
      Error: no matches found for attribute from_port

  Scenario Outline: Every component of the control plane which needs a LB # control_plane.feature:9
    Then attribute "from_port" equals "80" # control_plane.feature:16
      Error: no matches found for attribute from_port

  Scenario Outline: Every component of the control plane which needs a LB # control_plane.feature:9
    Then attribute "from_port" equals "8443" # control_plane.feature:16
      Error: no matches found for attribute from_port

  Scenario Outline: Every component of the control plane which needs a LB # control_plane.feature:9
    Then attribute "from_port" equals "2222" # control_plane.feature:16
      Error: no matches found for attribute from_port

  Scenario Outline: Every component of the control plane which needs a LB # control_plane.feature:9
    Then attribute "from_port" equals "8844" # control_plane.feature:16
      Error: no matches found for attribute from_port


6 scenarios (1 passed, 5 failed)
63 steps (23 passed, 5 failed, 35 skipped)
14.033302ms
```

## Contributions:
- all issues and PRs welcome.
