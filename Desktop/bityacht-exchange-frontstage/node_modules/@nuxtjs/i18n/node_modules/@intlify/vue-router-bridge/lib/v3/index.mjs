import VueRouter from 'vue-router'
import { getCurrentInstance, reactive } from 'vue-demi'

var isVueRouter3 = true
var isVueRouter4 = false

function useRouter() {
  const instance = getCurrentInstance()
  if (instance == null) {
    throw new Error(`should be used in setup`)
  }
  const vm = instance.proxy || instance
  return vm.$router
}

function useRoute() {
  const instance = getCurrentInstance()
  if (instance == null) {
    throw new Error(`should be used in setup`)
  }
  const vm = instance.proxy || instance
  return reactive(vm.$route)
}

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
  isVueRouter3,
  isVueRouter4,
  createRouter,
  createMemoryHistory,
  createRouterMatcher,
  createWebHashHistory,
  createWebHistory
}
export default VueRouter
