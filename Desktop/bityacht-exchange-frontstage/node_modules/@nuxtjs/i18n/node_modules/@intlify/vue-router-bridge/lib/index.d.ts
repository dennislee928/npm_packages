/* eslint-disable @typescript-eslint/no-explicit-any */
import type { RouteLocationNormalizedLoaded, Router, RouteRecord, RouteMeta, NavigationGuard } from 'vue-router'

/**
 * shim vue-router@3 typings
 */

// Stub VueRouter class interfaces
declare interface VueRouter {
  app: Vue
  options: RouterOptions
  mode: RouterMode
  currentRoute: Route
  beforeEach(guard: NavigationGuard): Function
  beforeResolve(guard: NavigationGuard): Function
  afterEach(hook: (to: Route, from: Route) => any): Function
  push(location: RawLocation): Promise<Route>
  push(location: RawLocation, onComplete?: Function, onAbort?: ErrorHandler): void
  replace(location: RawLocation): Promise<Route>
  replace(location: RawLocation, onComplete?: Function, onAbort?: ErrorHandler): void
  go(n: number): void
  back(): void
  forward(): void
  match(raw: RawLocation, current?: Route, redirectedFrom?: Location): Route
  getMatchedComponents(to?: RawLocation | Route): Component[]
  onReady(cb: Function, errorCb?: ErrorHandler): void
  onError(cb: ErrorHandler): void
  addRoutes(routes: RouteConfig[]): void
  AddRoute(parent: string, route: RouteConfig): void
  AddRoute(route: RouteConfig): void
  resolve(
    to: RawLocation,
    current?: Route,
    append?: boolean
  ): {
    location: Location
    route: Route
    href: string
    normalizedTo: Location
    resolved: Route
  }
}
declare type ErrorHandler = (err: Error) => void
declare type RouterMode = 'hash' | 'history' | 'abstract'
declare type Dictionary<T> = { [key: string]: T }
declare interface Location {
  name?: string
  path?: string
  hash?: string
  query?: Dictionary<string | (string | null)[] | null | undefined>
  params?: Dictionary<string>
  append?: boolean
  replace?: boolean
}
declare type RawLocation = string | Location
declare interface Route {
  path: string
  name?: string | null
  hash: string
  query: Dictionary<string | (string | null)[]>
  params: Dictionary<string>
  fullPath: string
  matched: RouteRecord[]
  redirectedFrom?: string
  meta?: RouteMeta
}
declare type RedirectOption = RawLocation | ((to: Route) => RawLocation)
declare interface PathToRegexpOptions {
  sensitive?: boolean
  strict?: boolean
  end?: boolean
}
declare type Vue = any // emulate Vue (v2)
declare type Component = any // emulate Vue Component (v2)
declare type Position = { x: number; y: number }
declare type PositionResult = Position | { selector: string; offset?: Position; behavior?: ScrollBehavior } | void
declare interface RouterOptions {
  routes?: RouteConfig[]
  mode?: RouterMode
  fallback?: boolean
  base?: string
  linkActiveClass?: string
  linkExactActiveClass?: string
  parseQuery?: (query: string) => Record<string, any>
  stringifyQuery?: (query: Record<string, any>) => string
  scrollBehavior?: (
    to: Route,
    from: Route,
    savedPosition: Position | void
  ) => PositionResult | Promise<PositionResult> | undefined | null
}
declare type RoutePropsFunction = (route: Route) => Record<string, any>
declare interface _RouteConfigBase {
  path: string
  name?: string
  children?: RouteConfig[]
  redirect?: RedirectOption
  alias?: string | string[]
  meta?: RouteMeta
  beforeEnter?: NavigationGuard
  caseSensitive?: boolean
  pathToRegexpOptions?: PathToRegexpOptions
}
declare interface RouteConfigSingleView extends _RouteConfigBase {
  component?: Component
  props?: boolean | Record<string, any> | RoutePropsFunction
}
declare interface RouteConfigMultipleViews extends _RouteConfigBase {
  components?: Dictionary<Component>
  props?: Dictionary<boolean | Record<string, any> | RoutePropsFunction>
}
declare type RouteConfig = RouteConfigSingleView | RouteConfigMultipleViews
declare interface RouteRecordPublic {
  path: string
  components: Dictionary<Component>
  instances: Dictionary<Vue>
  name?: string
  redirect?: RedirectOption
  meta: any
  beforeEnter?: (route: Route, redirect: (location: RawLocation) => void, next: () => void) => any
  props:
    | boolean
    | Record<string, any>
    | RoutePropsFunction
    | Dictionary<boolean | Record<string, any> | RoutePropsFunction>
}

/**
 * Returns the router instance. Equivalent to using `$router` inside templates.
 */
declare function useRouter<T = Router>(): T
/**
 * Returns the current route location. Equivalent to using `$route` inside templates.
 */
declare function useRoute<T = RouteLocationNormalizedLoaded>(): T

/**
 * Wheter the current vue-router version is 3
 */
declare const isVueRouter3: boolean
/**
 * Wheter the current vue-router version is 4
 */
declare const isVueRouter4: boolean

/**
 * Define stub vue-router 3 class
 */
declare class VueRouterLegacy {
  constructor(options?: RouterOptions)
}

// export vue-router@4 typings
export * from 'vue-router'

export {
  VueRouter,
  useRouter,
  useRoute,
  isVueRouter3,
  isVueRouter4,
  RouterMode,
  RawLocation,
  RedirectOption,
  RouteConfig,
  RouteRecordPublic,
  Location,
  Route
}

// Export dummy VueRouter class
export default VueRouterLegacy

/* eslint-enable @typescript-eslint/no-explicit-any */
