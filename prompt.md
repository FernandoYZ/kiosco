# Kiosco Mobile — Prompt de Desarrollo

**Versión**: MVP 1.0  
**Stack**: Android 7.0+ (API 24+), Android Studio Panda 4, Kotlin  
**Arquitectura**: Offline-first con sincronización local  
**Última actualización**: 2026-05-19

---

## Visión General

Kiosco Mobile es la versión nativa de Android de la plataforma de control de consumo escolar. Los usuarios descarga la BD una sola vez, consultan datos localmente (sin conexión a internet requerida), y pueden compartir su resumen de consumos por WhatsApp.

**No es un replicate del web en mobile.** Es una app optimizada para:
- Lectura rápida de consumos personales
- Compartir resumen por WhatsApp (imagen generada localmente)
- Funcionar offline 100%
- Sincronización manual cuando hay conexión

---

## Arquitectura de Datos

### Sincronización: Descarga Inicial + Offline-First

1. **Primer launch**:
   - App solicita `GET /api/database/download` con token de autenticación
   - Backend responde con `.db` comprimido (gzip)
   - App descomprime y almacena en `Context.getDatabasesPath()`

2. **Después del primer launch**:
   - Todas las queries usan SQLite local
   - Botón "Sincronizar" en UI → re-descarga .db completo
   - No hay conexión a internet requerida para navegar/consultar

3. **Autenticación local**:
   - Durante login, app recibe token HMAC-SHA256 en cookie
   - Token se almacena en **DataStore encriptado** (reemplaza EncryptedSharedPreferences deprecado desde 2025)
   - Token se pasa a cada request hacia `/api/database/download`

### Base de Datos Local

**Almacenamiento**: SQLite en `Context.getDatabasesPath()` (manejado por Room)

**Esquema**: Idéntico al del servidor (kiosco/internal/config/schema.sql)
- Tablas: `estudiantes`, `productos`, `consumos`, `pagos`
- Relaciones: consumos/pagos se ligan a estudiante_id
- Índices: mismo esquema que producción

**Herramienta ORM**: Room (Android Jetpack)
- DAOs para cada tabla (EstudianteDao, ConsumoDao, etc.)
- Migraciones automáticas al descargar nueva versión de .db

---

## Stack Técnico

### Frontend
- **Lenguaje**: Kotlin
- **UI Framework**: XML + ViewBinding (no Jetpack Compose)
- **Layouts**: ConstraintLayout para móvil + tablet responsivo
- **Componentes**: Material Design 3 (AndroidX)

### Persistencia
- **Room 3.0** (AndroidX): ORM para SQLite
- **DataStore** (AndroidX): almacenar token de forma segura con Tink (reemplaza EncryptedSharedPreferences deprecado)
- **WorkManager** (opcional): sincronización en background

### Networking
- **Retrofit 2** + **OkHttp 4**: descargar .db
- **Gson**: parseo de JSON (si hay endpoint futuro)

### Compartir
- **Intent.ACTION_SEND**: compartir por WhatsApp
- **Canvas/Bitmap**: generar imagen de resumen

---

## Design System: Apple-like UI con Material Design 3

**¿Qué UI framework?** Material Design 3 (AndroidX Material library)  
**¿Cómo se ve Apple-like?** Custom theming + componentes Material customizados

Kiosco Mobile usa **Material Design 3 como base** pero aplicamos un tema personalizado para lograr la **estética limpia y minimalista del web**, manteniendo compatibilidad con Android 7.0+.

### Estrategia: Material Design 3 + Custom Theme

Material Design 3 proporciona:
- Componentes modernos (buttons, cards, input fields, navigation)
- System manejo de theming (light/dark, dinamamic color)
- Accesibilidad integrada
- Compatibilidad con AndroidX

Customización Apple-like:
- Colores: Blanco, grises suaves, azul limpio (no colores vibrantes)
- Tipografía: Roboto con sizes generosas (Apple = legibilidad)
- Bordes: Material rounded corners (12-16dp) + sin sombras dramáticas
- Espaciado: Padding generoso (16-24dp), sin aglomeración
- Deshabilitar dark mode automático (Apple = light mode classic)

### Principios de Diseño
- **Minimalismo**: Espacio en blanco generoso, sin ornamentos
- **Tipografía clara**: Roboto, tamaños grandes y legibles
- **Colores neutros**: Blanco, gris suave, acentos azul claro
- **Redondeado**: Material rounded corners (12-16dp)
- **Sombras sutiles**: Elevación mínima, no dramática
- **Interacción natural**: Animaciones smooth, feedback táctil

### Paleta de Colores

```xml
<!-- res/values/colors.xml -->
<color name="colorPrimary">#0A84FF</color>        <!-- Azul iOS -->
<color name="colorSecondary">#5AC8FA</color>     <!-- Azul claro -->
<color name="colorBackground">#FFFFFF</color>   <!-- Blanco puro -->
<color name="colorSurface">#F5F5F7</color>      <!-- Gris muy claro -->
<color name="colorError">#FF3B30</color>        <!-- Rojo alerta -->
<color name="colorSuccess">#34C759</color>      <!-- Verde éxito -->
<color name="textPrimary">#1C1C1E</color>       <!-- Negro oscuro -->
<color name="textSecondary">#8E8E93</color>     <!-- Gris medio -->
<color name="textTertiary">#D1D1D6</color>      <!-- Gris claro -->
<color name="divider">#E5E5EA</color>           <!-- Líneas grises -->
```

### Configuración del Theme (AndroidManifest.xml)

```xml
<!-- AndroidManifest.xml -->
<application
    android:theme="@style/Theme.Kiosco"
    ...>
```

### Theme Principal (Material Design 3 Custom)

```xml
<!-- res/values/themes.xml -->
<resources>
    <!-- Tema base para Kiosco (Light only, sin dark mode) -->
    <style name="Theme.Kiosco" parent="Theme.Material3.Light">
        <!-- Colores primarios -->
        <item name="colorPrimary">@color/colorPrimary</item>
        <item name="colorSecondary">@color/colorSecondary</item>
        <item name="colorTertiary">@color/colorSecondary</item>
        
        <!-- Fondo -->
        <item name="android:colorBackground">@color/colorBackground</item>
        <item name="colorSurface">@color/colorSurface</item>
        
        <!-- Text colors -->
        <item name="android:textColorPrimary">@color/textPrimary</item>
        <item name="android:textColorSecondary">@color/textSecondary</item>
        
        <!-- Shape (bordes redondeados) -->
        <item name="shapeAppearanceSmallComponent">@style/ShapeAppearance.Kiosco.SmallComponent</item>
        <item name="shapeAppearanceMediumComponent">@style/ShapeAppearance.Kiosco.MediumComponent</item>
        <item name="shapeAppearanceLargeComponent">@style/ShapeAppearance.Kiosco.LargeComponent</item>
        
        <!-- Botón style -->
        <item name="buttonStyle">@style/Widget.Kiosco.Button</item>
        
        <!-- Dark mode deshabilitado (Apple-like = light only) -->
        <item name="android:forceDarkAllowed">false</item>
        
        <!-- Action bar (ocultar si no se usa) -->
        <item name="windowNoTitle">true</item>
        <item name="windowActionBar">false</item>
    </style>

    <!-- Shape styles (bordes redondeados) -->
    <style name="ShapeAppearance.Kiosco.SmallComponent" parent="ShapeAppearance.Material3.SmallComponent">
        <item name="cornerSize">8dp</item>
    </style>
    <style name="ShapeAppearance.Kiosco.MediumComponent" parent="ShapeAppearance.Material3.MediumComponent">
        <item name="cornerSize">12dp</item>
    </style>
    <style name="ShapeAppearance.Kiosco.LargeComponent" parent="ShapeAppearance.Material3.LargeComponent">
        <item name="cornerSize">16dp</item>
    </style>

    <!-- Button style customizado -->
    <style name="Widget.Kiosco.Button" parent="Widget.Material3.Button">
        <item name="cornerRadius">12dp</item>
        <item name="android:textSize">16sp</item>
        <item name="android:paddingStart">16dp</item>
        <item name="android:paddingEnd">16dp</item>
        <item name="android:paddingTop">12dp</item>
        <item name="android:paddingBottom">12dp</item>
    </style>
</resources>
```

### Tipografía

```xml
<!-- res/values/styles.xml -->
<style name="TextHeading1" parent="TextAppearance.Material3.Headline1">
    <item name="fontFamily">@font/roboto_bold</item>
    <item name="android:textSize">32sp</item>
    <item name="android:textColor">@color/textPrimary</item>
    <item name="android:lineSpacingExtra">4dp</item>
</style>

<style name="TextHeading2" parent="TextAppearance.MaterialComponents.Headline2">
    <item name="fontFamily">@font/roboto_bold</item>
    <item name="android:textSize">24sp</item>
    <item name="android:textColor">@color/textPrimary</item>
</style>

<style name="TextBody" parent="TextAppearance.MaterialComponents.Body1">
    <item name="fontFamily">@font/roboto_regular</item>
    <item name="android:textSize">16sp</item>
    <item name="android:textColor">@color/textPrimary</item>
    <item name="android:lineSpacingExtra">2dp</item>
</style>

<style name="TextCaption" parent="TextAppearance.MaterialComponents.Caption">
    <item name="fontFamily">@font/roboto_regular</item>
    <item name="android:textSize">12sp</item>
    <item name="android:textColor">@color/textSecondary</item>
</style>
```

### Componentes Principales

#### 1. Botón (Apple-style)
```xml
<!-- res/drawable/button_primary.xml -->
<shape xmlns:android="http://schemas.android.com/apk/res/android">
    <solid android:color="@color/colorPrimary" />
    <corners android:radius="12dp" />
</shape>

<!-- En layout: -->
<Button
    android:id="@+id/btnLogin"
    android:layout_width="match_parent"
    android:layout_height="48dp"
    android:text="Iniciar sesión"
    android:textStyle="bold"
    android:textSize="16sp"
    android:background="@drawable/button_primary"
    android:textColor="@android:color/white"
    android:layout_margin="16dp"
    android:stateListAnimator="@animator/button_elevation" />
```

#### 2. Card (Consumo/Pago)
```xml
<com.google.android.material.card.MaterialCardView
    android:layout_width="match_parent"
    android:layout_height="wrap_content"
    android:layout_margin="8dp"
    app:cardCornerRadius="12dp"
    app:cardElevation="2dp"
    app:cardBackgroundColor="@color/colorBackground"
    app:strokeColor="@color/divider"
    app:strokeWidth="1dp">

    <LinearLayout
        android:layout_width="match_parent"
        android:layout_height="wrap_content"
        android:orientation="vertical"
        android:padding="16dp">

        <TextView
            android:id="@+id/productName"
            android:layout_width="match_parent"
            android:layout_height="wrap_content"
            style="@style/TextBody"
            android:textStyle="bold" />

        <TextView
            android:id="@+id/productPrice"
            android:layout_width="match_parent"
            android:layout_height="wrap_content"
            style="@style/TextCaption"
            android:layout_marginTop="8dp"
            android:textColor="@color/textSecondary" />
    </LinearLayout>
</com.google.android.material.card.MaterialCardView>
```

#### 3. Entrada de Texto (Apple-style)
```xml
<com.google.android.material.textfield.TextInputLayout
    android:id="@+id/emailLayout"
    android:layout_width="match_parent"
    android:layout_height="wrap_content"
    android:layout_margin="16dp"
    app:boxCornerRadiusBottomEnd="12dp"
    app:boxCornerRadiusBottomStart="12dp"
    app:boxCornerRadiusTopEnd="12dp"
    app:boxCornerRadiusTopStart="12dp"
    app:boxStrokeColor="@color/divider"
    app:boxBackgroundColor="@color/colorSurface"
    app:hintTextColor="@color/textSecondary">

    <com.google.android.material.textfield.TextInputEditText
        android:id="@+id/emailInput"
        android:layout_width="match_parent"
        android:layout_height="48dp"
        android:hint="Correo electrónico"
        android:inputType="textEmailAddress"
        android:textSize="16sp"
        android:fontFamily="@font/roboto_regular" />
</com.google.android.material.textfield.TextInputLayout>
```

#### 4. Bottom Navigation (iOS-like)
```xml
<!-- res/menu/bottom_nav_menu.xml -->
<menu xmlns:android="http://schemas.android.com/apk/res/android">
    <item
        android:id="@+id/nav_consumos"
        android:icon="@drawable/ic_consumos"
        android:title="Consumos" />
    <item
        android:id="@+id/nav_pagos"
        android:icon="@drawable/ic_pagos"
        android:title="Pagos" />
    <item
        android:id="@+id/nav_resumen"
        android:icon="@drawable/ic_resumen"
        android:title="Resumen" />
    <item
        android:id="@+id/nav_settings"
        android:icon="@drawable/ic_settings"
        android:title="Ajustes" />
</menu>

<!-- En MainActivity layout: -->
<com.google.android.material.bottomnavigation.BottomNavigationView
    android:id="@+id/bottomNav"
    android:layout_width="match_parent"
    android:layout_height="56dp"
    android:layout_alignParentBottom="true"
    android:background="@color/colorBackground"
    app:menu="@menu/bottom_nav_menu"
    app:labelVisibilityMode="label_visibility_labeled"
    app:itemIconTint="@color/textSecondary"
    app:itemTextColor="@color/textSecondary" />
```

### Responsive Design (Mobile + Tablet)

#### Layout para Mobile (normal)
```xml
<!-- res/layout/activity_main.xml -->
<FrameLayout
    android:layout_width="match_parent"
    android:layout_height="match_parent">

    <FrameLayout
        android:id="@+id/fragmentContainer"
        android:layout_width="match_parent"
        android:layout_height="match_parent"
        android:layout_marginBottom="56dp" />

    <com.google.android.material.bottomnavigation.BottomNavigationView
        android:id="@+id/bottomNav"
        android:layout_width="match_parent"
        android:layout_height="56dp"
        android:layout_gravity="bottom" />
</FrameLayout>
```

#### Layout para Tablet (landscape)
```xml
<!-- res/layout-sw720dp/activity_main.xml -->
<LinearLayout
    android:layout_width="match_parent"
    android:layout_height="match_parent"
    android:orientation="horizontal">

    <!-- Sidebar Navigation (drawer/rail) -->
    <androidx.recyclerview.widget.RecyclerView
        android:id="@+id/navRail"
        android:layout_width="80dp"
        android:layout_height="match_parent"
        android:background="@color/colorSurface"
        android:padding="8dp" />

    <!-- Content -->
    <FrameLayout
        android:id="@+id/fragmentContainer"
        android:layout_width="0dp"
        android:layout_height="match_parent"
        android:layout_weight="1" />
</LinearLayout>
```

**Nota**: Material Design 3 en AndroidX maneja automáticamente temas claros/oscuros. Para mantener design Apple-like, deshabilitar dark mode en `styles.xml`:

```xml
<item name="android:forceDarkAllowed">false</item>
```

---

## Flujo Visual del MVP

```
LoginActivity
    ↓
  [Email] [Password]
    ↓
  [Login Button]
    ↓
SyncFragment (Progress bar)
    ↓
  Descarga .db...
    ↓
MainActivity (Bottom Nav)
    ├─ ConsumosFragment
    │   ├─ Card por consumo
    │   └─ Total semanal
    ├─ PagosFragment
    │   └─ Card por pago
    ├─ ResumenFragment
    │   ├─ Gráfico simple (o números grandes)
    │   └─ [Compartir por WhatsApp]
    └─ SettingsFragment
        └─ [Sincronizar BD]
```

---

## Módulos y Estructura

```
app/
├── src/main/
│   ├── java/com/kiosco/
│   │   ├── ui/
│   │   │   ├── activities/
│   │   │   │   ├── LoginActivity.kt
│   │   │   │   ├── MainActivity.kt
│   │   │   │   └── ResumenActivity.kt
│   │   │   ├── fragments/
│   │   │   │   ├── ConsumosFragment.kt
│   │   │   │   ├── PagosFragment.kt
│   │   │   │   ├── ResumenFragment.kt
│   │   │   │   └── SincronizarFragment.kt
│   │   │   └── adapters/
│   │   │       ├── ConsumosAdapter.kt
│   │   │       └── PagosAdapter.kt
│   │   ├── data/
│   │   │   ├── db/
│   │   │   │   ├── KioscoDatabase.kt
│   │   │   │   ├── dao/
│   │   │   │   │   ├── EstudianteDao.kt
│   │   │   │   │   ├── ConsumoDao.kt
│   │   │   │   │   ├── PagoDao.kt
│   │   │   │   │   └── ProductoDao.kt
│   │   │   │   └── entities/
│   │   │   │       ├── EstudianteEntity.kt
│   │   │   │       ├── ConsumoEntity.kt
│   │   │   │       ├── PagoEntity.kt
│   │   │   │       └── ProductoEntity.kt
│   │   │   ├── repositories/
│   │   │   │   └── KioscoRepository.kt
│   │   │   └── api/
│   │   │       ├── RetrofitClient.kt
│   │   │       └── KioscoApiService.kt
│   │   ├── viewmodels/
│   │   │   ├── LoginViewModel.kt
│   │   │   ├── ConsumosViewModel.kt
│   │   │   ├── ResumenViewModel.kt
│   │   │   └── SincronizarViewModel.kt
│   │   ├── utils/
│   │   │   ├── AuthManager.kt
│   │   │   ├── DatabaseManager.kt
│   │   │   ├── ResumenGenerator.kt
│   │   │   └── Constants.kt
│   │   └── MyApplication.kt
│   └── res/
│       ├── layout/
│       ├── drawable/
│       ├── values/
│       └── menu/
├── build.gradle.kts
└── AndroidManifest.xml
```

---

## Features Principales

### 1. Autenticación (LoginActivity)
- Input: email + password
- POST a `{SERVER}/login`
- Si éxito: recibe cookie con token HMAC-SHA256
- Almacena en EncryptedSharedPreferences
- Navega a descarga de .db

### 2. Descarga e Instalación de BD (SincronizarFragment)
- GET `/api/database/download` (autenticado con cookie)
- Descomprime gzip en tiempo real
- Copia a `getDatabasesPath()/database.db`
- Inicializa Room
- Navega a MainActivity (consumos)

### 3. Consumos (ConsumosFragment)
- Lista de consumos del estudiante logueado
- Filtro por fecha/semana (como web)
- Cálculo de total semanal en tiempo real desde BD local

### 4. Pagos (PagosFragment)
- Lista de pagos registrados
- Estado (pagado/pendiente)
- Detalles por consumo

### 5. Resumen (ResumenFragment)
- Resumen visual de consumos actuales
- Detalles por sector/producto
- **Compartir por WhatsApp**:
  - Genera imagen (Canvas) con:
    - Nombre estudiante
    - Total consumido
    - Desglose por sector
    - Fecha de generación
  - Abre Intent para WhatsApp con mensaje + imagen

### 6. Sincronizar (Settings/Menu)
- Botón "Actualizar BD"
- Re-descarga `.db` completo
- Reemplaza base de datos local
- Toast de confirmación

---

## Autenticación: Token HMAC-SHA256

**Formato en servidor**: `base64url(idUsuario:puede_editar:expiry_unix).base64url(HMAC-SHA256)`

**Flujo app**:
1. POST `/login` con credenciales
2. Servidor responde con cookie `kiosco_token` (HttpOnly en web, pero app lo lee)
3. App extrae cookie → almacena en **DataStore encriptado** (Google Tink)
4. Para cada request: pasa cookie en header `Cookie: kiosco_token=...`

**Implementación con DataStore (2026 standard)**:

```kotlin
// utils/AuthManager.kt
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import kotlinx.coroutines.flow.map

class AuthManager(private val dataStore: DataStore<Preferences>) {
    companion object {
        val TOKEN_KEY = stringPreferencesKey("kiosco_token")
    }

    suspend fun saveToken(token: String) {
        dataStore.edit { preferences ->
            preferences[TOKEN_KEY] = token
        }
    }

    fun getToken() = dataStore.data.map { preferences ->
        preferences[TOKEN_KEY] ?: ""
    }

    suspend fun clearToken() {
        dataStore.edit { preferences ->
            preferences.remove(TOKEN_KEY)
        }
    }
}
```

**En MainActivity.kt**:
```kotlin
// Crear DataStore
val dataStore = createDataStore(
    name = "kiosco_prefs",
    serializer = PreferencesSerializer
)
val authManager = AuthManager(dataStore)

// Después de login exitoso
authManager.saveToken(cookieValue)

// Cuando necesites el token (para descargar BD)
lifecycleScope.launch {
    authManager.getToken().collect { token ->
        if (token.isNotEmpty()) {
            descargarBD(token)
        }
    }
}
```

**No es necesario** parsear el token en la app. Solo:
- Almacenarlo en DataStore
- Pasarlo en requests
- Validar que existe antes de descargar BD

---

## Permisos (AndroidManifest.xml)

```xml
<uses-permission android:name="android.permission.INTERNET" />
<uses-permission android:name="android.permission.READ_EXTERNAL_STORAGE" />
<uses-permission android:name="android.permission.WRITE_EXTERNAL_STORAGE" />
<!-- Para compartir por WhatsApp -->
```

**Nota**: Storage permisos para Android 6.0+ requieren runtime permissions.

---

## Dependencias (build.gradle.kts) — 2026 Versiones Actuales

```kotlin
dependencies {
    // Core Android
    implementation("androidx.core:core-ktx:1.13.0")
    implementation("androidx.appcompat:appcompat:1.7.0")
    implementation("androidx.constraintlayout:constraintlayout:2.2.0")

    // Material Design 3
    implementation("com.google.android.material:material:1.12.0")

    // Jetpack Lifecycle & Navigation (Actualizado: 2026)
    implementation("androidx.lifecycle:lifecycle-viewmodel-ktx:2.11.0")
    implementation("androidx.lifecycle:lifecycle-livedata-ktx:2.11.0")
    implementation("androidx.lifecycle:lifecycle-runtime-ktx:2.11.0")
    implementation("androidx.navigation:navigation-fragment-ktx:2.8.0")
    implementation("androidx.navigation:navigation-ui-ktx:2.8.0")

    // Room 3.0 (Nuevo en 2026)
    implementation("androidx.room:room-runtime:3.0.0")
    implementation("androidx.room:room-ktx:3.0.0")
    kapt("androidx.room:room-compiler:3.0.0")

    // DataStore (reemplaza EncryptedSharedPreferences deprecado)
    implementation("androidx.datastore:datastore-preferences:1.1.1")
    implementation("androidx.datastore:datastore-preferences-core:1.1.1")
    
    // Encryption para DataStore (Tink)
    implementation("androidx.security:security-crypto:1.1.0-alpha07")

    // Networking (OkHttp 5.x)
    implementation("com.squareup.retrofit2:retrofit:2.11.0")
    implementation("com.squareup.retrofit2:converter-gson:2.11.0")
    implementation("com.squareup.okhttp3:okhttp:5.2.0")
    implementation("com.squareup.okhttp3:logging-interceptor:5.2.0")

    // JSON
    implementation("com.google.code.gson:gson:2.11.0")

    // Graphics (para generar imagen de resumen)
    implementation("androidx.graphics:graphics-core:1.0.0-beta01")

    // Testing
    testImplementation("junit:junit:4.13.2")
    androidTestImplementation("androidx.test.ext:junit:1.1.5")
    androidTestImplementation("androidx.test.espresso:espresso-core:3.5.1")
}

plugins {
    id("com.android.application")
    kotlin("android")
    kotlin("kapt")
}
```

**Cambios importantes vs. versión anterior:**
- ✅ `EncryptedSharedPreferences` deprecado → Migrado a `DataStore` con Tink
- ✅ `Room 2.6.1` → `Room 3.0.0` (mejor soporte Kotlin)
- ✅ `OkHttp 4.11` → `OkHttp 5.2.0` (mejor performance)
- ✅ `Lifecycle 2.7.0` → `Lifecycle 2.11.0` (nuevas features)
- ✅ `Navigation 2.7.7` → `Navigation 2.8.0`

---

## Configuración del Proyecto

### 1. Crear proyecto en Android Studio
```bash
File > New > New Android Project
- Language: Kotlin
- Minimum SDK: API 24 (Android 7.0)
- Target SDK: API 34+ (Android 14+)
- Template: Empty Views Activity
```

### 2. Configurar build.gradle.kts
```kotlin
android {
    compileSdk = 34
    defaultConfig {
        minSdk = 24
        targetSdk = 34
        applicationId = "com.kiosco"
        versionCode = 1
        versionName = "1.0.0"
    }
    buildFeatures {
        viewBinding = true
    }
}
```

### 3. Crear constantes
Archivo: `utils/Constants.kt`
```kotlin
object Constants {
    const val BASE_URL = "http://192.168.1.100:3200/" // Cambiar a tu IP/dominio
    const val DB_NAME = "database.db"
    const val PREF_TOKEN = "kiosco_token"
    const val SHARED_PREF_NAME = "kiosco_prefs"
}
```

---

## Flujo de Usuario MVP

### Primera vez
1. **LoginActivity**: Email + password
2. **SincronizarFragment**: Descarga .db (muestra progress bar)
3. **MainActivity**: Tab de Consumos (por defecto)

### Uso diario
1. **ConsumosFragment**: Ver consumos, total semanal
2. **ResumenFragment**: Ver resumen visual
3. **Compartir por WhatsApp**: 
   - Click "Compartir"
   - Genera imagen
   - Abre Intent de WhatsApp
   - Usuario elige contacto/grupo

### Sincronizar
- Settings > "Actualizar datos"
- App re-descarga .db
- Reemplaza local
- Toast de confirmación

---

## Notas Técnicas Importantes

### Offline-First
- **No usar** Firestore/Cloud Sync. Solo SQLite local.
- Si usuario offline: app funciona 100% con datos descargados.
- Sincronización es **manual** (botón), no automática.

### Generación de Imagen (Resumen)
- Usar `Canvas` para dibujar en `Bitmap`
- Renderizar: nombre, total, desglose sector, fecha
- Guardar en `getCacheDir()` temporalmente
- Pasar a Intent.ACTION_SEND

### Navegación
- Usar **Navigation Component** (Jetpack)
- Bottom Navigation o Drawer para tabs
- Activities: LoginActivity (login) → MainActivity (resto)

### ViewModels y State
- Un ViewModel por pantalla (LoginViewModel, ConsumosViewModel, etc.)
- LiveData para observar cambios BD
- No usar static variables

### Testing
MVP no requiere tests automatizados. Validar manualmente:
1. Login con credenciales válidas/inválidas
2. Descarga .db (medir tiempo, verificar gzip)
3. Consultas locales (filtros, sumas)
4. Compartir por WhatsApp (imagen se ve bien)
5. Sincronizar (reemplaza datos)

---

## Roadmap Post-MVP

- [ ] **V1.1**: Caché de imágenes compartidas, historial
- [ ] **V1.2**: Sincronización incremental (solo deltas, no BD completa)
- [ ] **V1.3**: Push notifications (cambios en consumos)
- [ ] **V1.4**: Compartir por email/SMS (no solo WhatsApp)
- [ ] **V2.0**: Edición local de consumos + sincronización bidireccional

---

## Referencias de Código (Backend)

Cuando integres con el backend:

- **Login**: POST `/login` → cookie `kiosco_token`
- **Descargar BD**: GET `/api/database/download` (con cookie)
  - Response: `application/gzip` (archivo comprimido)
  - Requiere autenticación válida
  - Header: `Cookie: kiosco_token=...`

Schema BD: Ver `internal/config/schema.sql` en el repo kiosco

---

## Troubleshooting Rápido

| Problema | Causa | Solución |
|----------|-------|----------|
| "Autenticación fallida" | Token expirado | Re-login |
| "No se descarga .db" | Servidor offline | Verificar IP en Constants.BASE_URL |
| "BD corrompida después descarga" | Gzip mal descomprimido | Ver DatabaseManager.kt logs |
| "Room no inicializa" | .db no se copió | Verificar permisos de storage |
| "Imagen resumen se ve mal" | Canvas dimensions | Ajustar en ResumenGenerator.kt |

---

## Checklist MVP

- [ ] Proyecto creado en Android Studio Panda 4
- [ ] build.gradle.kts configurado (minSdk 24, Room, Retrofit)
- [ ] Entidades Room creadas (Estudiante, Consumo, Pago, Producto)
- [ ] DAOs y Database.kt implementados
- [ ] LoginActivity con POST /login
- [ ] AuthManager para guardar token en SharedPreferences
- [ ] EncryptedSharedPreferences configurado
- [ ] RetrofitClient + KioscoApiService para descargar .db
- [ ] DatabaseManager para descompresión e instalación
- [ ] MainActivity con Navigation Component
- [ ] ConsumosFragment + ConsumosViewModel
- [ ] PagosFragment + PagosViewModel
- [ ] ResumenFragment + generador de imagen (Canvas)
- [ ] Intent.ACTION_SEND para WhatsApp
- [ ] SincronizarFragment con botón de actualización
- [ ] Permisos en AndroidManifest.xml
- [ ] Testing manual en dispositivo/emulador
- [ ] APK generado y funcional

---

**Preguntas?** Referencia este prompt durante desarrollo. Está diseñado para MVP rápido sin overdoing architecture.

¡A codear! 🚀
