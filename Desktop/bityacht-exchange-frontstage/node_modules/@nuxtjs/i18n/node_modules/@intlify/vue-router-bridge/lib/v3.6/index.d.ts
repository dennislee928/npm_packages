/* eslint-disable @typescript-eslint/no-explicit-any */
import VueRouter from 'vue-router'

import type { ComponentPublicInstance } from 'vue-demi'
import type { Route, RouterOptions } from 'vue-router'
import type {
  useRoute,
  useRouter,
  useLink,
  onBeforeRouteUpdate,
  onBeforeRouteLeave,
  RouterLinkOptions
} from 'vue-router/composables'

/**
 * shim vue-router@4 typings
 */

declare type PathParserOptions = Pick<_PathParserOptions, 'end' | 'sensitive' | 'strict'>

declare interface _PathParserOptions {
  /**
   * Makes the RegExp case sensitive. Defaults to false
   */
  sensitive?: boolean
  /**
   * Should we disallow a trailing slash. Defaults to false
   */
  strict?: boolean
  /**
   * Should the RegExp match from the beginning by prepending a `^` to it. Defaults to true
   * @internal
   */
  start?: boolean
  /**
   * Should the RegExp match until the end by appending a `$` to it. Defaults to true
   */
  end?: boolean
}

declare type Lazy<T> = () => Promise<T>

/**
 * Allowed Component definitions in route records provided by the user
 */
declare type RawRouteComponent = RouteComponent | Lazy<RouteComponent>

declare type RouteComponent = /* Component */ any

/**
 * Possible values in normalized {@link LocationQuery}. `null` renders the query
 * param but without an `=`.
 *
 * @example
 * ```
 * ?isNull&isEmpty=&other=other
 * gives
 * `{ isNull: null, isEmpty: '', other: 'other' }`.
 * ```
 */
declare type LocationQueryValue = string | null

/**
 * Normalized query object that appears in {@link RouteLocationNormalized}
 */
declare type LocationQuery = Record<string, LocationQueryValue | LocationQueryValue[]>

/**
 * Possible values when defining a query.
 */
declare type LocationQueryValueRaw = LocationQueryValue | number | undefined

declare type LocationQueryRaw = Record<string | number, LocationQueryValueRaw | LocationQueryValueRaw[]>

declare interface LocationAsPath {
  path: string
}

declare interface LocationAsRelativeRaw {
  name?: RouteRecordName
  params?: RouteParamsRaw
}

declare interface RouteQueryAndHash {
  query?: LocationQueryRaw
  hash?: string
}

declare interface RouteRecordSingleView extends _RouteRecordBase {
  /**
   * Component to display when the URL matches this route.
   */
  component: RawRouteComponent
  components?: never
  /**
   * Allow passing down params as props to the component rendered by `router-view`.
   */
  props?: _RouteRecordProps
}

declare interface RouteRecordMultipleViews extends _RouteRecordBase {
  /**
   * Components to display when the URL matches this route. Allow using named views.
   */
  components: Record<string, RawRouteComponent>
  component?: never
  /**
   * Allow passing down params as props to the component rendered by
   * `router-view`. Should be an object with the same keys as `components` or a
   * boolean to be applied to every component.
   */
  props?: Record<string, _RouteRecordProps> | boolean
}

/**
 * Route Record that defines a redirect. Cannot have `component` or `components`
 * as it is never rendered.
 */
declare interface RouteRecordRedirect extends _RouteRecordBase {
  redirect: RouteRecordRedirectOption
  component?: never
  components?: never
}

declare type _RouteRecordProps = boolean | Record<string, any> | ((to: RouteLocationNormalized) => Record<string, any>)

declare type RouteRecordRaw = RouteRecordSingleView | RouteRecordMultipleViews | RouteRecordRedirect

/**
 * Possible values for a user-defined route record's name
 */
declare type RouteRecordName = string | symbol

declare type RouteParamValue = string

declare type RouteParamValueRaw = RouteParamValue | number | null | undefined

declare type RouteParams = Record<string, RouteParamValue | RouteParamValue[]>

declare type RouteParamsRaw = Record<string, RouteParamValueRaw | Exclude<RouteParamValueRaw, null | undefined>[]>

declare type RouteRecord = RouteRecordNormalized

type HistoryStateArray = Array<HistoryStateValue>

declare type HistoryStateValue = string | number | boolean | null | undefined | HistoryState | HistoryStateArray

/**
 * Allowed HTML history.state
 */
declare interface HistoryState {
  [x: number]: HistoryStateValue
  [x: string]: HistoryStateValue
}

declare interface RouteLocationOptions {
  /**
   * Replace the entry in the history instead of pushing a new entry
   */
  replace?: boolean
  /**
   * Triggers the navigation even if the location is the same as the current one
   */
  force?: boolean
  /**
   * State to save using the History API. This cannot contain any reactive
   * values and some primitives like Symbols are forbidden. More info at
   * https://developer.mozilla.org/en-US/docs/Web/API/History/state
   */
  state?: HistoryState
}

/**
 * Interface to type `meta` fields in route records.
 *
 * @example
 *
 * ```ts
 * // typings.d.ts or router.ts
 * import 'vue-router';
 *
 * declare module 'vue-router' {
 *   interface RouteMeta {
 *     requiresAuth?: boolean
 *   }
 *  }
 * ```
 */
declare type RouteMeta = Record<string | number | symbol, unknown>

/**
 * User-level route location
 */
declare type RouteLocationRaw =
  | string
  | (RouteQueryAndHash & LocationAsPath & RouteLocationOptions)
  | (RouteQueryAndHash & LocationAsRelativeRaw & RouteLocationOptions)

declare type RouteRecordRedirectOption = RouteLocationRaw | ((to: RouteLocation) => RouteLocationRaw)

declare interface RouteLocation extends _RouteLocationBase {
  /**
   * Array of {@link RouteRecord} containing components as they were
   * passed when adding records. It can also contain redirect records. This
   * can't be used directly
   */
  matched: RouteRecord[]
}

declare interface _RouteLocationBase {
  /**
   * Percentage encoded pathname section of the URL.
   */
  path: string
  /**
   * The whole location including the `search` and `hash`. This string is
   * percentage encoded.
   */
  fullPath: string
  /**
   * Object representation of the `search` property of the current location.
   */
  query: LocationQuery
  /**
   * Hash of the current location. If present, starts with a `#`.
   */
  hash: string
  /**
   * Name of the matched record
   */
  name: RouteRecordName | null | undefined
  /**
   * Object of decoded params extracted from the `path`.
   */
  params: RouteParams
  /**
   * Contains the location we were initially trying to access before ending up
   * on the current location.
   */
  redirectedFrom: RouteLocation | undefined
  /**
   * Merged `meta` properties from all of the matched route records.
   */
  meta: RouteMeta
}

declare interface _RouteRecordBase extends PathParserOptions {
  /**
   * Path of the record. Should start with `/` unless the record is the child of
   * another record.
   *
   * @example `/users/:id` matches `/users/1` as well as `/users/posva`.
   */
  path: string
  /**
   * Where to redirect if the route is directly matched. The redirection happens
   * before any navigation guard and triggers a new navigation with the new
   * target location.
   */
  redirect?: RouteRecordRedirectOption
  /**
   * Array of nested routes.
   */
  children?: RouteRecordRaw[]
  /**
   * Aliases for the record. Allows defining extra paths that will behave like a
   * copy of the record. Allows having paths shorthands like `/users/:id` and
   * `/u/:id`. All `alias` and `path` values must share the same params.
   */
  alias?: string | string[]
  /**
   * Name for the route record.
   */
  name?: RouteRecordName
  /**
   * Before Enter guard specific to this record. Note `beforeEnter` has no
   * effect if the record has a `redirect` property.
   */
  beforeEnter?: NavigationGuardWithThis<undefined> | NavigationGuardWithThis<undefined>[]
  /**
   * Arbitrary data attached to the record.
   */
  meta?: RouteMeta
}

declare type NavigationGuardNextCallback = (vm: ComponentPublicInstance) => any

declare type NavigationGuardReturn = void | Error | RouteLocationRaw | boolean | NavigationGuardNextCallback

export declare interface NavigationGuard {
  (to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext):
    | NavigationGuardReturn
    | Promise<NavigationGuardReturn>
}

declare interface NavigationGuardNext {
  (): void
  (error: Error): void
  (location: RouteLocationRaw): void
  (valid: boolean | undefined): void
  (cb: NavigationGuardNextCallback): void
}

declare interface NavigationGuardWithThis<T> {
  (this: T, to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext):
    | NavigationGuardReturn
    | Promise<NavigationGuardReturn>
}

export declare interface RouteLocationNormalized extends _RouteLocationBase {
  /**
   * Array of {@link RouteRecordNormalized}
   */
  matched: RouteRecordNormalized[]
}

declare interface RouteLocationMatched extends RouteRecordNormalized {
  components: Record<string, RouteComponent>
}
declare interface RouteRecordNormalized {
  path: _RouteRecordBase['path']
  redirect: _RouteRecordBase['redirect'] | undefined
  name: _RouteRecordBase['name']
  components: RouteRecordMultipleViews['components']
  children: Exclude<_RouteRecordBase['children'], void>
  meta: Exclude<_RouteRecordBase['meta'], void>
  props: Record<string, _RouteRecordProps>
  beforeEnter: _RouteRecordBase['beforeEnter']
  leaveGuards: Set<NavigationGuard>
  updateGuards: Set<NavigationGuard>
  enterCallbacks: Record<string, NavigationGuardNextCallback[]>
  instances: Record<string, ComponentPublicInstance | undefined | null>
  aliasOf: RouteRecordNormalized | undefined
}

declare interface RouteLocationNormalizedLoaded extends _RouteLocationBase {
  /**
   * Array of {@link RouteLocationMatched} containing only plain components (any
   * lazy-loaded components have been loaded and were replaced inside of the
   * `components` object) so it can be directly used to display routes. It
   * cannot contain redirect records either
   */
  matched: RouteLocationMatched[]
}

declare type Router = any
declare type RouterHistory = any
declare type RouterMatcher = any

/**
 * Creates a Router instance that can be used by a Vue app.
 *
 * @param options - {@link RouterOptions}
 */
declare function createRouter(options: RouterOptions): Router

/**
 * Creates a in-memory based history. The main purpose of this history is to handle SSR. It starts in a special location that is nowhere.
 * It's up to the user to replace that location with the starter location by either calling `router.push` or `router.replace`.
 *
 * @param base - Base applied to all urls, defaults to '/'
 * @returns a history object that can be passed to the router constructor
 */
declare function createMemoryHistory(base?: string): RouterHistory

/**
 * Creates a Router Matcher.
 *
 * @internal
 * @param routes - array of initial routes
 * @param globalOptions - global route options
 */
declare function createRouterMatcher(routes: RouteRecordRaw[], globalOptions: PathParserOptions): RouterMatcher

/**
 * Creates a hash history. Useful for web applications with no host (e.g.
 * `file://`) or when configuring a server to handle any URL is not possible.
 *
 * @param base - optional base to provide. Defaults to `location.pathname +
 * location.search` If there is a `<base>` tag in the `head`, its value will be
 * ignored in favor of this parameter **but note it affects all the
 * history.pushState() calls**, meaning that if you use a `<base>` tag, it's
 * `href` value **has to match this parameter** (ignoring anything after the
 * `#`).
 *
 * @example
 * ```js
 * // at https://example.com/folder
 * createWebHashHistory() // gives a url of `https://example.com/folder#`
 * createWebHashHistory('/folder/') // gives a url of `https://example.com/folder/#`
 * // if the `#` is provided in the base, it won't be added by `createWebHashHistory`
 * createWebHashHistory('/folder/#/app/') // gives a url of `https://example.com/folder/#/app/`
 * // you should avoid doing this because it changes the original url and breaks copying urls
 * createWebHashHistory('/other-folder/') // gives a url of `https://example.com/other-folder/#`
 *
 * // at file:///usr/etc/folder/index.html
 * // for locations with no `host`, the base is ignored
 * createWebHashHistory('/iAmIgnored') // gives a url of `file:///usr/etc/folder/index.html#`
 * ```
 */
declare function createWebHashHistory(base?: string): RouterHistory

/**
 * Creates an HTML5 history. Most common history for single page applications.
 *
 * @param base -
 */
declare function createWebHistory(base?: string): RouterHistory

/**
 * Wheter the current vue-router version is 3
 */
declare const isVueRouter3: boolean
/**
 * Wheter the current vue-router version is 4
 */
declare const isVueRouter4: boolean

// export vue-router@4 typings
export * from 'vue-router'

export {
  Route,
  VueRouter,
  useRouter,
  useRoute,
  useLink,
  onBeforeRouteUpdate,
  onBeforeRouteLeave,
  RouterLinkOptions,
  isVueRouter3,
  isVueRouter4,
  createRouter,
  createMemoryHistory,
  createRouterMatcher,
  createWebHashHistory,
  createWebHistory,
  RouteLocationNormalizedLoaded
}

export default VueRouter

/* eslint-enable @typescript-eslint/no-explicit-any */
