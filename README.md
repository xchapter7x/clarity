# Clarity
[![CircleCI](https://circleci.com/gh/xchapter7x/clarity.svg?style=svg)](https://circleci.com/gh/xchapter7x/clarity)

## A declaritive test framework for Terraform
- reason: B/c unit testing terraform needs to be a thing

## Info:
- its a single binary, thats it. no need to write steps or go code. just write a declarative feature file.
- good for writing fitness function tests on your HCL
- its gherkin bdd inspired
- provides its own matchers and hcl parser
- must be run from the directory where your terraform files live
- its a wrapper on [godog](https://github.com/DATA-DOG/godog), which provides built in terraform specific steps

## Install

```bash
$ export OS="osx" #(osx | unix)
$ curl -sL https://github.com/xchapter7x/clarity/releases/latest/download/clarity_${OS} -o /usr/local/bin/clarity && chmod +x /usr/local/bin/clarity
```

### Download Binaries
[HERE](https://github.com/xchapter7x/clarity/releases/latest)

### CLI Options

```
-> % clarity --help
Usage:
  godog [options] [<features>]

Builds a test package and runs given feature files.
Command should be run from the directory of tested package and contain buildable go source.

Arguments:
  features             Optional feature(s) to run. Can be:
                           - dir (features/)
                           - feature (*.feature)
                           - scenario at specific line (*.feature:10)
                       If no feature paths are listed, suite tries features path by default.

Options:
  -c, --concurrency=1  Run the test suite with concurrency level:
                           - = 1: supports all types of formats.
                           - >= 2: only supports progress. Note, that
                           your context needs to support parallel execution.
  -d, --definitions    Print all available step definitions.
  -f, --format=pretty  How to format tests output. Built-in formats:
                           - events: Produces JSON event stream, based on spec: 0.1.0.
                           - junit: Prints junit compatible xml to stdout
                           - pretty: Prints every feature with runtime statuses.
                           - progress: Prints a character per step.
                           - cucumber: Produces cucumber JSON format output.
  --no-colors          Disable ansi colors.
  --random[=SEED]      Randomly shuffle the scenario execution order.
                       Specify SEED to reproduce the shuffling from a previous run.
                           e.g. --random or --random=5738
  --stop-on-failure    Stop processing on first failed scenario.
  --strict             Fail suite when there are pending or undefined steps.
  -t, --tags=          Filter scenarios by tags. Expression can be:
                           - "@wip": run all scenarios with wip tag
                           - "~@wip": exclude all scenarios with wip tag
                           - "@wip && ~@new": run wip scenarios, but exclude new
                           - "@wip,@undone": run wip or undone scenarios
```

### Writting your terraform tests

Simply put a .feature file in the directory where the terraform you wish to test resides
```gherkin
$ ls
control_plane.feature dns.tf                lb.tf                 network.tf            outputs.tf            variables.tf

$ cat terraform/modules/control_plane.feature

Feature: We should have a LB for our control plane and its components and as
  such we should configure the proper security groups and listeners

  Scenario: we are using a single LB to route to all control plane components
    Given Terraform
    And a "aws_lb" of type "resource"
    Then attribute "load_balancer_type" equals "network"
    And it occurs exactly 1 times

  Scenario Outline: Every component of the control plane which needs a LB
    should be properly configured to have one

    Given Terraform
    And a "aws_security_group" of type "resource"
    And "our component is <component>"
    When attribute "ingress" exists
    Then attribute "ingress" matches regex "from_port.*<port>"
    And attribute "ingress" matches regex "to_port.*<port>"

    Given Terraform
    And a "aws_lb_listener" of type "resource"
    And "our component is <component>"
    Then attribute "port" equals <port>

    Given Terraform
    And  a "aws_lb_target_group" of type "resource"
    And "our component is <component>"
    Then attribute "port" equals <port>

    Examples:
    | port | component |
    | 443  | ATC       |
    | 80   | ATC       |
    | 8443 | UAA       |
    | 2222 | TSA       |
    | 8844 | CredHub   |
```

### Running your terraform tests

Use the clarity cli to run any feature files you wish
```bash
-> % clarity control_plane.feature
Feature: We should have a LB for our control plane and its components and as
  such we should configure the proper security groups and listeners

  Scenario: we are using a single LB to route to all control plane components # control_plane.feature:4
    Given Terraform                                                           # clarity_test.go:15 -> *Match
    And a "aws_lb" of type "resource"                                         # clarity_test.go:16 -> *Match
    Then attribute "load_balancer_type" equals "network"                      # clarity_test.go:19 -> *Match
    And it occurs exactly 1 times                                             # clarity_test.go:24 -> *Match

  Scenario Outline: Every component of the control plane which needs a LB # control_plane.feature:10
    Given Terraform                                                       # clarity_test.go:15 -> *Match
    And a "aws_security_group" of type "resource"                         # clarity_test.go:16 -> *Match
    And "our component is <component>"                                    # clarity_test.go:9 -> noopComment
    When attribute "ingress" exists                                       # clarity_test.go:21 -> *Match
    Then attribute "ingress" matches regex "from_port.*<port>"            # clarity_test.go:25 -> *Match
    And attribute "ingress" matches regex "to_port.*<port>"               # clarity_test.go:25 -> *Match
    Given Terraform                                                       # clarity_test.go:15 -> *Match
    And a "aws_lb_listener" of type "resource"                            # clarity_test.go:16 -> *Match
    And "our component is <component>"                                    # clarity_test.go:9 -> noopComment
    Then attribute "port" equals <port>                                   # clarity_test.go:17 -> *Match
    Given Terraform                                                       # clarity_test.go:15 -> *Match
    And a "aws_lb_target_group" of type "resource"                        # clarity_test.go:16 -> *Match
    And "our component is <component>"                                    # clarity_test.go:9 -> noopComment
    Then attribute "port" equals <port>                                   # clarity_test.go:17 -> *Match

    Examples:
      | port | component |
      | 443  | ATC       |
      | 80   | ATC       |
      | 8443 | UAA       |
      | 2222 | TSA       |
      | 8844 | CredHub   |

6 scenarios (6 passed)
74 steps (74 passed)
35.093777ms
```

### gherkin step matchers available
	| '([^"]*)'                                       		| noop to insert context into behavior def                     |
	| 'Terraform'                                     		| parses the terraform from your local dir                     |
	| 'pending "([^"]*)"$'                            		| mark pending and skip all following steps in the scenario    |
	| 'a "([^"]*)" of type "([^"]*)"'                 		| matches on types such as resource,data and the resource name |
	| 'a "([^"]*)" of type "([^"]*)" named "([^"]*)"' 		| matches on types, resource names and instance names          |
	| 'the value of "([^"]*)" always equals (\d+)'          	| matches on the value given and the value of the attribute    |
	| 'the value of "([^"]*)" never equals (\d+)'      		| inverse match on attr value and given value                  |
	| 'the value of "([^"]*)" always equals "([^"]*)"'      	| matches on the value given and the value of the attribute    |
	| 'the value of "([^"]*)" never equals "([^"]*)"'  		| inverse match on attr value and given value                  |
	| 'the value of "([^"]*)" always matches regex "([^"]*)"'   	| matches the attributes value on the given regex              |
	| 'the value of "([^"]*)" is always greater than (\d+)'     	| matches on gt against the given value and attr value         |
	| 'the value of "([^"]*)" is always less than (\d+)'        	| matches on lt against the given value and attr value         |
	| 'attribute "([^"]*)" always equals (\d+)'              	| matches on the value given and the value of the attribute    |
	| 'attribute "([^"]*)" never equals (\d+)'      		| inverse match on attr value and given value                  |
	| 'attribute "([^"]*)" always equals "([^"]*)"'          	| matches on the value given and the value of the attribute    |
	| 'attribute "([^"]*)" never equals "([^"]*)"'  		| inverse match on attr value and given value                  |
	| 'attribute "([^"]*)" always matches regex "([^"]*)"'   	| matches the attributes value on the given regex              |
	| 'attribute "([^"]*)" is always greater than (\d+)'     	| matches on gt against the given value and attr value         |
	| 'attribute "([^"]*)" is always less than (\d+)'        	| matches on lt against the given value and attr value         |
	| 'the value of "([^"]*)" equals (\d+)'              		| matches on the value given and the value of the attribute    |
	| 'the value of "([^"]*)" does not equal (\d+)'      		| inverse match on attr value and given value                  |
	| 'the value of "([^"]*)" equals "([^"]*)"'          		| matches on the value given and the value of the attribute    |
	| 'the value of "([^"]*)" does not equal "([^"]*)"'  		| inverse match on attr value and given value                  |
	| 'the value of "([^"]*)" matches regex "([^"]*)"'   		| matches the attributes value on the given regex              |
	| 'the value of "([^"]*)" is greater than (\d+)'     		| matches on gt against the given value and attr value         |
	| 'the value of "([^"]*)" is less than (\d+)'        		| matches on lt against the given value and attr value         |
	| 'attribute "([^"]*)" equals (\d+)'              		| matches on the value given and the value of the attribute    |
	| 'attribute "([^"]*)" does not equal (\d+)'      		| inverse match on attr value and given value                  |
	| 'attribute "([^"]*)" equals "([^"]*)"'          		| matches on the value given and the value of the attribute    |
	| 'attribute "([^"]*)" does not equal "([^"]*)"'  		| inverse match on attr value and given value                  |
	| 'attribute "([^"]*)" matches regex "([^"]*)"'   		| matches the attributes value on the given regex              |
	| 'attribute "([^"]*)" is greater than (\d+)'     		| matches on gt against the given value and attr value         |
	| 'attribute "([^"]*)" is less than (\d+)'        		| matches on lt against the given value and attr value         |
	| 'attribute "([^"]*)" exists'                    		| if the given attribute exists in the matching objects        |
	| 'it occurs at least (\d+) times'                		| if the match set contains at least the given number          |
	| 'it occurs at most (\d+) times'                 		| if the match set contains at most the given number           |
	| 'it occurs exactly (\d+) times'                 		| if the match set continas exactly the given number           |

## Development & Contributions:
- all issues and PRs welcome.


### Run the Tests
```bash
$ make unit
```

### Build the binary
```bash
$ make build
```




## For more complete testing support using OPA policies check out:
https://github.com/instrumenta/conftest
