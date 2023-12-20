/** @type {import('tailwindcss').Config} */

module.exports = {
  content: [`./components/**/*.{vue,js,ts}`, `./layouts/**/*.vue`, `./pages/**/*.vue`, `./composables/**/*.{js,ts}`, `./plugins/**/*.{js,ts}`, `./utils/**/*.{js,ts}`, `./App.{js,ts,vue}`, `./app.{js,ts,vue}`, `./Error.{js,ts,vue}`, `./error.{js,ts,vue}`, `./app.config.{js,ts}`],
  theme: {
    fontFamily: {
      sans: `PingFang TC, 'Lato', 'Noto Sans TC', 'Noto Sans', "微軟正黑體", "Open Sans", Roboto, Arial, sans-serif`,
    },
    fontSize: {
      h1: [
        '50px',
        {
          fontWeight: '700',
        },
      ],
      h2: [
        '44px',
        {
          fontWeight: '700',
        },
      ],
      h3: [
        '40px',
        {
          fontWeight: '700',
        },
      ],
      h4: [
        '36px',
        {
          fontWeight: '700',
        },
      ],
      h5: [
        '32px',
        {
          fontWeight: '700',
        },
      ],
      subTitle: ['18px'],
      normal: ['16px'],
      subText: ['14px'],
      small: ['12px'],
    },
    extend: {
      screens: {
        xxxs: '400px',
        xxs: '460px',
        xs: '576px',
        // 'sm': '640px',
        // 'md': '768px',
        md_lg: '900px',
        // 'lg': '1024px',
        // 'xl': '1280px',
        // '2xl': '1536px',
        '3xl': '1700px',
      },
      colors: {
        white: '#ffffff',
        black: '#000000',
        darkBlueBg: '#060A2F',
        grayBorder: '#c5c6c9',
        grayText: '#a4a4a4',
        grayBg: '#d9d9d9 ',
        gray_50: '#fefffe',
        gray_100: '#FCFDFD',
        gray_150: '#F9FAFB',
        gray_200: '#F4F7FA',
        gray_300: '#DDE1E5',
        gray_400: '#8B939A',
        gray_500: '#8A8C8F',
        gray_600: '#6B6C6C',
        gray_700: '#181B22',
        gray_800: '#7589A4',
        shallow_gray: '#c8c8c8',
        primary: '#19253A',
        primary_400: '#394B6A',
        secondary: '#8E9BA5',
        waterBlueOld: '#25BBEE',
        // waterBlue: '#19253a', // old dark blue
        waterBlue: '#7484a2',
        darkBlue: '#13458c',
        red: '#DD4949',
        lightRed: '#FF0303',
        green: '#44CF2E',
        leaf: '#CDCD83',
        yellow: '#F8b62d',
        hover: '#18e1be',
        skyBlue: '#179CD6',
        skyBlue2: '#e1eff8',
        pink: '#F9A7BA',
        pink2: '#fea38f',
        orange: '#ffa412',
        brown: '#DA8353',
      },
    },
  },
  plugins: [],
};
