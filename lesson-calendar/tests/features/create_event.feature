Feature: Create event
  In order to use calendar
  As a user
  I want create event

  Scenario: create event
    When I create event
    Then There are no create errors
    And Event uuid should not be empty
  
  Scenario: create event on busy time
    When I create event on busy time
    Then I receive date already busy error
