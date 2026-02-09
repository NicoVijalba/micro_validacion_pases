# Security Policy

## Reporting a Vulnerability
Email: security@example.com
Include reproduction details, impact, affected versions.

## Supported Versions
- `1.x`: supported

## Secure Development Rules
- No secrets in repository.
- Rotate credentials and JWT signing keys periodically.
- Least privilege DB user (`INSERT` only for app account when feasible).
- Dependabot enabled and `govulncheck` required in CI.
