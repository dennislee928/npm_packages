import { getRouteBaseName, localePath, localeRoute, switchLocalePath, localeHead } from 'vue-i18n-routing';
type GetRouteBaseName = typeof getRouteBaseName;
type LocalePath = typeof localePath;
type LocaleRoute = typeof localeRoute;
type LocaleHead = typeof localeHead;
type SwitchLocalePath = typeof switchLocalePath;
declare const _default: any;
export default _default;
declare module '#app' {
    interface NuxtApp {
        /**
         * Returns base name of current (if argument not provided) or passed in route.
         *
         * @remarks
         * Base name is name of the route without locale suffix and other metadata added by nuxt i18n module
         *
         * @param givenRoute - A route.
         *
         * @returns The route base name. if cannot get, `undefined` is returned.
         */
        $getRouteBaseName: (...args: Parameters<GetRouteBaseName>) => ReturnType<GetRouteBaseName>;
        /**
         * Returns localized path for passed in route.
         *
         * @remarks
         * If locale is not specified, uses current locale.
         *
         * @param route - A route.
         * @param locale - A locale, optional.
         *
         * @returns A path of the current route.
         */
        $localePath: (...args: Parameters<LocalePath>) => ReturnType<LocalePath>;
        /**
         * Returns localized route for passed in `route` parameters.
         *
         * @remarks
         * If `locale` is not specified, uses current locale.
         *
         * @param route - A route.
         * @param locale - A {@link Locale | locale}, optional.
         *
         * @returns A route. if cannot resolve, `undefined` is returned.
         */
        $localeRoute: (...args: Parameters<LocaleRoute>) => ReturnType<LocaleRoute>;
        /**
         * Returns localized head properties for locale-related aspects.
         *
         * @param options - An options, see about details [I18nHeadOptions](https://github.com/intlify/routing/blob/main/packages/vue-i18n-routing/api.md#i18nheadoptions).
         *
         * @returns The localized [head properties](https://github.com/intlify/routing/blob/main/packages/vue-i18n-routing/api.md#i18nheadmetainfo).
         */
        $localeHead: (...args: Parameters<LocaleHead>) => ReturnType<LocaleHead>;
        /**
         * Returns path of the current route for specified locale
         *
         * @param locale - A {@link Locale}
         *
         * @returns A path of the current route
         */
        $switchLocalePath: (...args: Parameters<SwitchLocalePath>) => ReturnType<SwitchLocalePath>;
    }
}
