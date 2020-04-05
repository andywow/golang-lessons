Feature: List date events
  In order to use calendar
  As a user
  I want list events for date

  Scenario: list events for date
    When I list events for date
    Then There are no list errors
    And Event list should not be empty

  Scenario: list events for week
    When I list events for week
    Then There are no list errors
    And Event list should not be empty

  Scenario: list events for month
    When I list events for month
    Then There are no list errors
    And Event list should not be empty
