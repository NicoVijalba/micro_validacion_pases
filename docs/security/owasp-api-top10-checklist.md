# OWASP API Security Top 10 Checklist

- API1 Broken Object Level Auth: no expone acceso por ID sin auth; endpoint protegido con JWT.
- API2 Broken Authentication: validacion de firma y claims; sin tokens en logs.
- API3 Broken Object Property Level Auth: payload estricto y campos permitidos.
- API4 Unrestricted Resource Consumption: rate limit, timeout, body max.
- API5 Broken Function Level Auth: endpoint de escritura solo autenticado.
- API6 Unrestricted Access to Sensitive Business Flows: throttling por IP.
- API7 SSRF: sin fetch dinamico de URLs de usuario.
- API8 Security Misconfiguration: headers hardening + CORS controlado.
- API9 Improper Inventory Management: OpenAPI y versionado `/v1`.
- API10 Unsafe Consumption of APIs: dependencia JWKS con TLS; validacion estricta.
