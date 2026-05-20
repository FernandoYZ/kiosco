// Kiosco Escolar - lógica del cliente

// ---------------------------------------------------------------------------
// Registro del Service Worker (PWA)
// ---------------------------------------------------------------------------
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js').catch((err) => {
      console.warn('Service Worker no pudo registrarse:', err);
    });
  });
}

// ---------------------------------------------------------------------------
// Copiar/descargar el comprobante como imagen PNG usando html-to-image
// ---------------------------------------------------------------------------
window.copiarComprobante = async function copiarComprobante(idElemento, nombreArchivo) {
  const elemento = document.getElementById(idElemento);
  const btnCopiar = document.getElementById("btnCopiar");
  const iconoCopiar = document.getElementById("iconoCopiar");
  const textoCopiar = document.getElementById("textoCopiar");

  if (!elemento) return;

  btnCopiar.disabled = true;
  textoCopiar.textContent = "Generando...";
  iconoCopiar.textContent = "⏳";

  try {
    const blob = await htmlToImage.toBlob(elemento, {
      backgroundColor: "#ffffff",
      pixelRatio: 2,
    });

    if (!blob) throw new Error("No se pudo generar la imagen");

    const filename = nombreArchivo || "comprobante.png";
    const file = new File([blob], filename, { type: "image/png" });

    // 1. Prioridad: Intentar usar Web Share API (Ideal para iPad/iOS)
    if (navigator.canShare && navigator.canShare({ files: [file] })) {
      await navigator.share({
        files: [file],
        title: "Comprobante",
        // IMPORTANTE: No enviamos las propiedades 'text' ni 'url' 
        // para evitar que WhatsApp adjunte texto basura.
      });
      textoCopiar.textContent = "¡Compartido!";
      iconoCopiar.textContent = "✅";
      btnCopiar.style.backgroundColor = "#16a34a";
    } 
    // 2. Fallback 1: Copiar al portapapeles (Ideal para Desktop)
    else if (navigator.clipboard?.write) {
      await navigator.clipboard.write([
        new ClipboardItem({ "image/png": blob }),
      ]);
      textoCopiar.textContent = "¡Copiado!";
      iconoCopiar.textContent = "✅";
      btnCopiar.style.backgroundColor = "#16a34a";
    } 
    // 3. Fallback 2: Descarga clásica (Navegadores antiguos o sin soporte)
    else {
      const url = URL.createObjectURL(blob);
      const enlace = document.createElement("a");
      enlace.href = url;
      enlace.download = filename;
      enlace.click();
      URL.revokeObjectURL(url);
      textoCopiar.textContent = "Descargado";
      iconoCopiar.textContent = "💾";
      btnCopiar.style.backgroundColor = "#16a34a";
    }
  } catch (error) {
    console.error("Error al generar imagen:", error);
    // Si el usuario cancela el menú de compartir en iOS, lanzará un error (AbortError).
    // Puedes evitar que se muestre como "Error" en rojo si fue solo una cancelación.
    if (error.name === "AbortError") {
      textoCopiar.textContent = "Cancelado";
      iconoCopiar.textContent = "ℹ️";
      btnCopiar.style.backgroundColor = "#f59e0b"; // Naranja/Warning
    } else {
      textoCopiar.textContent = "Error";
      iconoCopiar.textContent = "❌";
      btnCopiar.style.backgroundColor = "#dc2626";
    }
  } finally {
    setTimeout(() => {
      textoCopiar.textContent = "Copiar como imagen";
      iconoCopiar.textContent = "📋";
      btnCopiar.style.backgroundColor = "";
      btnCopiar.disabled = false;
    }, 2000);
  }
}

// ---------------------------------------------------------------------------
// Control de BottomNav: Mostrar solo en páginas autenticadas
// ---------------------------------------------------------------------------
document.addEventListener('DOMContentLoaded', () => {
  const bottomNav = document.getElementById('bottomNav');
  if (!bottomNav) return;

  // Esconder el bottom nav en páginas públicas como /login
  if (window.location.pathname === '/login') {
    bottomNav.style.display = 'none';
  }

  // Hide-on-Scroll: Esconde el Bottom Nav al hacer scroll hacia abajo
  let lastScrollTop = 0;
  const scrollThreshold = 50; // Píxeles de scroll antes de reaccionar

  window.addEventListener('scroll', () => {
    // Si está escondido, no hacer nada
    if (bottomNav.style.display === 'none') return;

    const currentScroll = window.pageYOffset || document.documentElement.scrollTop;

    // Scroll hacia abajo: esconde
    if (currentScroll > lastScrollTop && currentScroll > scrollThreshold) {
      bottomNav.style.transform = 'translateY(100%)';
    }
    // Scroll hacia arriba: muestra
    else if (currentScroll < lastScrollTop) {
      bottomNav.style.transform = 'translateY(0)';
    }

    lastScrollTop = currentScroll <= 0 ? 0 : currentScroll;
  }, false);
});

// ---------------------------------------------------------------------------
// Registro de Consumos - Lógica para grid
// ---------------------------------------------------------------------------
document.addEventListener('DOMContentLoaded', () => {
  inicializarRegistroConsumos();
});

function inicializarRegistroConsumos() {
  // Búsqueda y filtro por grado
  const searchInput = document.getElementById('search-input');
  if (searchInput) {
    // Generar botones de grados únicos
    const gradosSet = new Set();
    document.querySelectorAll('.estudiante-fila').forEach(fila => {
      const gradoP = fila.querySelector('p:last-of-type');
      if (gradoP) gradosSet.add(gradoP.textContent.trim());
    });

    const container = document.querySelector('.grado-btn')?.parentElement;
    if (container && gradosSet.size > 0) {
      for (const grado of Array.from(gradosSet).sort()) {
        if (grado === 'Todos') continue;
        const btn = document.createElement('button');
        btn.type = 'button';
        btn.className = 'grado-btn px-3 sm:px-4 py-1.5 sm:py-2 rounded-full text-[13px] sm:text-[14px] font-semibold transition-all whitespace-nowrap shadow-sm border bg-white border-gray-200 text-gray-600 active:scale-95';
        btn.dataset.grado = grado;
        btn.textContent = grado;
        container.appendChild(btn);
      }
    }

    // Estado del filtro
    let filtroState = { grado: 'Todos', search: '' };

    function aplicarFiltro() {
      const container = document.getElementById('estudiantes-container');
      if (!container) return;
      const filas = container.querySelectorAll('.estudiante-fila');
      let visibles = 0;

      filas.forEach(fila => {
        const apellidosNombres = fila.textContent.toLowerCase();
        const gradoP = fila.querySelector('p:last-of-type');
        const grado = gradoP ? gradoP.textContent.trim() : '';

        const matchGrado = filtroState.grado === 'Todos' || grado === filtroState.grado;
        const matchSearch = filtroState.search === '' || apellidosNombres.includes(filtroState.search.toLowerCase());

        if (matchGrado && matchSearch) {
          fila.style.display = '';
          visibles++;
        } else {
          fila.style.display = 'none';
        }
      });

      const contador = document.getElementById('contador');
      if (contador) contador.textContent = `encontrados: ${visibles}`;
    }

    // Evento de búsqueda
    searchInput.addEventListener('input', (e) => {
      filtroState.search = e.target.value;
      aplicarFiltro();
    });

    // Eventos de grados
    document.querySelectorAll('.grado-btn').forEach(btn => {
      btn.addEventListener('click', (e) => {
        e.preventDefault();
        document.querySelectorAll('.grado-btn').forEach(b => {
          b.classList.remove('bg-[#007AFF]', 'border-[#007AFF]', 'text-white');
          b.classList.add('bg-white', 'border-gray-200', 'text-gray-600');
        });
        e.target.classList.add('bg-[#007AFF]', 'border-[#007AFF]', 'text-white');
        e.target.classList.remove('bg-white', 'border-gray-200', 'text-gray-600');
        filtroState.grado = e.target.dataset.grado;
        aplicarFiltro();
      });
    });
  }

  // Cambio de fecha
  const selectorFecha = document.getElementById('selector-fecha');
  if (selectorFecha) {
    selectorFecha.addEventListener('change', (e) => {
      const pathParts = window.location.pathname.split('/');
      const base = pathParts[1]; // 'registro' | 'resumen'
      const sector = pathParts[2];
      if (sector === 'menor' || sector === 'mayor') {
        window.location = `/${base}/${sector}?fecha=${e.target.value}`;
      }
    });
  }
}
