Object.defineProperty(exports, '__esModule', {
  value: true
})

var VueI18n = require('vue-i18n')

// stub vue-i18n 8 class
class VueI18nLegacy {
  static install() {}
  static version = ''
}

Object.keys(VueI18n).forEach(function (key) {
  VueI18nLegacy[key] = VueI18n[key]
})

VueI18nLegacy.isVueI18n8 = false
VueI18nLegacy.isVueI18n9 = true

exports.default = VueI18nLegacy
module.exports = VueI18nLegacy
