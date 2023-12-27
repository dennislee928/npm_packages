# @intlify/vue-router-bridge

Vue Router bridging for Vue 2 & Vue 3

> This library is inspired by [vue-demi](https://github.com/vueuse/vue-demi)

## 🌟 Features

- Vue Router composable APIs available on Vue 2 & Vue 3
  - `useRouter`
  - `useRoute`
- Stubbed Vue Router 4 APIs on Vue Router 3
  - About details [here](https://github.com/intlify/bridging/blob/63b1fb5dab4be772f4f564f023e9028979d37196/packages/vue-router-bridge/lib/v3/index.d.ts#L7-L75)
- Auto detect Vue Router version on bundling
- Manual switch versions

## 💿 Installation

```sh
# npm
npm install @intlify/vue-router-bridge

# yarn
yarn add @intlify/vue-router-bridge

# pnpm
pnpm add @intlify/vue-router-bridge
```

## ⛓️ Dependencies

You need to add `vue-router` and `@vue/composition-api` to your plugin's peer dependencies to specify what versions you support.

```js
{
  "dependencies": {
    "@intlify/vue-router-bridge": "latest"
  },
  "peerDependencies": {
    "@vue/composition-api": "^1.0.0-rc.1",
    "vue-router": "^3.0.0" // or "^4.0.0" base on your preferred working environment
  },
  "peerDependenciesMeta": {
    "@vue/composition-api": {
      "optional": true
    }
  },
  "devDependencies": {
    "vue-router": "^3.0.0" // or "^4.0.0" base on your preferred working environment
  },
}
```

Import everything related to Vue Router from it, it will redirect to `vue-router@3` + `@vue/composition-api` or `vue-router@4` based on users' environments.

```js
import { useRouter, useRoute } from '@intlify/vue-router-bridge'
```

When using with [Vite](https://vitejs.dev), you will need to opt-out the pre-bundling to get `@intlify/vue-router-bridge` work properly by

```js
// vite.config.js
export default defineConfig({
  optimizeDeps: {
    exclude: ['@intlify/vue-router-bridge']
  }
})
```

## 🤝 Extra APIs

`@intlify/vue-router-bridge` provides extra APIs to help distinguish users' environments and to do some version-specific logic.

### `isVueRouter3` / `isVueRouter4`

```js
import { isVueRouter3, isVueRouter4 } from '@intlify/vue-router-bridge'

if (isVueRouter3) {
  // Vue Router 3 only
} else {
  // Vue Router 4 only
}
```

## 📺 CLI

To explicitly switch the redirecting version, you can use these commands in your project's root:

### 🤏 Manually Switch Versions

```sh
npx vue-router-switch 3
# or
npx vue-router-switch 4
```

### 📦 Package Aliasing

If you would like to import `vue-router` under an alias, you can use the following command:

```sh
npx vue-router-switch 3 vue-router3
# or
npx vue-router-switch 4 vue-router4
```

### 🩹 Auto Fix

If the postinstall hook doesn't get triggered or you have updated the Vue Router version, try to run the following command to resolve the redirecting:

```sh
npx vue-router-fix
```

### ✳️ Isomorphic Testings

You can support testing for both versions by adding npm alias in your dev dependencies. For example:

```js
{
  "scripts": {
    "test:3": "vue-router-switch 3 vue-router3 && jest",
    "test:4": "vue-router-switch 4 && jest",
  },
  "devDependencies": {
    "vue-router": "^4.0.0",
    "vue-router3": "npm:vue-router@3"
  },
}
```

or

```js
{
  "scripts": {
    "test:3": "vue-router-switch 3 && jest",
    "test:4": "vue-router-switch 4 vue-router4 && jest",
  },
  "devDependencies": {
    "vue-router": "^3.0.0",
    "vue-router4": "npm:vue-router@4"
  },
}
```

## 🍭 Examples

- [Vue 2 + Vue Router v3.6 before](../../examples/vue-router-v3-vue2)
- [Vue 2 + Vue Router v3.6 or v3.x later](../../examples/vue-router-v36-vue2)
- [Vue 3 + Vue Router v4](../../examples/vue-router-v4-vue3)

## 💖 Thanks

This package idea was inspired from [vue-demi](https://github.com/vueuse/vue-demi), [@antfu](https://github.com/antfu)'s great work!

## ©️ License

[MIT](http://opensource.org/licenses/MIT)
