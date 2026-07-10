// @ts-check
import {defineConfig, fontProviders} from 'astro/config';

export default defineConfig({
  base: "/lumingrate",
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
