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
