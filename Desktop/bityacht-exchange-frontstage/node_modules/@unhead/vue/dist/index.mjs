import { whitelistSafeInput, unpackMeta, composableNames } from 'unhead';
export { createHeadCore } from 'unhead';
import { u as useHead, h as headSymbol, r as resolveUnrefHeadInput, s as serverUseHead } from './shared/vue.505b826e.mjs';
export { V as VueReactiveUseHeadPlugin, a as createHead, c as createServerHead, i as injectHead } from './shared/vue.505b826e.mjs';
import { getCurrentInstance, ref, watchEffect } from 'vue';
import { asArray } from '@unhead/shared';

const VueHeadMixin = {
  created() {
    const instance = getCurrentInstance();
    if (!instance)
      return;
    const options = instance.type;
    if (!options || !("head" in options))
      return;
    const source = typeof options.head === "function" ? () => options.head.call(instance.proxy) : options.head;
    useHead(source);
  }
};

const Vue2ProvideUnheadPlugin = function(_Vue, head) {
  _Vue.mixin({
    beforeCreate() {
      const options = this.$options;
      const origProvide = options.provide;
      options.provide = function() {
        let origProvideResult;
        if (typeof origProvide === "function")
          origProvideResult = origProvide.call(this);
        else
          origProvideResult = origProvide || {};
        return {
          ...origProvideResult,
          [headSymbol]: head
        };
      };
    }
  });
};

function useHeadSafe(input, options = {}) {
  return useHead(input, { ...options, transform: whitelistSafeInput });
}

function useSeoMeta(input, options) {
  const headInput = ref({});
  watchEffect(() => {
    const resolvedMeta = resolveUnrefHeadInput(input);
    const { title, titleTemplate, ...meta } = resolvedMeta;
    headInput.value = {
      title,
      titleTemplate,
      meta: unpackMeta(meta)
    };
  });
  return useHead(headInput, options);
}

function useServerHead(input, options = {}) {
  return serverUseHead(input, { ...options, mode: "server" });
}

function useServerHeadSafe(input, options = {}) {
  return useHeadSafe(input, { ...options, mode: "server" });
}

function useServerSeoMeta(input, options) {
  return useSeoMeta(input, { ...options || {}, mode: "server" });
}

function useTagTitle(title) {
  return useHead({ title });
}
function useTitleTemplate(titleTemplate) {
  return useHead({ titleTemplate });
}
function useTagMeta(meta) {
  return useHead({ meta: asArray(meta) });
}
function useTagMetaFlat(meta) {
  const input = ref({});
  watchEffect(() => {
    input.value = unpackMeta(resolveUnrefHeadInput(meta));
  });
  return useHead({ meta: input });
}
function useTagLink(link) {
  return useHead({ link: asArray(link) });
}
function useTagScript(script) {
  return useHead({ script: asArray(script) });
}
function useTagStyle(style) {
  return useHead({ style: asArray(style) });
}
function useTagNoscript(noscript) {
  return useHead({ noscript: asArray(noscript) });
}
function useTagBase(base) {
  return useHead({ base });
}
function useHtmlAttrs(attrs) {
  return useHead({ htmlAttrs: attrs });
}
function useBodyAttrs(attrs) {
  return useHead({ bodyAttrs: attrs });
}
function useServerTagTitle(title) {
  return useServerHead({ title });
}
function useServerTitleTemplate(titleTemplate) {
  return useServerHead({ titleTemplate });
}
function useServerTagMeta(meta) {
  return useServerHead({ meta: asArray(meta) });
}
function useServerTagMetaFlat(meta) {
  const input = ref({});
  watchEffect(() => {
    input.value = unpackMeta(resolveUnrefHeadInput(meta));
  });
  return useServerHead({ meta: input });
}
function useServerTagLink(link) {
  return useServerHead({ link: asArray(link) });
}
function useServerTagScript(script) {
  return useServerHead({ script: asArray(script) });
}
function useServerTagStyle(style) {
  return useServerHead({ style: asArray(style) });
}
function useServerTagNoscript(noscript) {
  return useServerHead({ noscript: asArray(noscript) });
}
function useServerTagBase(base) {
  return useServerHead({ base });
}
function useServerHtmlAttrs(attrs) {
  return useServerHead({ htmlAttrs: attrs });
}
function useServerBodyAttrs(attrs) {
  return useHead({ bodyAttrs: attrs });
}

const coreComposableNames = [
  "injectHead"
];
const unheadVueComposablesImports = {
  "@unhead/vue": [...coreComposableNames, ...composableNames]
};

export { Vue2ProvideUnheadPlugin, VueHeadMixin, headSymbol, resolveUnrefHeadInput, unheadVueComposablesImports, useBodyAttrs, useHead, useHeadSafe, useHtmlAttrs, useSeoMeta, useServerBodyAttrs, useServerHead, useServerHeadSafe, useServerHtmlAttrs, useServerSeoMeta, useServerTagBase, useServerTagLink, useServerTagMeta, useServerTagMetaFlat, useServerTagNoscript, useServerTagScript, useServerTagStyle, useServerTagTitle, useServerTitleTemplate, useTagBase, useTagLink, useTagMeta, useTagMetaFlat, useTagNoscript, useTagScript, useTagStyle, useTagTitle, useTitleTemplate };
