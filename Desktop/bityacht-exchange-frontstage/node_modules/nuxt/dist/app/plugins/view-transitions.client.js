import { useRouter } from "#app/composables/router";
import { defineNuxtPlugin } from "#app/nuxt";
export default defineNuxtPlugin((nuxtApp) => {
  if (!document.startViewTransition) {
    return;
  }
  let finishTransition;
  let abortTransition;
  const router = useRouter();
  router.beforeResolve((to, from) => {
    if (to === from || to.matched.every((comp, index) => comp.components && comp.components?.default === from.matched[index]?.components?.default)) {
      return;
    }
    const promise = new Promise((resolve, reject) => {
      finishTransition = resolve;
      abortTransition = reject;
    });
    let changeRoute;
    const ready = new Promise((resolve) => changeRoute = resolve);
    const transition = document.startViewTransition(() => {
      changeRoute();
      return promise;
    });
    transition.finished.then(() => {
      abortTransition = void 0;
      finishTransition = void 0;
    });
    return ready;
  });
  nuxtApp.hook("vue:error", () => {
    abortTransition?.();
    abortTransition = void 0;
  });
  nuxtApp.hook("page:finish", () => {
    finishTransition?.();
    finishTransition = void 0;
  });
});
