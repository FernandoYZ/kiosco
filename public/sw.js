const CACHE = 'kiosco-v1';

const ASSETS_ESTATICOS = [
  '/dist/styles.css',
  '/dist/htmx.min.js',
  '/dist/alpine.min.js',
  '/dist/bundle.min.js',
  '/dist/canvas.min.js',
  '/favicon.webp',
];

// Instalar: cachear assets estáticos
self.addEventListener('install', (e) => {
  e.waitUntil(
    caches.open(CACHE).then((cache) => cache.addAll(ASSETS_ESTATICOS))
  );
  self.skipWaiting();
});

// Activar: limpiar cachés anteriores
self.addEventListener('activate', (e) => {
  e.waitUntil(
    caches.keys().then((keys) =>
      Promise.all(
        keys.filter((k) => k !== CACHE).map((k) => caches.delete(k))
      )
    )
  );
  self.clients.claim();
});

// Fetch: cache-first para estáticos
self.addEventListener('fetch', (e) => {
  // Solo manejar GET
  if (e.request.method !== 'GET') return;

  const url = new URL(e.request.url);

  // Assets estáticos
  const esEstatico =
    url.pathname.startsWith('/dist/') ||
    url.pathname.startsWith('/fonts/') ||
    url.pathname === '/favicon.webp' ||
    url.pathname === '/manifest.json';

  if (esEstatico) {
    e.respondWith(
      caches.match(e.request).then((cached) => {
        const network = fetch(e.request).then((res) => {
          // Solo cachear si el scheme es http/https (no chrome-extension, etc)
          if (res.ok && (url.protocol === 'http:' || url.protocol === 'https:')) {
            const clone = res.clone();
            caches.open(CACHE).then((cache) => cache.put(e.request, clone));
          }
          return res;
        });
        return cached || network;
      })
    );
    return;
  }

  // Páginas HTML
  e.respondWith(
    fetch(e.request).catch(() => caches.match(e.request))
  );
});
