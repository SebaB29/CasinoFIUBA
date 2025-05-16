import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,     // aquí indicás que Vite corra en el puerto 3000
    host: true      // para que esté accesible desde fuera del contenedor
  }
})
