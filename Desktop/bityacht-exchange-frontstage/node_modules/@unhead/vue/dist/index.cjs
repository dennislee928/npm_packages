'use strict';

const unhead = require('unhead');
const useHead = require('./shared/vue.a57702be.cjs');
const vue = require('vue');
const shared = require('@unhead/shared');

const VueHeadMixin = {
  created() {
    const instance = vue.getCurrentInstance();
    if (!instance)
      return;
    const options = instance.type;
    if (!options || !("head" in options))
      return;
    const source = typeof options.head === "function" ? () => options.head.call(instance.proxy) : options.head;
    useHead.useHead(source);
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
          [useHead.headSymbol]: head
        };
      };
    }
  });
};

function useHeadSafe(input, options = {}) {
  return useHead.useHead(input, { ...options, transform: unhead.whitelistSafeInput });
}

function useSeoMeta(input, options) {
  const headInput = vue.ref({});
  vue.watchEffect(() => {
    const resolvedMeta = useHead.resolveUnrefHeadInput(input);
    const { title, titleTemplate, ...meta } = resolvedMeta;
    headInput.value = {
      title,
      titleTemplate,
      meta: unhead.unpackMeta(meta)
    };
  });
  return useHead.useHead(headInput, options);
}

function useServerHead(input, options = {}) {
  return useHead.serverUseHead(input, { ...options, mode: "server" });
}

function useServerHeadSafe(input, options = {}) {
  return useHeadSafe(input, { ...options, mode: "server" });
}

function useServerSeoMeta(input, options) {
  return useSeoMeta(input, { ...options || {}, mode: "server" });
}

function useTagTitle(title) {
  return useHead.useHead({ title });
}
function useTitleTemplate(titleTemplate) {
  return useHead.useHead({ titleTemplate });
}
function useTagMeta(meta) {
  return useHead.useHead({ meta: shared.asArray(meta) });
}
function useTagMetaFlat(meta) {
  const input = vue.ref({});
  vue.watchEffect(() => {
    input.value = unhead.unpackMeta(useHead.resolveUnrefHeadInput(meta));
  });
  return useHead.useHead({ meta: input });
}
function useTagLink(link) {
  return useHead.useHead({ link: shared.asArray(link) });
}
function useTagScript(script) {
  return useHead.useHead({ script: shared.asArray(script) });
}
function useTagStyle(style) {
  return useHead.useHead({ style: shared.asArray(style) });
}
function useTagNoscript(noscript) {
  return useHead.useHead({ noscript: shared.asArray(noscript) });
}
function useTagBase(base) {
  return useHead.useHead({ base });
}
function useHtmlAttrs(attrs) {
  return useHead.useHead({ htmlAttrs: attrs });
}
function useBodyAttrs(attrs) {
  return useHead.useHead({ bodyAttrs: attrs });
}
function useServerTagTitle(title) {
  return useServerHead({ title });
}
function useServerTitleTemplate(titleTemplate) {
  return useServerHead({ titleTemplate });
}
function useServerTagMeta(meta) {
  return useServerHead({ meta: shared.asArray(meta) });
}
function useServerTagMetaFlat(meta) {
  const input = vue.ref({});
  vue.watchEffect(() => {
    input.value = unhead.unpackMeta(useHead.resolveUnrefHeadInput(meta));
  });
  return useServerHead({ meta: input });
}
function useServerTagLink(link) {
  return useServerHead({ link: shared.asArray(link) });
}
function useServerTagScript(script) {
  return useServerHead({ script: shared.asArray(script) });
}
function useServerTagStyle(style) {
  return useServerHead({ style: shared.asArray(style) });
}
function useServerTagNoscript(noscript) {
  return useServerHead({ noscript: shared.asArray(noscript) });
}
function useServerTagBase(base) {
  return useServerHead({ base });
}
function useServerHtmlAttrs(attrs) {
  return useServerHead({ htmlAttrs: attrs });
}
function useServerBodyAttrs(attrs) {
  return useHead.useHead({ bodyAttrs: attrs });
}

const coreComposableNames = [
  "injectHead"
];
const unheadVueComposablesImports = {
  "@unhead/vue": [...coreComposableNames, ...unhead.composableNames]
};

exports.createHeadCore = unhead.createHeadCore;
exports.VueReactiveUseHeadPlugin = useHead.VueReactiveUseHeadPlugin;
exports.createHead = useHead.createHead;
exports.createServerHead = useHead.createServerHead;
exports.headSymbol = useHead.headSymbol;
exports.injectHead = useHead.injectHead;
exports.resolveUnrefHeadInput = useHead.resolveUnrefHeadInput;
exports.useHead = useHead.useHead;
exports.Vue2ProvideUnheadPlugin = Vue2ProvideUnheadPlugin;
exports.VueHeadMixin = VueHeadMixin;
exports.unheadVueComposablesImports = unheadVueComposablesImports;
exports.useBodyAttrs = useBodyAttrs;
exports.useHeadSafe = useHeadSafe;
exports.useHtmlAttrs = useHtmlAttrs;
exports.useSeoMeta = useSeoMeta;
exports.useServerBodyAttrs = useServerBodyAttrs;
exports.useServerHead = useServerHead;
exports.useServerHeadSafe = useServerHeadSafe;
exports.useServerHtmlAttrs = useServerHtmlAttrs;
exports.useServerSeoMeta = useServerSeoMeta;
exports.useServerTagBase = useServerTagBase;
exports.useServerTagLink = useServerTagLink;
exports.useServerTagMeta = useServerTagMeta;
exports.useServerTagMetaFlat = useServerTagMetaFlat;
exports.useServerTagNoscript = useServerTagNoscript;
exports.useServerTagScript = useServerTagScript;
exports.useServerTagStyle = useServerTagStyle;
exports.useServerTagTitle = useServerTagTitle;
exports.useServerTitleTemplate = useServerTitleTemplate;
exports.useTagBase = useTagBase;
exports.useTagLink = useTagLink;
exports.useTagMeta = useTagMeta;
exports.useTagMetaFlat = useTagMetaFlat;
exports.useTagNoscript = useTagNoscript;
exports.useTagScript = useTagScript;
exports.useTagStyle = useTagStyle;
exports.useTagTitle = useTagTitle;
exports.useTitleTemplate = useTitleTemplate;
