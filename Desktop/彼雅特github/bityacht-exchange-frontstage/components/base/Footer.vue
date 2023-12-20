<template>
  <section class="relative w-full bg-darkBlueBg overflow-hidden z-[15]">
    <div class="2xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[50px] py-10">
      <div class="flex flex-col">
        <img src="/img/logoWhite.svg" class="w-[333px]" />
        <p class="mt-9">
          <span class="text-h5 text-white">{{ $t('index.download') }} </span><span class="text-h3 text-waterBlue ml-3">{{ $t('index.bityacht') }}</span>
        </p>
        <p class="text-subTitle text-white mt-3">{{ $t('index.bestParner') }}</p>
        <div class="flex justify-between md:flex-row flex-col z-30">
          <div class="flex flex-col">
            <!-- <span class="flex mt-5">
              <img src="/assets/img/appStore.png" class="w-[165px] cursor-pointer" />
              <img src="/assets/img/googlePlay.png" class="w-[190px] ml-4 cursor-pointer" />
            </span> -->
            <!-- <p class="text-subText text-white mt-10">{{ $t('index.connectPhone') }} : 0800-638-888</p>
            <p class="text-subText text-white mt-4">{{ $t('index.phonePleaseCall') }} : 02-2588-2588</p> -->
            <p class="text-subText text-white mt-4">{{ $t('index.LineFooter') }} : <a class="text-white hover:text-waterBlue transition-all duration-300" href="https://lin.ee/WV9O4yx" target="_blank">@023ccabq</a></p>
            <p class="text-subText text-white mt-4">
              {{ $t('index.EmailFooter') }} :
              <a class="text-white hover:text-waterBlue transition-all duration-300" href="mailto:support@bityacht.io">support@bityacht.io</a>
            </p>
          </div>
          <div class="flex md:flex-col flex-row mt-5">
            <p class="text-subTitle text-waterBlue mb-7 sm:w-auto xs:w-[80px] w-[50px]">{{ $t('index.aboutBitYacht') }}</p>
            <div class="md:block flex flex-wrap md:mt-0 mt-3 md:ml-0 ml-10">
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">
                <a href="https://medium.com/bityacht/%E9%97%9C%E6%96%BC%E6%88%91%E5%80%91-d6d88bc9c5a5" target="_blank">{{ $t('index.aboutUs') }}</a>
              </p>
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">
                <a href="https://medium.com/bityacht/%E4%BA%BA%E6%89%8D%E6%8B%9B%E5%8B%9F-2e56f9fe7bdf" target="_blank">{{ $t('index.findPeople') }}</a>
              </p>
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">
                <a href="https://medium.com/bityacht/%E4%BD%BF%E7%94%A8%E8%80%85%E5%8D%94%E8%AD%B0-aa8bfcf36c86" target="_blank">{{ $t('index.userPolicy') }}</a>
              </p>
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">
                <a href="https://medium.com/bityacht/%E9%9A%B1%E7%A7%81%E6%AC%8A%E6%94%BF%E7%AD%96-1d1a54909f33" target="_blank">{{ $t('index.privacyPolicy') }}</a>
              </p>
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">
                <a href="https://medium.com/bityacht/%E9%A2%A8%E9%9A%AA%E8%81%B2%E6%98%8E-5bbb7c7abbdd" target="_blank">{{ $t('index.disclaimer') }}</a>
              </p>
            </div>
          </div>
          <div class="flex md:flex-col flex-row mt-5">
            <p class="text-subTitle text-waterBlue mb-7 xs:w-[80px] w-[50px]">{{ $t('index.information') }}</p>
            <div class="md:block flex flex-wrap md:mt-0 mt-3 md:ml-0 ml-10">
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">
                <a href="https://medium.com/bityacht/%E7%94%A2%E5%93%81%E8%B3%87%E8%A8%8A-64b8b7ce890a" target="_blank">{{ $t('index.rate') }}</a>
              </p>
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">
                <a href="https://medium.com/bityacht/%E6%8F%90%E9%A0%98%E9%99%90%E5%88%B6-8fb0942259a8" target="_blank">{{ $t('index.limit') }}</a>
              </p>
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">
                <a href="https://medium.com/bityacht/%E4%BA%A4%E6%98%93%E8%A6%8F%E5%89%87-1ae039e568f1" target="_blank">{{ $t('index.tradeRules') }}</a>
              </p>
            </div>
          </div>
          <div class="flex md:flex-col flex-row mt-5">
            <p class="text-subTitle text-waterBlue mb-7 xs:w-[80px] w-[50px]">{{ $t('index.help') }}</p>
            <div class="md:block flex flex-wrap md:mt-0 mt-3 md:ml-0 ml-10 md:w-full xxs:w-[30%] w-[65%]">
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">
                <a href="https://medium.com/bityacht/%E5%B8%B8%E8%A6%8B%E5%95%8F%E9%A1%8C/home" target="_blank">{{ $t('index.normalQuestion') }}</a>
              </p>
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">
                <a href="https://medium.com/bityacht/%E5%B9%B3%E5%8F%B0%E5%85%AC%E5%91%8A/home" target="_blank">{{ $t('index.announcement') }}</a>
              </p>
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">
                <a href="https://medium.com/bityacht/%E6%9C%80%E6%96%B0%E6%B4%BB%E5%8B%95/home" target="_blank">{{ $t('index.newAction') }}</a>
              </p>
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">
                <a href="https://medium.com/bityacht/%E6%B3%95%E4%BA%BA%E5%AE%A2%E6%88%B6%E7%94%B3%E8%AB%8B/home" target="_blank">{{ $t('index.apply') }}</a>
              </p>
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">
                <a href="https://medium.com/bityacht/%E8%81%AF%E7%B5%A1%E6%88%91%E5%80%91-7975575da325" target="_blank">{{ $t('index.connectUs') }}</a>
              </p>
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">
                <a href="https://medium.com/bityacht/%E9%98%B2%E8%A9%90%E6%AC%BA%E5%AE%A3%E5%B0%8E/home" target="_blank">{{ $t('index.fraudPrevention') }}</a>
              </p>
            </div>
          </div>
          <!-- <div class="flex md:flex-col flex-row mt-5">
            <p class="text-subTitle text-waterBlue mb-7">{{ $t('index.about') }}</p>
            <div class="md:block flex flex-wrap md:mt-0 mt-3 md:ml-0 ml-10">
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">{{ $t('index.aboutUs') }}</p>
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">{{ $t('index.CEOStory') }}</p>
              <p class="text-small text-white mb-6 cursor-pointer hover:text-waterBlue transition-all duration-300 md:w-full w-[50%]">{{ $t('index.officialWebsite') }}</p>
            </div>
          </div> -->
        </div>
      </div>
      <!-- <div class="flex justify-center md:mt-[75px] mt-10">
        <span v-for="item of iconList" class="md:mx-4 mx-2">
          <NuxtLink :to="item.href" target="_blank">
            <img :src="item.img" />
          </NuxtLink>
        </span>
      </div> -->
      <div class="text-subText text-white md:mt-5 mt-9 text-center -mx-2">Copyright © 2023 {{ $t('index.footerCopyRight') }}</div>
    </div>
    <img src="/assets/img/Group16.png" class="absolute top-0 right-0 md:w-[500px] w-[350px]" />
    <img src="/assets/img/Group37.png" class="absolute top-0 left-0 w-[500px] hidden xl:block" />
  </section>
</template>

<script setup>
// const iconList = [
//   {
//     img: 'icon/facebook.svg',
//     href: 'https://yahoo.com.tw',
//   },
//   {
//     img: 'icon/youtube.svg',
//     href: 'https://yahoo.com.tw',
//   },
//   {
//     img: 'icon/IG.svg',
//     href: 'https://yahoo.com.tw',
//   },
//   {
//     img: 'icon/line.svg',
//     href: 'https://yahoo.com.tw',
//   },
//   {
//     img: 'icon/lineChat.svg',
//     href: 'https://yahoo.com.tw',
//   },
//   {
//     img: 'icon/podcast.svg',
//     href: 'https://yahoo.com.tw',
//   },
// ];
</script>
