# Threat Model (STRIDE)

## Activos
- Token de acceso JWT
- Datos de records en MySQL
- Imagen de contenedor y pipeline
- Credenciales de despliegue

## STRIDE y mitigaciones
- Spoofing: validacion JWT RS256, issuer/audience strict, TLS extremo a extremo.
- Tampering: firma JWT, SQL parametrizado, imagen escaneada con Trivy.
- Repudiation: logs JSON con request-id, timestamps UTC.
- Information disclosure: no loggear tokens, headers de seguridad, secretos fuera del repo.
- Denial of Service: rate limit, body limit, timeouts, throttle.
- Elevation of privilege: principle of least privilege para DB y CI secrets, branch protection.

## Riesgos residuales
- Caida del proveedor JWKS puede afectar refresco de claves.
- Ataques volumetricos requieren capa adicional (WAF/CDN) fuera de este servicio.
