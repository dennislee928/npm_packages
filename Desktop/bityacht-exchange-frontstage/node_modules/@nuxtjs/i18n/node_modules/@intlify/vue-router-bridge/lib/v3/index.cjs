Object.defineProperty(exports, '__esModule', {
  value: true
})

var VueRouter = require('vue-router')
var VueDemi = require('vue-demi')

const useRouter = () => {
  const instance = VueDemi.getCurrentInstance()
  if (instance == null) {
    throw new Error(`should be used in setup`)
  }
  const vm = instance.proxy || instance
  return vm.$router
}

const useRoute = () => {
  const instance = VueDemi.getCurrentInstance()
  if (instance == null) {
    throw new Error(`should be used in setup`)
  }
  const vm = instance.proxy || instance
  return VueDemi.reactive(vm.$route)
}

VueRouter.useRouter = useRouter
VueRouter.useRoute = useRoute
VueRouter.isVueRouter3 = true
VueRouter.isVueRouter4 = false

/**
 * shim vue-router@4 typings
 */
const STUB = () => ({})
VueRouter.createRouter = STUB
VueRouter.createMemoryHistory = STUB
VueRouter.createRouterMatcher = STUB
VueRouter.createWebHashHistory = STUB
VueRouter.createWebHistory = STUB

exports.default = VueRouter
module.exports = VueRouter
