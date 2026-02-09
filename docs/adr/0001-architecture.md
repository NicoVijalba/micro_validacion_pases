# ADR 0001: Clean Architecture (Hexagonal)

## Status
Accepted - 2026-02-09

## Context
Se requiere un microservicio mantenible, testeable y con seguridad reforzada, desacoplado de framework y BD.

## Decision
Usar capas:
- Dominio (`internal/domain`): entidades y puertos.
- Casos de uso (`internal/usecase`): reglas de negocio.
- Adaptadores (`internal/transport/http`, `internal/repository/mysql`).
- Wiring (`internal/app`).

## Consequences
- Facilita unit testing con mocks.
- Cambios de framework/DB con impacto acotado.
- Mayor boilerplate inicial aceptado por estabilidad operativa.
