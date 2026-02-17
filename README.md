# validacion-pases

Microservicio REST en Go (chi) para:
1. Emitir token JWT por `usuario/password` en `POST /v1/token`.
2. Validar Bearer JWT.
3. Insertar registros de pase con campos de negocio en MySQL.

## Endpoints
- `GET /healthz`
- `GET /readyz`
- `POST /v1/token` (sin token)
- `POST /v1/records` (requiere Bearer token)

## Flujo de autenticacion
1. Cliente llama `POST /v1/token` con `username` y `password`.
2. API valida contra `TOKEN_USERS`.
3. API emite JWT HS256 (`JWT_HS_SECRET`).
4. Cliente usa `Authorization: Bearer <token>` para `POST /v1/records`.

## Variables importantes
Ver `.env.example`.
- `JWT_ALG=HS256`
- `JWT_HS_SECRET=...`
- `TOKEN_USERS=user1:pass1,user2:pass2`
- `JWT_TOKEN_TTL=1h`

## Regla de negocio aplicada en guardado
- `EMISION`: `time.Now().UTC()`.
- `NAVE`: texto del request.
- `VIAJE`: texto del request.
- `CLIENTE`: texto del request.
- `BOOKING`: texto del request.
- `CONTENEDOR`:
  - `rama=internacional` -> `contenedor_serie`.
  - `rama=nacional` -> `1 X <codigo_iso>`.
  - si `rama` no viene, el backend la infiere (`contenedor_serie` => internacional, `codigo_iso+transportista` => nacional).
- `LIBRE_DE_RETENCION_HASTA`: `fecha_real + dias_libre`.
- `dias_libre`: si no viene, `0`.
- `TRANSPORTISTA`: requerido solo para `rama=nacional`.
- `TITULO_TERMINAL`: derivado de `puerto_descargue`.
- `USUARIO_FIRMA`: `sub` del JWT.

## Modelo MySQL (`records`)
Columnas:
- `id`
- `emision`
- `nave`
- `viaje`
- `cliente`
- `booking`
- `rama`
- `contenedor`
- `puerto_descargue`
- `libre_retencion_hasta`
- `dias_libre`
- `transportista`
- `titulo_terminal`
- `usuario_firma`
- `created_at`

## Ejemplo: emitir token
```bash
curl -X POST http://localhost:8080/v1/token \
  -H 'Content-Type: application/json' \
  -d '{"username":"apiuser","password":"change-me"}'
```

## Ejemplo: guardar record
```bash
curl -X POST http://localhost:8080/v1/records \
  -H "Authorization: Bearer <TOKEN>" \
  -H 'Content-Type: application/json' \
  -d '{
    "nave":"NAVE ABC",
    "viaje":"VJ-001",
    "cliente":"CLIENTE XYZ",
    "booking":"BK-123",
    "rama":"internacional",
    "contenedor_serie":"ABCU1234567",
    "fecha_real":"2026-02-09",
    "dias_libre":2,
    "puerto_descargue":"Balboa"
  }'
```

## Ejemplo: guardar record (payload legado compatible)
```bash
curl -X POST http://localhost:8080/v1/records \
  -H "Authorization: Bearer <TOKEN>" \
  -H 'Content-Type: application/json' \
  -d '{
    "emision":"2026-02-17 09:41:45",
    "nave":"NYK DENEB",
    "viaje":"072E",
    "cliente":"CAPITAL PACIFICO, S.A.",
    "booking":"YMLUL160382911",
    "rama":"internacional",
    "contenedor":"YMLU5374938",
    "puerto_descargue":"RODMAN",
    "libre_retencion_hasta":"2021-03-06",
    "dias_libre":0,
    "transportista":"GLOBERUNNERS, INC",
    "titulo_terminal":"PANAMA PORTS COMPANY (RODMAN)",
    "usuario_firma":"Admin"
  }'
```

## Desarrollo local
```bash
cp .env.example .env
docker compose up --build
```

## Calidad
```bash
make fmt
make lint
make test
make test-integration
make coverage
```
