# Secrets, Rotation and Least Privilege

## Storage
- Use GitHub Environments Secrets for CI/CD (`DEPLOY_HOST`, `DEPLOY_USER`, `DEPLOY_SSH_KEY`).
- Use Vault/secret manager for runtime secrets (`DB_DSN`, signing material if HS256).
- Never store secrets in git, Docker image, or logs.

## Rotation
- Rotate DB passwords every 90 days.
- Rotate JWT signing keys with overlap window and JWKS publication.
- Revoke leaked deploy keys immediately and regenerate.

## Least privilege
- DB user for app with minimal grants: `INSERT, SELECT` on target schema if possible.
- Separate CI identity from runtime identity.
- Restrict production environment approvals and branch protections.
