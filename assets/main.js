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
// Registro de Consumos - Lógica de Alpine.js
// ---------------------------------------------------------------------------
window.registroApp = function() {
  const app = {
    view: 'sectores',
    sector: null,
    cargando: false,
    search: '',
    gradoActivo: 'Todos',
    estudiantes: [],
    productos: [],
    openId: null,
    bandeja: {},

    async seleccionarSector(s) {
      this.cargando = true;
      this.sector = s;
      try {
        const res = await fetch(`/registro/estudiantes?sector=${s}`);
        const data = await res.json();
        this.estudiantes = data.estudiantes || [];
        this.productos = data.productos || [];
        this.view = 'registro';
        this.bandeja = {};
        this.openId = null;
        this.search = '';
        this.gradoActivo = 'Todos';
      } catch (e) {
        console.error('Error cargando sector:', e);
        alert('Error al cargar los datos');
      } finally {
        this.cargando = false;
      }
    },

    volverASectores() {
      this.view = 'sectores';
      this.sector = null;
      this.estudiantes = [];
      this.productos = [];
      this.bandeja = {};
      this.openId = null;
      this.search = '';
      this.gradoActivo = 'Todos';
    },

    toggle(id) {
      this.openId = this.openId === id ? null : id;
    },

    agregarABandeja(estId, prod) {
      if (!this.bandeja[estId]) this.bandeja[estId] = [];
      const item = this.bandeja[estId].find(p => p.id === prod.IdProducto);
      if (item) {
        item.qty++;
      } else {
        this.bandeja[estId].push({
          id: prod.IdProducto,
          nombre: prod.Nombre,
          precio: prod.PrecioUnitario,
          qty: 1
        });
      }
    },

    getQty(estId, prodId) {
      return this.bandeja[estId]?.find(p => p.id === prodId)?.qty || 0;
    },

    limpiarBandeja(estId) {
      this.bandeja[estId] = [];
    },

    calcularTotal(estId) {
      const items = this.bandeja[estId] || [];
      return items.reduce((a, b) => a + (b.precio * b.qty), 0);
    },

    async confirmar(estId) {
      const items = this.bandeja[estId] || [];
      if (!items.length) return;

      const payload = {
        items: items.map(p => ({
          id_estudiante: estId,
          id_producto: p.id,
          cantidad: p.qty
        })),
        fecha: new Date().toISOString().split('T')[0]
      };

      try {
        const res = await fetch('/registro/guardar', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        });

        if (!res.ok) throw new Error('Error en la respuesta');

        // Éxito: limpiar y cerrar
        this.limpiarBandeja(estId);
        this.openId = null;
        alert('Consumos registrados correctamente');
      } catch (e) {
        console.error('Error guardando:', e);
        alert('Error al registrar consumos');
      }
    },

    filtrarEstudiantes() {
      return this.estudiantes.filter(e => {
        const coincideBusqueda = e.Apellidos.includes(this.search) ||
          e.Nombres.includes(this.search);
        const coincideGrado = this.gradoActivo === 'Todos' || e.NombreGrado === this.gradoActivo;
        return coincideBusqueda && coincideGrado;
      });
    },

    gradosDisponibles() {
      const grados = new Set();
      this.estudiantes.forEach(e => {
        if (e.NombreGrado) grados.add(e.NombreGrado);
      });
      return Array.from(grados).sort();
    }
  };
  return app;
};
