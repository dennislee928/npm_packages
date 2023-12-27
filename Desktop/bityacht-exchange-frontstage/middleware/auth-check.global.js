export default defineNuxtRouteMiddleware((to, from) => {
  const isLogin = useCookie('isLogin');
  const nuxt = useNuxtApp();
  if (!isLogin.value) {
    if (to.path.includes('Members')) {
      const isVerify = to.fullPath.includes('verify');
      if(isVerify){
        return navigateTo({
          path: nuxt.$localePath('/login'),
          query: {
            isVerify:true
          },
        });
      }else{
        return navigateTo(nuxt.$localePath('/login'));
      }
    }
    if (to.path.includes('/Trade')) {
      return navigateTo(nuxt.$localePath('/login'));
    }
    if (to.path.includes('/MyAssets')) {
      return navigateTo(nuxt.$localePath('/login'));
    }
  }
  window.scrollTo(0, 0);
});
