import { defineConfig } from 'vitepress'
import { withMermaid } from 'vitepress-plugin-mermaid'

export default withMermaid(defineConfig({
  title: 'auditlog',
  description: 'A generic, type-safe audit-log domain for Go services',
  base: '/auditlog/',
  cleanUrls: true,
  markdown: {
    lineNumbers: true,
  },
  ignoreDeadLinks: [
    /^http:\/\/localhost/,
  ],
  themeConfig: {
    logo: '/logo.png',
    nav: [
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Reference', link: '/reference/configuration' },
      { text: 'godoc', link: 'https://pkg.go.dev/github.com/foomo/auditlog' },
    ],
    sidebar: {
      '/guide/': [
        {
          text: 'Guide',
          items: [
            { text: 'Getting started', link: '/guide/getting-started' },
            { text: 'Architecture', link: '/guide/architecture' },
            { text: 'Retention', link: '/guide/retention' },
            { text: 'gotsrpc integration', link: '/guide/gotsrpc' },
          ],
        },
      ],
      '/reference/': [
        {
          text: 'Reference',
          items: [
            { text: 'Configuration', link: '/reference/configuration' },
            { text: 'API', link: '/reference/api' },
          ],
        },
      ],
      '/': [
        {
          text: 'Contributing',
          items: [
            { text: 'Contributing', link: '/CONTRIBUTING' },
            { text: 'Code of Conduct', link: '/CODE_OF_CONDUCT' },
            { text: 'Security', link: '/SECURITY' },
          ],
        },
      ],
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/foomo/auditlog' },
    ],
    editLink: {
      pattern: 'https://github.com/foomo/auditlog/edit/main/docs/:path',
      text: 'Edit this page on GitHub',
    },
    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Made with ♥ by foomo / bestbytes',
    },
  },
  mermaid: {},
  mermaidPlugin: {
    class: 'mermaid',
  },
}))
