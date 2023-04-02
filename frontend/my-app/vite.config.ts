import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/upload": {
        target: "https://team24-f6kqjvqzqa-uc.a.run.app",
        changeOrigin: true,
        secure: false,
      },
    },
  },
});
