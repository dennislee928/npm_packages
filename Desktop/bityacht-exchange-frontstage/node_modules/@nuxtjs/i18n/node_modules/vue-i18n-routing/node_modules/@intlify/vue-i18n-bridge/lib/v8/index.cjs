Object.defineProperty(exports, '__esModule', {
  value: true
})

var VueI18n = require('vue-i18n')
var VueI18nBridge = require('vue-i18n-bridge')

Object.keys(VueI18nBridge).forEach(function (key) {
  VueI18n[key] = VueI18nBridge[key]
})
VueI18n.isVueI18n8 = true
VueI18n.isVueI18n9 = false

exports.default = VueI18n
module.exports = VueI18n
