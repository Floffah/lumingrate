// @ts-check
import {defineConfig, fontProviders} from 'astro/config';

import tailwindcss from '@tailwindcss/vite';
import wasm from "vite-plugin-wasm";

// https://astro.build/config
export default defineConfig({
  base: "/lumingrate",
  vite: {
    plugins: [tailwindcss()]
  },
  fonts: [{
    provider: fontProviders.local(),
    name: "Monaspace Neon Var",
    cssVariable: "--font-mono",

    options: {
      variants: [{
        src: ['./src/fonts/Monaspace Neon Var.woff2'],
        weight: '1 1000',
      }]
    }
  }]
});