// Kiosco Escolar - lógica del cliente

// ---------------------------------------------------------------------------
// Alpine.data('comprobante') — copiar/descargar el comprobante como imagen PNG
// ---------------------------------------------------------------------------
document.addEventListener('alpine:init', () => {
  Alpine.data('comprobante', () => ({
    estado: '',
    isLoading: false,

    async copiar(idElemento, nombreArchivo) {
      if (this.isLoading) return;

      const elemento = document.getElementById(idElemento);
      if (!elemento) return;

      this.isLoading = true;
      this.estado = 'Generando...';

      try {
        const blob = await htmlToImage.toBlob(elemento, {
          backgroundColor: '#ffffff',
          pixelRatio: 2,
        });

        if (!blob) throw new Error('No se pudo generar la imagen');

        const filename = nombreArchivo || 'comprobante.png';
        const file = new File([blob], filename, { type: 'image/png' });

        // 1. Prioridad: Web Share API (Ideal para iPad/iOS)
        if (navigator.canShare && navigator.canShare({ files: [file] })) {
          await navigator.share({
            files: [file],
            title: 'Comprobante',
            // IMPORTANTE: No enviamos 'text' ni 'url' para evitar que
            // WhatsApp adjunte texto basura.
          });
          this.estado = '¡Compartido!';
          this.$el.style.backgroundColor = '#16a34a';
        }
        // 2. Fallback 1: Copiar al portapapeles (Desktop)
        else if (navigator.clipboard?.write) {
          await navigator.clipboard.write([
            new ClipboardItem({ 'image/png': blob }),
          ]);
          this.estado = '¡Copiado!';
          this.$el.style.backgroundColor = '#16a34a';
        }
        // 3. Fallback 2: Descarga clásica (Navegadores sin soporte)
        else {
          const url = URL.createObjectURL(blob);
          const enlace = document.createElement('a');
          enlace.href = url;
          enlace.download = filename;
          enlace.click();
          URL.revokeObjectURL(url);
          this.estado = 'Descargado 💾';
          this.$el.style.backgroundColor = '#16a34a';
        }
      } catch (error) {
        console.error('Error al generar imagen:', error);
        // Si el usuario cancela el menú de compartir en iOS → AbortError.
        // No se muestra como error rojo, sino como cancelación amber.
        if (error.name === 'AbortError') {
          this.estado = 'Cancelado';
          this.$el.style.backgroundColor = '#f59e0b';
        } else {
          this.estado = 'Error';
          this.$el.style.backgroundColor = '#dc2626';
        }
      } finally {
        setTimeout(() => {
          this.estado = '';
          this.$el.style.backgroundColor = '';
          this.isLoading = false;
        }, 2000);
      }
    },
  }));
});
