import VueRouter from 'vue-router'
import { useRoute, useRouter, useLink, onBeforeRouteUpdate, onBeforeRouteLeave } from 'vue-router/composables'

var isVueRouter3 = true
var isVueRouter4 = false

/**
 * shim vue-router@4 typings
 */
const STUB = () => ({})
const createRouter = STUB
const createMemoryHistory = STUB
const createRouterMatcher = STUB
const createWebHashHistory = STUB
const createWebHistory = STUB

export {
  useRouter,
  useRoute,
  useLink,
  onBeforeRouteUpdate,
  onBeforeRouteLeave,
  isVueRouter3,
  isVueRouter4,
  createRouter,
  createMemoryHistory,
  createRouterMatcher,
  createWebHashHistory,
  createWebHistory
}
export default VueRouter
