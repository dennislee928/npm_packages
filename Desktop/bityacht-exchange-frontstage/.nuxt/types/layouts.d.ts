import { ComputedRef, Ref } from 'vue'
export type LayoutKey = "default" | "members"
declare module "/Users/lipeichen/Desktop/彼雅特github/bityacht-exchange-frontstage/node_modules/nuxt/dist/pages/runtime/composables" {
  interface PageMeta {
    layout?: false | LayoutKey | Ref<LayoutKey> | ComputedRef<LayoutKey>
  }
}