<template>
  <NuxtLayout name="members">
    <section v-if="dataReady" class="relative bg-white 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[30px]">
      <template v-if="userData.type === 1">
        <div class="flex justify-between lg:flex-row flex-col relative z-30">
          <MembersLevelCard :level="1" :real-level="realLevel">
            <div class="flex xs:flex-row flex-col">
              <div class="flex items-start 3xl:mr-[180px] 2xl:mr-[80px] xl:mr-10 mr-5 xs:mb-0 mb-8">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
                <div class="flex flex-col">
                  <p class="text-black text-subTitle mb-2">
                    {{ $t('memberCenter.levelText1') }}
                  </p>
                  <p class="text-waterBlue text-small">
                    {{ $t('memberCenter.levelOneText1') }}
                  </p>
                </div>
              </div>
              <div class="flex items-start 3xl:mr-5 mr-0">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
                <div class="flex flex-col">
                  <p class="text-black text-subTitle mb-2">
                    {{ $t('memberCenter.levelText2') }}
                  </p>
                  <p class="text-brown text-small">
                    {{ $t('memberCenter.levelOneText2') }}
                  </p>
                  <p class="text-brown text-small">
                    {{ $t('memberCenter.levelOneText3') }}
                  </p>
                </div>
              </div>
            </div>
          </MembersLevelCard>
          <MembersLevelCard :level="2" :real-level="realLevel">
            <div class="flex xs:flex-row flex-col">
              <div class="flex items-start 3xl:mr-10 2xl:mr-6 mr-4 xs:mb-0 mb-8">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
                <div class="flex flex-col">
                  <p class="text-black text-subTitle mb-2">
                    {{ $t('memberCenter.levelText1') }}
                  </p>
                  <p class="text-waterBlue text-small">
                    {{ $t('memberCenter.levelTwoText1') }}
                  </p>
                </div>
              </div>
              <div class="flex items-start">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
                <div class="flex flex-col">
                  <p class="text-black text-subTitle">
                    {{ $t('memberCenter.levelText2') }}
                  </p>
                  <div class="flex md:flex-row flex-col">
                    <div>
                      <p class="text-waterBlue text-small my-2">
                        {{ $t('memberCenter.levelTwoText2') }}
                      </p>
                      <div class="flex">
                        <div class="flex flex-col 3xl:mr-10 2xl:mr-6 mr-4">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/saveMoney.svg" class="mr-2" />
                            {{ $t('memberCenter.levelTwoText4') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelTwoText6') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelTwoText7') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelTwoText8') }}</p>
                        </div>
                        <div class="flex flex-col">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/getMoney.svg" class="mr-2 h-[18px]" />
                            {{ $t('memberCenter.levelTwoText5') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelTwoText9') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelTwoText10') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelTwoText11') }}</p>
                        </div>
                      </div>
                    </div>
                    <div class="3xl:ml-10 lg:ml-4 md:ml-3 ml-0">
                      <p class="text-waterBlue text-small my-2">
                        {{ $t('memberCenter.levelTwoText3') }}
                      </p>
                      <div class="flex">
                        <div class="flex flex-col 3xl:mr-10 2xl:mr-6 mr-4">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/saveMoney.svg" class="mr-2" />
                            {{ $t('memberCenter.levelTwoText4') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelTwoText12') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelTwoText13') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelTwoText14') }}</p>
                        </div>
                        <div class="flex flex-col">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/getMoney.svg" class="mr-2 h-[18px]" />
                            {{ $t('memberCenter.levelTwoText5') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelTwoText15') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelTwoText16') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelTwoText17') }}</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </MembersLevelCard>
        </div>
        <div class="flex justify-between lg:flex-row flex-col relative z-30">
          <MembersLevelCard :level="3" :real-level="realLevel">
            <div class="flex lg:flex-col xs:flex-row flex-col 3xl:pr-[55px] 2xl:pr-6 xl:pr-3 pr-0">
              <div class="flex items-start 3xl:mr-10 2xl:mr-6 mr-4 mb-8">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
                <div class="flex flex-col">
                  <p class="text-black text-subTitle mb-2">
                    {{ $t('memberCenter.levelText1') }}
                  </p>
                  <p class="text-waterBlue text-small">
                    {{ $t('memberCenter.levelThreeText1') }}
                  </p>
                  <p class="text-waterBlue text-small">
                    {{ $t('memberCenter.levelThreeText2') }}
                  </p>
                </div>
              </div>
              <div class="flex items-start">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
                <div class="flex flex-col">
                  <p class="text-black text-subTitle">
                    {{ $t('memberCenter.levelText2') }}
                  </p>
                  <div class="flex lg:flex-row md:flex-row flex-col">
                    <div>
                      <p class="text-waterBlue text-small my-2">
                        {{ $t('memberCenter.levelTwoText2') }}
                      </p>
                      <div class="flex lg:flex-col flex-row">
                        <div class="flex flex-col 3xl:mr-10 2xl:mr-6 mr-4 mb-4">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/saveMoney.svg" class="mr-2" />
                            {{ $t('memberCenter.levelTwoText4') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelThreeText3') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelThreeText4') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelThreeText5') }}</p>
                        </div>
                        <div class="flex flex-col">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/getMoney.svg" class="mr-2 h-[18px]" />
                            {{ $t('memberCenter.levelTwoText5') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelThreeText6') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelThreeText7') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelThreeText8') }}</p>
                        </div>
                      </div>
                    </div>
                    <div class="3xl:ml-10 lg:ml-4 md:ml-3 ml-0">
                      <p class="text-waterBlue text-small my-2">
                        {{ $t('memberCenter.levelTwoText3') }}
                      </p>
                      <div class="flex lg:flex-col flex-row">
                        <div class="flex flex-col 3xl:mr-10 2xl:mr-6 mr-4 mb-4">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/saveMoney.svg" class="mr-2" />
                            {{ $t('memberCenter.levelTwoText4') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelThreeText9') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelThreeText10') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelThreeText11') }}</p>
                        </div>
                        <div class="flex flex-col">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/getMoney.svg" class="mr-2 h-[18px]" />
                            {{ $t('memberCenter.levelTwoText5') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelThreeText12') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelThreeText13') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelThreeText14') }}</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </MembersLevelCard>
          <MembersLevelCard :level="4" :real-level="realLevel">
            <div class="flex lg:flex-col xs:flex-row flex-col 3xl:pr-[55px] 2xl:pr-6 xl:pr-3 pr-0">
              <div class="flex items-start 3xl:mr-10 2xl:mr-6 mr-4 mb-8">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
                <div class="flex flex-col">
                  <p class="text-black text-subTitle mb-2">
                    {{ $t('memberCenter.levelText1') }}
                  </p>
                  <p class="text-waterBlue text-small">
                    {{ $t('memberCenter.levelFourText1') }}
                  </p>
                  <p class="text-waterBlue text-small">
                    {{ $t('memberCenter.levelFourText2') }}
                  </p>
                </div>
              </div>
              <div class="flex items-start">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
                <div class="flex flex-col">
                  <p class="text-black text-subTitle">
                    {{ $t('memberCenter.levelText2') }}
                  </p>
                  <div class="flex lg:flex-row md:flex-row flex-col">
                    <div>
                      <p class="text-waterBlue text-small my-2">
                        {{ $t('memberCenter.levelTwoText2') }}
                      </p>
                      <div class="flex lg:flex-col flex-row">
                        <div class="flex flex-col 3xl:mr-10 2xl:mr-6 mr-4 mb-4">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/saveMoney.svg" class="mr-2" />
                            {{ $t('memberCenter.levelTwoText4') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFourText3') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFourText4') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFourText5') }}</p>
                        </div>
                        <div class="flex flex-col">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/getMoney.svg" class="mr-2 h-[18px]" />
                            {{ $t('memberCenter.levelTwoText5') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFourText6') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFourText7') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFourText8') }}</p>
                        </div>
                      </div>
                    </div>
                    <div class="3xl:ml-10 lg:ml-4 md:ml-3 ml-0">
                      <p class="text-waterBlue text-small my-2">
                        {{ $t('memberCenter.levelTwoText3') }}
                      </p>
                      <div class="flex lg:flex-col flex-row">
                        <div class="flex flex-col 3xl:mr-10 2xl:mr-6 mr-4 mb-4">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/saveMoney.svg" class="mr-2" />
                            {{ $t('memberCenter.levelTwoText4') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFourText9') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFourText10') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFourText11') }}</p>
                        </div>
                        <div class="flex flex-col">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/getMoney.svg" class="mr-2 h-[18px]" />
                            {{ $t('memberCenter.levelTwoText5') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFourText12') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFourText13') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFourText14') }}</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </MembersLevelCard>
          <MembersLevelCard :level="5" :real-level="realLevel">
            <div class="flex lg:flex-col xs:flex-row flex-col 3xl:pr-[55px] 2xl:pr-6 xl:pr-3 pr-0">
              <div class="flex items-start 3xl:mr-10 2xl:mr-6 mr-4 mb-8">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
                <div class="flex flex-col">
                  <p class="text-black text-subTitle mb-2">
                    {{ $t('memberCenter.levelText1') }}
                  </p>
                  <p class="text-waterBlue text-small">
                    {{ $t('memberCenter.levelFiveText1') }}
                  </p>
                  <p class="text-waterBlue text-small">
                    {{ $t('memberCenter.levelFiveText2') }}
                  </p>
                </div>
              </div>
              <div class="flex items-start">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
                <div class="flex flex-col">
                  <p class="text-black text-subTitle">
                    {{ $t('memberCenter.levelText2') }}
                  </p>
                  <div class="flex lg:flex-row md:flex-row flex-col">
                    <div>
                      <p class="text-waterBlue text-small my-2">
                        {{ $t('memberCenter.levelTwoText2') }}
                      </p>
                      <div class="flex lg:flex-col flex-row">
                        <div class="flex flex-col 3xl:mr-10 2xl:mr-6 mr-4 mb-4">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/saveMoney.svg" class="mr-2" />
                            {{ $t('memberCenter.levelTwoText4') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFiveText3') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFiveText4') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFiveText5') }}</p>
                        </div>
                        <div class="flex flex-col">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/getMoney.svg" class="mr-2 h-[18px]" />
                            {{ $t('memberCenter.levelTwoText5') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFiveText6') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFiveText7') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFiveText8') }}</p>
                        </div>
                      </div>
                    </div>
                    <div class="3xl:ml-10 lg:ml-4 md:ml-3 ml-0">
                      <p class="text-waterBlue text-small my-2">
                        {{ $t('memberCenter.levelTwoText3') }}
                      </p>
                      <div class="flex lg:flex-col flex-row">
                        <div class="flex flex-col 3xl:mr-10 2xl:mr-6 mr-4 mb-4">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/saveMoney.svg" class="mr-2" />
                            {{ $t('memberCenter.levelTwoText4') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFiveText9') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFiveText10') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFiveText11') }}</p>
                        </div>
                        <div class="flex flex-col">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/getMoney.svg" class="mr-2 h-[18px]" />
                            {{ $t('memberCenter.levelTwoText5') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFiveText12') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFiveText13') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.levelFiveText14') }}</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </MembersLevelCard>
        </div>
      </template>
      <template v-else>
        <div class="flex justify-between lg:flex-row flex-col relative z-30">
          <MembersLevelCard :level="1" :is-blue="true" :real-level="realLevel">
            <div class="flex lg:flex-col md:flex-row flex-col 3xl:pr-[55px] 2xl:pr-6 xl:pr-3 pr-0">
              <div class="flex items-start 3xl:mr-10 2xl:mr-6 mr-4 xs:mb-0 mb-8">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
                <div class="flex flex-col mb-4">
                  <p class="text-black text-subTitle mb-2">
                    {{ $t('memberCenter.levelText1') }}
                  </p>
                  <p class="text-waterBlue text-small">
                    {{ $t('memberCenter.blueLevelOneText1') }}
                  </p>
                </div>
              </div>
              <div class="flex items-start">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
                <div class="flex flex-col">
                  <p class="text-black text-subTitle">
                    {{ $t('memberCenter.levelText2') }}
                  </p>
                  <div class="flex md:flex-row flex-col">
                    <div>
                      <p class="text-waterBlue text-small my-2">
                        {{ $t('memberCenter.levelTwoText2') }}
                      </p>
                      <div class="flex">
                        <div class="flex flex-col 3xl:mr-10 2xl:mr-6 mr-4">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/saveMoney.svg" class="mr-2" />
                            {{ $t('memberCenter.levelTwoText4') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelOneText2') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelOneText3') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelOneText4') }}</p>
                        </div>
                        <div class="flex flex-col">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/getMoney.svg" class="mr-2 h-[18px]" />
                            {{ $t('memberCenter.levelTwoText5') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelOneText5') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelOneText6') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelOneText7') }}</p>
                        </div>
                      </div>
                    </div>
                    <div class="3xl:ml-10 lg:ml-4 md:ml-3 ml-0">
                      <p class="text-waterBlue text-small my-2">
                        {{ $t('memberCenter.levelTwoText3') }}
                      </p>
                      <div class="flex">
                        <div class="flex flex-col 3xl:mr-10 2xl:mr-6 mr-4">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/saveMoney.svg" class="mr-2" />
                            {{ $t('memberCenter.levelTwoText4') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelOneText8') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelOneText9') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelOneText10') }}</p>
                        </div>
                        <div class="flex flex-col">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/getMoney.svg" class="mr-2 h-[18px]" />
                            {{ $t('memberCenter.levelTwoText5') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelOneText11') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelOneText12') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelOneText13') }}</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </MembersLevelCard>
          <MembersLevelCard :level="2" :is-blue="true" :real-level="realLevel">
            <div class="flex lg:flex-col md:flex-row flex-col 3xl:pr-[55px] 2xl:pr-6 xl:pr-3 pr-0">
              <div class="flex items-start 3xl:mr-10 2xl:mr-6 mr-4 xs:mb-0 mb-8">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
                <div class="flex flex-col mb-4">
                  <p class="text-black text-subTitle mb-2">
                    {{ $t('memberCenter.levelText1') }}
                  </p>
                  <p class="text-waterBlue text-small">
                    {{ $t('memberCenter.blueLevelTwoText1') }}
                  </p>
                  <p class="text-waterBlue text-small">
                    {{ $t('memberCenter.blueLevelTwoText2') }}
                  </p>
                </div>
              </div>
              <div class="flex items-start">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
                <div class="flex flex-col">
                  <p class="text-black text-subTitle">
                    {{ $t('memberCenter.levelText2') }}
                  </p>
                  <div class="flex md:flex-row flex-col">
                    <div>
                      <p class="text-waterBlue text-small my-2">
                        {{ $t('memberCenter.levelTwoText2') }}
                      </p>
                      <div class="flex">
                        <div class="flex flex-col 3xl:mr-10 2xl:mr-6 mr-4">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/saveMoney.svg" class="mr-2" />
                            {{ $t('memberCenter.levelTwoText4') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelTwoText3') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelTwoText4') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelTwoText5') }}</p>
                        </div>
                        <div class="flex flex-col">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/getMoney.svg" class="mr-2 h-[18px]" />
                            {{ $t('memberCenter.levelTwoText5') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelTwoText6') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelTwoText7') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelTwoText8') }}</p>
                        </div>
                      </div>
                    </div>
                    <div class="3xl:ml-10 lg:ml-4 md:ml-3 ml-0">
                      <p class="text-waterBlue text-small my-2">
                        {{ $t('memberCenter.levelTwoText3') }}
                      </p>
                      <div class="flex">
                        <div class="flex flex-col 3xl:mr-10 2xl:mr-6 mr-4">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/saveMoney.svg" class="mr-2" />
                            {{ $t('memberCenter.levelTwoText4') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelTwoText9') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelTwoText10') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelTwoText11') }}</p>
                        </div>
                        <div class="flex flex-col">
                          <p class="flex text-small text-waterBlue mb-1">
                            <img src="/icon/getMoney.svg" class="mr-2 h-[18px]" />
                            {{ $t('memberCenter.levelTwoText5') }}
                          </p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelTwoText12') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelTwoText13') }}</p>
                          <p class="text-small text-gray_800 mb-1">{{ $t('memberCenter.blueLevelTwoText14') }}</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </MembersLevelCard>
        </div>
      </template>
    </section>
  </NuxtLayout>
</template>
<script setup>
definePageMeta({
  middleware: defineNuxtRouteMiddleware(() => {
    const isLogin = useCookie('isLogin');
    if (!isLogin.value) {
      return navigateTo('/login');
    }
  }),
});

const dataReady = ref(false);
const userData = ref();
userData.value = JSON.parse(localStorage.getItem('userInfo'));
const realLevel = userData.value.level;

onMounted(() => {
  dataReady.value = true;
});
</script>
