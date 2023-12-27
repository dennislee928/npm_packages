Object.defineProperty(exports, '__esModule', {
  value: true
})

var VueRouter = require('vue-router')
var { useRoute, useRouter, useLink, onBeforeRouteUpdate, onBeforeRouteLeave } = require('vue-router/composables')

VueRouter.useRouter = useRouter
VueRouter.useRoute = useRoute
VueRouter.useLink = useLink
VueRouter.onBeforeRouteUpdate = onBeforeRouteUpdate
VueRouter.onBeforeRouteLeave = onBeforeRouteLeave
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
