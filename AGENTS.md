# Code Review Rules

## Go
- Error handling explícito
- Variable names claros y descriptivos
- Funciones pequeñas y enfocadas
- Validación en boundaries (user input, external APIs)

## Templates (Templ)
- Componentes reutilizables para UI
- Validar disponibilidad de contexto
- Templates enfocados en presentación

## Middleware
- Separación clara de responsabilidades
- Documentar flujos de auth/CSRF
- Loguear eventos de seguridad
- Validar CSRF token en POST

## Security
- HTTPS en producción
- Argumenton2id para contraseñas (no plaintext)
- CSRF protection en formularios POST
- Rate limiting en login
- No exponer detalles internos en errores
