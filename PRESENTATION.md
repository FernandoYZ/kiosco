# Kiosco Mobile — Presentación Canva

**Autor**: Fernando  
**Fecha**: 19 de Mayo 2026  
**Proyecto**: Kiosco Mobile App (Android MVP)

---

## SLIDE 1: Portada

**KIOSCO MOBILE**  
*Control de Consumo Escolar en tu Bolsillo*

Versión 1.0 MVP  
Android 7.0+ | Offline-First | Apple-like Design

---

## SLIDE 2: El Problema

❌ **Situación actual:**
- Plataforma Kiosco funciona solo en web
- Estudiantes necesitan acceso desde móvil/tablet
- Sin conexión a internet = sin acceso a su resumen
- Compartir datos requiere manualmente

**Solución:** App nativa Android que funciona offline

---

## SLIDE 3: Visión del Producto

✅ **Kiosco Mobile**
- App nativa Android (Kotlin)
- Acceso 100% offline a consumos personales
- Interfaz Apple-like (minimalista, limpia)
- Compartir resumen por WhatsApp en 1 tap
- Sincronización manual de datos

**Casos de uso:**
1. Estudiante abre app → ve su resumen → comparte por WhatsApp
2. Padre recibe comprobante visual de consumos
3. Sin internet = sigue funcionando con datos descargados

---

## SLIDE 4: Arquitectura: Offline-First

```
Primer Launch:
[Login] → [Descargar BD] → [App funciona 100% local]
             ↓ (gzip)
        SQLite local

Después:
[Botón "Sincronizar"] → [Re-descarga .db completo]
                           (solo si hay conexión)
```

**Ventajas:**
- ⚡ Rápido (sin latencia de red)
- 📱 Funciona sin internet
- 🔄 Sincronización manual (control del usuario)
- 🎯 MVP simple (hoy lanzamos esto)

---

## SLIDE 5: Stack Técnico Backend

**Endpoint nuevo**: `GET /api/database/download`

```
Go 1.26.2
└─ internal/controllers/api_controller.go
   └─ DescargarBD()
      ├─ Lee .db comprimido (gzip)
      ├─ Requiere autenticación
      └─ Headers para descarga
```

**Integración:**
- Usa middleware de autenticación existente
- Compatible con token HMAC-SHA256 actual
- 5 minutos de implementación ✓ HECHO

---

## SLIDE 6: Stack Técnico Frontend

**Lenguaje**: Kotlin  
**UI Framework**: Material Design 3 (customizado Apple-like)  
**Persistencia**: Room 3.0 + SQLite local  
**Autenticación**: DataStore encriptado (Google Tink)  
**Networking**: Retrofit 2 + OkHttp 5  

**Versiones (2026 actualizadas)**:
- Room 3.0.0 (nuevo)
- OkHttp 5.2.0 (5x más rápido)
- Lifecycle 2.11.0
- Navigation 2.8.0

---

## SLIDE 7: Design System

**Material Design 3 + Custom Theme**

```
Material Design 3 (oficial Android)
        ↓
Custom Theme: Theme.Kiosco
        ↓
Apple-like visual (minimalista, limpio)
```

**Colores:**
- Blanco puro (#FFFFFF)
- Azul iOS (#0A84FF)
- Grises suaves (#F5F5F7, #8E8E93)
- Bordes redondeados (12-16dp)

**Typography:**
- Roboto (Android standard)
- Sizes grandes y legibles
- Line spacing generoso

**Componentes:**
- ✅ Bottom Navigation (iOS-style)
- ✅ Material Cards (redondeadas)
- ✅ Material Buttons (12dp radius)
- ✅ Text Inputs (redondeados, sin bordes visibles)

---

## SLIDE 8: Features MVP

### 1️⃣ Autenticación
- Login con email + password
- POST `/login` (backend existente)
- Token guardado en DataStore encriptado

### 2️⃣ Descarga de BD
- `GET /api/database/download`
- Descompresión automática (gzip)
- Guardado en SQLite local

### 3️⃣ Consumos
- Lista de consumos personales
- Filtro por semana/fecha
- Total semanal calculado localmente

### 4️⃣ Pagos
- Lista de pagos registrados
- Estado (pagado/pendiente)
- Detalles por consumo

### 5️⃣ Resumen
- Vista visual de consumos
- Desglose por sector/producto
- **[Compartir por WhatsApp]** ← Feature star ⭐

### 6️⃣ Sincronizar
- Botón "Actualizar datos"
- Re-descarga .db completo
- Reemplaza versión local

---

## SLIDE 9: Feature Star: Compartir por WhatsApp

**Flujo:**

```
[Resumen] 
    ↓
[Click "Compartir por WhatsApp"]
    ↓
[App genera imagen]
    (Canvas + Bitmap)
    ├─ Nombre del estudiante
    ├─ Total consumido
    ├─ Desglose por sector
    └─ Fecha generación
    ↓
[Intent.ACTION_SEND]
    ↓
[Abre WhatsApp]
    ↓
[Usuario elige contacto/grupo]
    ↓
[Envía imagen + mensaje]
```

**Resultado:** 
Padre recibe comprobante visual en 5 segundos

---

## SLIDE 10: Responsividad (Mobile + Tablet)

**Mobile (normal):**
- Bottom Navigation (4 tabs)
- Full-width cards
- Textos grandes, toques generosos

**Tablet (landscape):**
- Sidebar navigation (80dp)
- Content panel al lado
- Master-detail layout

**Tecnología:**
- ConstraintLayout flexible
- res/layout-sw720dp para tablets
- Material adaptive UI

---

## SLIDE 11: Estructura del Proyecto

```
app/src/main/java/com/kiosco/
├── ui/
│   ├── activities/
│   │   ├── LoginActivity.kt
│   │   └── MainActivity.kt
│   ├── fragments/
│   │   ├── ConsumosFragment.kt
│   │   ├── PagosFragment.kt
│   │   ├── ResumenFragment.kt
│   │   └── SincronizarFragment.kt
│   └── adapters/
├── data/
│   ├── db/
│   │   ├── KioscoDatabase.kt (Room)
│   │   ├── dao/
│   │   └── entities/
│   ├── repositories/
│   └── api/
├── viewmodels/
├── utils/
└── MyApplication.kt
```

**Patrón:** MVVM + Repository (Clean Architecture)

---

## SLIDE 12: Flujo de Usuario MVP

```
┌─────────────────────────────┐
│   1. LoginActivity          │
│   Email + Password          │
└──────────────┬──────────────┘
               ↓
┌─────────────────────────────┐
│   2. SyncFragment           │
│   Descargando BD...         │
└──────────────┬──────────────┘
               ↓
┌─────────────────────────────┐
│   3. MainActivity           │
│   Bottom Nav (4 opciones)   │
├─────────────────────────────┤
│ ✓ Consumos (default)        │
│ ✓ Pagos                     │
│ ✓ Resumen [Compartir] ⭐    │
│ ✓ Sincronizar              │
└─────────────────────────────┘
```

---

## SLIDE 13: Seguridad & Autenticación

**Token HMAC-SHA256** (existente en backend)

```
Formato:
base64url(idUsuario:puede_editar:expiry_unix) . 
base64url(HMAC-SHA256)
```

**En la app:**
- ✅ Almacenado en DataStore encriptado (Google Tink)
- ✅ Pasado en header Cookie a cada request
- ✅ Validado antes de descargar BD
- ✅ Borrado al logout

**No es necesario** parsear el token en la app.  
Solo: guardar → pasar → validar

---

## SLIDE 14: Dependencias 2026

**Actualizado a versiones 2026 (NOT 2024!)**

| Librería | Versión |
|----------|---------|
| Room | 3.0.0 ⭐ |
| OkHttp | 5.2.0 ⭐ |
| Retrofit | 2.11.0 |
| Lifecycle | 2.11.0 |
| Navigation | 2.8.0 |
| DataStore | 1.1.1 |
| Material Design | 1.12.0 |
| Kotlin | 2.0.x |

**Nota:** EncryptedSharedPreferences deprecado → DataStore (Google official)

---

## SLIDE 15: Herramientas & Versiones

**Android Studio**: Panda 4 ✅  
**Target SDK**: 34+ (Android 14)  
**Minimum SDK**: 24 (Android 7.0) ✅  
**Build System**: Gradle 8.x  
**Kotlin**: 2.0.x  
**Compile Target**: API 34  

**Template recomendado:** Empty Activity (control total)

---

## SLIDE 16: Timeline MVP

```
Fase 1: Setup (1-2 horas)
├─ Crear proyecto Android Studio
├─ Configurar build.gradle.kts
├─ Copiar colores + theme
└─ Copiar dependencias

Fase 2: Core (4-6 horas)
├─ LoginActivity + autenticación
├─ DatabaseManager (descargar/descomprimir)
├─ Room setup + entidades
└─ MainActivity + bottom nav

Fase 3: Features (4-6 horas)
├─ ConsumosFragment + ViewModel
├─ PagosFragment + ViewModel
├─ ResumenFragment + Canvas
└─ Compartir por WhatsApp

Fase 4: Polish (1-2 horas)
├─ Responsive design (tablet)
├─ Testing manual
└─ APK generado

Total: ~12-18 horas (1-2 días)
```

---

## SLIDE 17: What's Done ✅

**Backend**
- ✅ Endpoint `/api/database/download` creado
- ✅ Compresión gzip integrada
- ✅ Autenticación validada
- ✅ Code compila sin errores

**Documentación**
- ✅ `prompt.md` (600+ líneas)
  - Arquitectura completa
  - Design system Apple-like
  - Stack técnico 2026
  - Código Kotlin de ejemplo
  - Layouts XML listos para copiar
- ✅ `PRESENTATION.md` (esta diapositiva)

---

## SLIDE 18: What's Next ❓

**Usuario:**
1. Abre Android Studio Panda 4
2. File → New → New Android Project
3. Choose: Empty Activity
4. minSdk = 24, targetSdk = 34
5. Sigue el `prompt.md` punto por punto

**Semanal:** MVP launch en producción  
**Post-MVP:** Sincronización inteligente (deltas)

---

## SLIDE 19: Diferenciales

✨ **Por qué esto es especial:**

1. **Offline-first**: Funciona sin internet
2. **Design cohesivo**: Web + Mobile mismo visual
3. **Rápido**: MVP en 1-2 días
4. **Escalable**: Arquitectura preparada para crecer
5. **Seguro**: Token + DataStore encriptado
6. **Moderno**: 2026 stack (Room 3, OkHttp 5)
7. **Responsive**: Mobile + Tablet soportados
8. **Apple-like**: Minimalista, limpio, intuitivo

---

## SLIDE 20: Business Impact

**Estudiantes:**
- 📱 Acceso desde móvil (sin web)
- ⚡ Rápido, sin lag (offline)
- 📤 Comparten resumen en 5 segundos

**Padres:**
- 📸 Reciben comprobante visual (WhatsApp)
- ✅ Confianza en datos reales
- 🕐 Respuesta inmediata (no esperan email)

**Institución:**
- 🎯 Engagement mejorado
- 💪 Plataforma moderna (iOS-style UX)
- 🔄 Sincronización confiable (offline-first)

---

## SLIDE 21: Métricas & KPIs

**Esperado en 30 días:**
- 📲 80% de estudiantes con app instalada
- ⭐ 4.5+ rating (Apple-like UX)
- 📤 50% comparten resumen por WhatsApp
- ⚡ <2s para abrir consumos (offline)
- 🔄 0% errores de sincronización

---

## SLIDE 22: Riesgos & Mitigación

| Riesgo | Probabilidad | Mitigación |
|--------|--------------|-----------|
| Sincronización falla | Baja | Validación en descarga, logs detallados |
| Permisos Android | Baja | Runtime permissions implementadas |
| Performance tablet | Baja | ConstraintLayout flexible, testing |
| Token expirado | Media | Refresh flow + re-login |
| WhatsApp no instalado | Baja | Fallback a Intent fallback |

---

## SLIDE 23: Roadmap Post-MVP

**V1.1 (Semana 2)**
- Caché de imágenes compartidas
- Historial de sincronizaciones
- Modo oscuro (opcional)

**V1.2 (Semana 3)**
- Sincronización incremental (deltas)
- Menor uso de datos
- Más rápida

**V1.3 (Semana 4)**
- Push notifications (cambios en consumos)
- Alerts de límite de presupuesto

**V2.0 (Mes 2)**
- Edición local de consumos
- Sincronización bidireccional
- Estadísticas históricas

---

## SLIDE 24: Stack Completo (Resumen Visual)

```
┌─────────────────────────────────────────┐
│         FRONTEND (Android)              │
├─────────────────────────────────────────┤
│ Material Design 3 + Custom Theme        │
│ Kotlin + Room 3.0 + DataStore           │
│ Retrofit 2 + OkHttp 5 (Networking)      │
│ MVVM Architecture                       │
└────────────┬────────────────────────────┘
             ↓
┌─────────────────────────────────────────┐
│    BACKEND (Go) — Existente             │
├─────────────────────────────────────────┤
│ GET /api/database/download              │
│ SQLite (WAL mode)                       │
│ HMAC-SHA256 autenticación               │
└─────────────────────────────────────────┘
```

---

## SLIDE 25: Conclusión

**Kiosco Mobile es:**
- 🎯 MVP claro, definido, alcanzable
- ⚡ Rápido de desarrollar (1-2 días)
- 📱 Experiencia Apple-like en Android
- 🔒 Seguro, offline-first, escalable
- 📈 Impacto real en estudiantes & padres

**Status:** Ready to build 🚀  
**Próximo paso:** Abrir Android Studio

---

## SLIDE 26: Preguntas?

**Contacto:**
- Documentación: `prompt.md` (600+ líneas)
- Código: `/home/fernando/Works/kiosco`
- Backend: `GET /api/database/download` ✅

**¿Dudas sobre:**
- Arquitectura?
- Design system?
- Stack técnico?
- Timeline?

---

## SLIDE 27: Agradecimiento

✨ **Built with:**
- Android Jetpack (Room, DataStore, Navigation)
- Material Design 3
- Modern Kotlin practices
- Offline-first architecture

**Gracias por la atención.**

🚀 **Vamos a hacerlo realidad.**

