## Summary

Add user authentication to the system to allow users to securely log in and access protected resources.

## Motivation

The current system has no authentication mechanism. All routes are publicly accessible, which poses a security risk.

## Scope

- Login endpoint accepting email and password
- Session token generation on successful login
- Protected route middleware
