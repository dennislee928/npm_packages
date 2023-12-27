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
  console.error(`[vue-i18n-bridge] `, ...args)
}

function warn(...args) {
  console.warn(`[vue-i18n-bridge] `, ...args)
}

function log(...args) {
  console.log(`[vue-i18n-bridge] `, ...args)
}

function copy(name, version, i18n, esm = false) {
  const src = path.join(dir, `v${version}`, name)
  const dest = path.join(dir, name)
  let content = fs.readFileSync(src, 'utf-8')
  content = esm
    ? content.replace(/from 'vue-i18n'/g, `from '${i18n}'`)
    : content.replace(/require\('vue-i18n'\)/g, `require('${i18n}')`)
  try {
    fs.unlinkSync(dest)
  } catch (error) {}
  fs.writeFileSync(dest, content, 'utf-8')
}

function checkBridge() {
  const bridge = loadModule('vue-i18n-bridge')
  if (!bridge) {
    warn('Vue I18n Bridge plugin is not found. Please run "npm install vue-i18n-bridge" to install.')
    return false
  }
  return true
}

function checkVueI18n(pkg) {
  const i18n = loadModule(pkg)
  if (!i18n) {
    warn(`Vue I18n plugin is not found. Please run "npm install ${pkg}" to install.`)
    return false
  }
  return true
}

function switchVersion(version, i18n) {
  i18n = i18n || 'vue-i18n'
  if (!checkVueI18n(i18n)) {
    return
  }
  if (version === 8 && !checkBridge()) {
    return
  }
  copy('index.cjs', version, i18n)
  copy('index.mjs', version, i18n, true)
  copy('index.d.ts', version, i18n, true)
}

module.exports.warn = warn
module.exports.log = log
module.exports.loadModule = loadModule
module.exports.switchVersion = switchVersion
