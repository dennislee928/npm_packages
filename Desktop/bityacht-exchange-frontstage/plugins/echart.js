// import { use } from 'echarts/core';

// import { CanvasRenderer } from 'echarts/renderers';
// import { BarChart } from 'echarts/charts';
// import { GridComponent, TooltipComponent } from 'echarts/components';
import * as echarts from 'echarts';

export default defineNuxtPlugin((nuxtApp) => {
  // use([CanvasRenderer, BarChart, GridComponent, TooltipComponent]);
  nuxtApp.vueApp.config.globalProperties.$echarts = echarts;
});
