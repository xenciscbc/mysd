## Architecture

The authentication system uses a stateless JWT-based approach.

## Key Decisions

- Use JWT tokens instead of server-side sessions for horizontal scalability
- Tokens expire after 24 hours and cannot be revoked (acceptable for v1)
- Passwords hashed with bcrypt (cost factor 12)
