import { defineConfig } from "@rspress/core";
import { pluginSitemap } from "@rspress/plugin-sitemap";
import path from "node:path";
import pluginKatex from "rspress-plugin-katex";

export default defineConfig({
  root: "content",
  base: "/esquema/",
  lang: "en",
  title: "esquema",
  description:
    "A declarative runtime that turns HCL manifests into HTTP APIs.",
  logoText: "esquema",
  outDir: "dist",
  globalStyles: path.join(__dirname, "theme/index.css"),
  head: [["meta", { name: "theme-color", content: "#d97706" }]],
  plugins: [
    pluginSitemap({
      siteUrl: "https://ju4n97.github.io/esquema/",
    }),
    pluginKatex(),
  ],
  llms: true,
  themeConfig: {
    enableContentAnimation: false,
    enableAppearanceAnimation: false,
    lastUpdated: true,
    enableScrollToTop: true,
    llmsUI: true,
    editLink: {
      docRepoBaseUrl: "https://github.com/ju4n97/esquema/edit/main/",
    },
    socialLinks: [
      {
        icon: "github",
        mode: "link",
        content: "https://github.com/ju4n97/esquema",
      },
    ],
    footer: {
      message:
        '<a href="https://github.com/ju4n97/esquema/blob/main/LICENSE" target="_blank" rel="noreferrer">MIT License</a> © 2026 esquema contributors.',
    },
  },
  route: {
    useTransitions: false,
  },
});