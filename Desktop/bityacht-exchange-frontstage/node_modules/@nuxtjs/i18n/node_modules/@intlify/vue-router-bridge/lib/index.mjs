import { useRouter, useRoute } from 'vue-router'

var isVueRouter3 = false
var isVueRouter4 = true

// stub vue-router 3 class
class VueRouterLegacy {
  static install() {}
  static version = ''
}

export * from 'vue-router'
export { useRouter, useRoute, isVueRouter3, isVueRouter4 }

export default VueRouterLegacy
