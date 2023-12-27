<template>
  <!-- <div :id="id" :class="`h-[${height}px] 3xl:w-[${width}px] 2xl:w-[${width - 100}px] xl:w-[${width - 200}px] lg:w-[${width - 250}px] w-[${width - 340}px]`"></div> -->
  <template v-if="location === 'index'">
    <div :id="id" class="h-[80px] xl:w-[550px] lg:w-[350px] w-[300px]"></div>
  </template>
  <template v-if="location === 'trade'">
    <div :id="id" class="h-[80px] 2xl:w-[500px] xl:w-[300px] lg:w-[250px] w-[150px]"></div>
  </template>
  <template v-if="location === 'myAssets'">
    <div :id="id" class="h-[250px] 3xl:w-[900px] 2xl:w-[800px] lg:w-[650px] md:w-[500px] sm:w-[550px] xs:w-[500px] xxs:w-[450px] xxxs:w-[360px] w-[320px]"></div>
  </template>
</template>
<script>
import * as echarts from 'echarts';
import { useNow, useDateFormat } from '@vueuse/core';
import { formatValueByDigits } from '@/config/config';

export default {
  // props: ['id', 'title', 'xAxisData', 'yAxis', 'series', 'max', 'min', 'height', 'width'],
  props: {
    id: String || Number,
    title: String || Number,
    series: Array,
    max: Number,
    min: Number,
    location: String,
    show: Boolean,
    symbol: String,
    inverse: Boolean,
  },
  setup(props) {
    const lineChart = shallowRef(null);
    const initChartDom = () => {
      lineChart.value = echarts.getInstanceByDom(document.getElementById(props.id));
      if (!lineChart.value) lineChart.value = echarts.init(document.getElementById(props.id));
    };
    const resizeChart = () => {
      nextTick(() => {
        window.onresize = () => {
          lineChart.value.resize();
        };
      });
    };
    const initChart = () => {
      // responsive chart size
      const option = {
        title: {
          text: props.title || '',
        },
        tooltip: {
          // show: true,
          // order: 'seriesDesc',
          formatter: function (params) {
            // console.log('params :>> ', params);
            var html = '';
            let now = useNow();
            let fakeNow = new Date(now.value);
            fakeNow.setDate(fakeNow.getDate() - Number(params.name));
            fakeNow = useDateFormat(fakeNow, 'YYYY/MM/DD').value;
            html += `${fakeNow}：<br /><span style="color: #13458c">${formatValueByDigits(params.value, 2)}</span><br />`;
            return html;
          },
        },
        legend: {},
        xAxis: {
          type: 'category',
          // minInterval: 5,
          // splitNumber: 500,
          inverse: props.inverse,
          show: props.show,
          scale: true,
          splitLine: { show: false },
          axisLabel: {
            formatter: function (value) {
              // console.log('value :>> ', value);
            },
          },
        },
        yAxis: {
          max: props.max,
          min: props.min,
          show: false,
          scale: true,
          // offset: -5,
          splitLine: { show: false },
        },
        series:
          {
            type: 'line',
            symbol: props.symbol,
            symbolSize: 7,
            lineStyle: {
              color: '#FF7984',
            },
            itemStyle: {
              color: '#FF7984',
            },
            data: props.series,
          } || [],
      };
      setTimeout(() => {
        lineChart.value.hideLoading();
        lineChart.value.setOption(option, true);
      }, 200);
    };
    watch(
      () => props.series,
      (newValue) => {
        nextTick(() => {
          if (newValue && lineChart.value) {
            initChart();
            resizeChart();
          }
        });
      },
      { deep: true, immediate: true }
    );
    const isMobile = ref(false);
    onMounted(() => {
      const screen = window.innerWidth;
      if (screen < 768) {
        isMobile.value = true;
      }
      if (lineChart.value) {
        lineChart.value.dispose();
        return;
      }
      initChartDom();
      // console.log('props.series :>> ', props.series);
    });
  },
};
</script>
