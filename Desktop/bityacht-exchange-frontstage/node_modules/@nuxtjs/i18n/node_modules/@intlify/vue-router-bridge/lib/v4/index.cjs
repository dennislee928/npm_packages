Object.defineProperty(exports, '__esModule', {
  value: true
})

var VueRouter = require('vue-router')

// stub vue-router 3 class
class VueRouterLegacy {
  static install() {}
  static version = ''
}

Object.keys(VueRouter).forEach(function (key) {
  VueRouterLegacy[key] = VueRouter[key]
})

VueRouterLegacy.isVueRouter3 = false
VueRouterLegacy.isVueRouter4 = true

exports.default = VueRouterLegacy
module.exports = VueRouterLegacy
