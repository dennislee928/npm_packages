const fs = require('fs') // eslint-disable-line @typescript-eslint/no-var-requires
const path = require('path') // eslint-disable-line @typescript-eslint/no-var-requires

const DEBUG = process.env.DEBUG

const dir = path.resolve(__dirname, '..', 'lib')

function loadModule(name) {
  try {
    return require(name)
  } catch (e) {
    DEBUG && error(e)
    return undefined
  }
}

function error(...args) {
  console.error(`[vue-router-bridge] `, ...args)
}

function warn(...args) {
  console.warn(`[vue-router-bridge] `, ...args)
}

function log(...args) {
  console.log(`[vue-router-bridge] `, ...args)
}

function copy(name, version, router, esm = false) {
  const src = path.join(dir, `v${version}`, name)
  const dest = path.join(dir, name)
  let content = fs.readFileSync(src, 'utf-8')
  content = esm
    ? content.replace(/from 'vue-router'/g, `from '${router}'`)
    : content.replace(/require\('vue-router'\)/g, `require('${router}')`)
  try {
    fs.unlinkSync(dest)
  } catch (error) {}
  fs.writeFileSync(dest, content, 'utf-8')
}

function checkVueRouter(pkg) {
  const router = loadModule(pkg)
  if (!router) {
    warn('Vue Router plugin is not found. Please run "npm install vue-router" to install.')
    return false
  }
  return true
}

function checkVCA() {
  const demi = loadModule('vue-demi')
  if (!demi) {
    return false
  }

  if (demi.Vue.version.startsWith('2.7.')) {
    return true
  } else if (demi.Vue.version.startsWith('2.')) {
    const VCA = loadModule('@vue/composition-api')
    if (!VCA) {
      warn('Composition API plugin is not found. Please run "npm install @vue/composition-api" to install.')
      return false
    }
    return true
  } else {
    return false
  }
}

function switchVersion(version, router) {
  router = router || 'vue-router'
  if (!checkVueRouter(router)) {
    return
  }
  if (version === 3 && !checkVCA()) {
    return
  }
  copy('index.cjs', version, router)
  copy('index.mjs', version, router, true)
  copy('index.d.ts', version, router, true)
}

module.exports.warn = warn
module.exports.log = log
module.exports.loadModule = loadModule
module.exports.switchVersion = switchVersion
