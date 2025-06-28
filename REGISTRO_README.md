# Sistema de Autenticación Completo (Registro + Login + JWT)

## Descripción

Este documento describe el sistema completo de autenticación implementado en el proyecto Pokedex Backend Go, incluyendo registro, login y autenticación JWT.

## Características Implementadas

### ✅ Funcionalidades Completadas

#### Sistema de Registro
1. **Hash de Contraseñas**: Las contraseñas se almacenan usando bcrypt con salt automático
2. **Validación de Email**: Validación básica del formato de email
3. **Validación de Contraseña**: Longitud mínima de 6 caracteres
4. **Prevención de Duplicados**: Verificación de emails únicos
5. **JWT Token**: Generación automática de token JWT al registrarse

#### Sistema de Login
1. **Autenticación por Email**: Login usando email y contraseña
2. **Verificación de Contraseña**: Comparación segura con bcrypt
3. **JWT Token**: Generación de token JWT al hacer login exitoso
4. **Manejo de Errores**: Respuestas apropiadas para credenciales inválidas

#### Sistema JWT
1. **Generación de Tokens**: Tokens JWT con expiración de 24 horas
2. **Validación de Tokens**: Middleware de autenticación para rutas protegidas
3. **Claims Personalizados**: Incluye user_id, email, issuer, etc.
4. **Middleware de Autenticación**: RequireAuth y OptionalAuth

#### General
1. **Manejo de Errores**: Respuestas HTTP apropiadas para diferentes tipos de errores
2. **Logging**: Registro detallado de operaciones y errores
3. **Estructura Limpia**: Separación clara entre capas (handler, service, repository)
4. **Base de Datos**: Migraciones automáticas y constraints apropiados

## Estructura de Archivos

```
domain/
├── login/
│   ├── handler/
│   │   └── login_handler.go         # Manejo de requests HTTP de login
│   ├── service/
│   │   └── login_service.go         # Lógica de negocio de login
│   ├── repository/
│   │   └── login_repository.go      # Acceso a datos de login
│   └── login.go                     # Provider de dependencias
├── register/
│   ├── handler/
│   │   └── register_handler.go      # Manejo de requests HTTP de registro
│   ├── service/
│   │   └── register_service.go      # Lógica de negocio de registro
│   ├── repository/
│   │   ├── register_repository.go   # Acceso a datos de registro
│   │   └── register_repository_test.go # Pruebas unitarias
│   └── register.go                  # Provider de dependencias
└── profile/
    ├── handler/
    │   └── profile_handler.go       # Manejo de requests HTTP de perfil
    ├── service/
    │   ├── profile_service.go       # Lógica de negocio de perfil
    │   └── profile_service_test.go  # Pruebas unitarias
    ├── repository/
    │   ├── profile_repository.go    # Acceso a datos de perfil
    │   └── profile_repository_test.go # Pruebas unitarias
    └── profile.go                   # Provider de dependencias

pkg/
├── auth/
│   ├── jwt.go                       # Servicio JWT
│   └── middleware.go                # Middleware de autenticación
├── model/
│   └── user.go                      # Modelo de usuario
└── dto/
    ├── login_response.go            # DTO para respuestas de login
    └── register_response.go         # DTO para respuestas de registro
```

## API Endpoints

### POST /api/v1/register

**Request Body:**
```json
{
  "email": "usuario@ejemplo.com",
  "password": "mipassword123"
}
```

**Response (201 Created):**
```json
{
  "user": {
    "id": "50c6410c-c61c-4ba2-a401-1f5c54d4087c",
    "email": "nuevousuario@ejemplo.com",
    "username": null,
    "created_at": "2025-06-28T16:24:32.990988-04:00",
    "updated_at": "2025-06-28T16:24:32.990988-04:00"
  },
  "message": "User registered successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### POST /api/v1/login

**Request Body:**
```json
{
  "email": "usuario@ejemplo.com",
  "password": "mipassword123"
}
```

**Response (200 OK):**
```json
{
  "user": {
    "id": "04eefc71-0637-40ce-8fb4-a05d50b817f1",
    "email": "usuario2@ejemplo.com",
    "username": null,
    "created_at": "2025-06-28T16:11:11.162556-04:00",
    "updated_at": "2025-06-28T16:11:11.162556-04:00"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### GET /api/v1/profile (Protegido)

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
  "user": {
    "id": "04eefc71-0637-40ce-8fb4-a05d50b817f1",
    "name": "Juan Pérez",
    "username": null,
    "email": "usuario2@ejemplo.com",
    "phone": "+1234567890",
    "created_at": "2025-06-28T16:11:11.162556-04:00",
    "updated_at": "2025-06-28T16:37:29.584993-04:00"
  },
  "message": "Profile retrieved successfully"
}
```

### PUT /api/v1/profile (Protegido)

**Headers:**
```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body (todos los campos son opcionales):**
```json
{
  "name": "Juan Carlos Pérez",
  "phone": "+1987654321",
  "username": "juanperez"
}
```

**Ejemplos de uso:**

1. **Actualizar solo un campo:**
```json
{
  "username": "nuevousername"
}
```

2. **Actualizar algunos campos (otros se mantienen):**
```json
{
  "name": "Nuevo Nombre",
  "phone": null
}
```

3. **Campos enviados como null no se actualizan:**
```json
{
  "name": "Juan Carlos",
  "username": null,
  "phone": null
}
```

**Response (200 OK):**
```json
{
  "user": {
    "id": "04eefc71-0637-40ce-8fb4-a05d50b817f1",
    "name": "Juan Carlos Pérez",
    "username": "juanperez",
    "email": "usuario2@ejemplo.com",
    "phone": "+1987654321",
    "created_at": "2025-06-28T16:11:11.162556-04:00",
    "updated_at": "2025-06-28T17:05:21.260172-04:00"
  },
  "message": "Profile updated successfully"
}
```

**Campos Permitidos para Actualización:**
- `name`: Nombre completo del usuario (opcional)
- `phone`: Número de teléfono del usuario (opcional)
- `username`: Nombre de usuario único (opcional)

**Validaciones:**
- `username` debe ser único en el sistema
- Los campos vacíos (`""`) no se consideran válidos y no se actualizan

**Nota:** Los campos `email`, `password`, `id`, `created_at`, `updated_at` no pueden ser actualizados a través de este endpoint por seguridad.

**Errores Posibles:**
- `400 Bad Request`: Email o password faltantes, formato inválido, password muy corta, no hay updates válidos
- `401 Unauthorized`: Credenciales inválidas, token inválido o faltante
- `404 Not Found`: Usuario no encontrado (solo profile)
- `409 Conflict`: Email ya existe (registro), username ya existe (profile update)
- `500 Internal Server Error`: Error del servidor

## Configuración de Base de Datos

La tabla `users` se crea automáticamente con la siguiente estructura:

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    username VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

## JWT Token Information

### Token Structure
Los JWT tokens incluyen los siguientes claims:
- `user_id`: ID único del usuario
- `email`: Email del usuario
- `iss`: Issuer (pokedex_backend_go)
- `sub`: Subject (user_id)
- `exp`: Expiration time (24 horas desde creación)
- `iat`: Issued at time
- `nbf`: Not before time

### Token Usage
Para usar endpoints protegidos, incluir el token en el header:
```
Authorization: Bearer <jwt-token>
```

## Cómo Probar

1. **Iniciar la base de datos:**
   ```bash
   docker-compose up postgres -d
   ```

2. **Ejecutar la aplicación:**
   ```bash
   go run cmd/api/main.go
   ```

3. **Probar el registro:**
   ```bash
   curl -X POST http://localhost:3000/api/v1/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

4. **Probar el login:**
   ```bash
   curl -X POST http://localhost:3000/api/v1/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

5. **Probar endpoint protegido (GET profile):**
   ```bash
   # Usar el token obtenido del login/registro
   curl -X GET http://localhost:3000/api/v1/profile \
     -H "Authorization: Bearer <jwt-token>"
   ```

6. **Probar actualización de perfil (PUT profile):**
   ```bash
   # Actualizar todos los campos
   curl -X PUT http://localhost:3000/api/v1/profile \
     -H "Authorization: Bearer <jwt-token>" \
     -H "Content-Type: application/json" \
     -d '{"name":"Juan Carlos Pérez","phone":"+1987654321","username":"juanperez"}'

   # Actualizar solo username
   curl -X PUT http://localhost:3000/api/v1/profile \
     -H "Authorization: Bearer <jwt-token>" \
     -H "Content-Type: application/json" \
     -d '{"username":"nuevousername"}'

   # Actualizar algunos campos (null no se modifica)
   curl -X PUT http://localhost:3000/api/v1/profile \
     -H "Authorization: Bearer <jwt-token>" \
     -H "Content-Type: application/json" \
     -d '{"name":"Nuevo Nombre","phone":null,"username":null}'
   ```

## Próximos Pasos Sugeridos

1. **JWT Authentication**: Implementar tokens JWT reales en lugar del placeholder
2. **Validación de Email Avanzada**: Usar regex más robusta o librería de validación
3. **Rate Limiting**: Prevenir ataques de fuerza bruta
4. **Email Verification**: Envío de emails de confirmación
5. **Tests de Integración**: Pruebas con base de datos real
6. **Campos Adicionales**: Implementar name, username, phone en el registro

## Seguridad

- ✅ Contraseñas hasheadas con bcrypt
- ✅ Contraseñas no se retornan en respuestas JSON
- ✅ Validación de entrada
- ✅ Prevención de duplicados
- ⚠️ Falta implementar rate limiting
- ⚠️ Falta validación de email más robusta
