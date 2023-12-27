import type { NavigationGuard } from 'vue-router'
export type MiddlewareKey = string
declare module "/Users/lipeichen/Desktop/彼雅特github/bityacht-exchange-frontstage/node_modules/nuxt/dist/pages/runtime/composables" {
  interface PageMeta {
    middleware?: MiddlewareKey | NavigationGuard | Array<MiddlewareKey | NavigationGuard>
  }
}