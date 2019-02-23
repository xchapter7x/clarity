Feature: In order to test out syntax desired from
  the clarity declarative tests, we have this
  file as a step generator, so we can see how
  the language will flow before investing
  time in building it out.
  
  Scenario: Policy Structure
    Given Terraform
    And a "aws_security_group" of type "resource"
    Then attribute "attrelement" exists

  Scenario: We would like our resources to be able to talk to things on the internet
    Given Terraform
    And a "aws_security_group" of type "resource"
    Then attribute "blah" is greater than 1234

  Scenario: We would like our resources to be able to talk to things on the internet
    Given Terraform
    And "we need a fully functional something"
    And a "aws_security_group" of type "resource"
    And a "aws_security_group" of type "resource"
    And a "aws_security_group" of type "resource"
    When attribute "someelement" exists
    Then attribute "blah" is greater than 1234

    Given Terraform
    When attribute "someelement" exists
    And it occurs at least 2 times
    And attribute "someelement" exists 
    And it occurs at most 2 times
    And attribute "someelement" exists 
    And it occurs exactly 2 times
    Then attribute "blah" equals "1234"
    And attribute "blah" does not equal "1234"
    And attribute "blah" matches regex "bleh"
    And attribute "blah" matches regex "hi.*there"

  Scenario: We would like our resources to be able to talk to things on the internet
    Given Terraform
    And a "aws_security_group" of type "resource"
    And a "aws_security_group" of type "resource"
    And a "aws_security_group" of type "resource"
    And a "aws_security_group" of type "resource" named "control_plane_internal"
    Then attribute "blah" is less than 1234

