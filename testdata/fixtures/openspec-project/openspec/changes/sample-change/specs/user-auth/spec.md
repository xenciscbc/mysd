## Requirement: User Authentication

The system MUST validate user credentials before granting access.
The system SHOULD log all authentication attempts for auditing purposes.
The system MUST reject requests with invalid or expired tokens.

### Scenario: Successful Login

WHEN a user submits valid credentials
THEN the system returns a session token
AND the token is valid for 24 hours

### Scenario: Failed Login

WHEN a user submits invalid credentials
THEN the system returns an error response
AND the system SHOULD increment the failed attempt counter
AND the system MAY lock the account after 5 failed attempts
