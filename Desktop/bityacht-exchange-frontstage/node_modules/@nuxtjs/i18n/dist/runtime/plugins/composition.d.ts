import type { Composer, ExportedGlobalComposer } from 'vue-i18n';
import type { I18nRoutingCustomProperties } from 'vue-i18n-routing/dist/vue-i18n';
import type { NuxtI18nRoutingCustomProperties } from '../types';
declare const _default: any;
export default _default;
declare module '#app' {
    interface NuxtApp {
        $i18n: ExportedGlobalComposer & Composer & NuxtI18nRoutingCustomProperties & I18nRoutingCustomProperties;
    }
}
