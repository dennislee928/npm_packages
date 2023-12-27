var isVueI18n8 = false
var isVueI18n9 = true

// stub vue-i18n 8 class
class VueI18nLegacy {
  static install() {}
  static version = ''
}

export * from 'vue-i18n'
export { isVueI18n8, isVueI18n9 }
export default VueI18nLegacy
