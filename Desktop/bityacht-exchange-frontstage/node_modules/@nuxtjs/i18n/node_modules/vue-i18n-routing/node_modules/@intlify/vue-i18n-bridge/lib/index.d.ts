/* eslint-disable @typescript-eslint/no-explicit-any */
declare type Vue = any // emulate Vue (v2)
declare namespace _VueI18n {
  type Path = string
  type Locale = string
  type FallbackLocale = string | string[] | false | { [locale: string]: string[] }
  type Values = any[] | { [key: string]: any }
  type Choice = number
  interface MessageContext {
    list(index: number): unknown
    named(key: string): unknown
    linked(key: string): TranslateResult
    values: any
    path: string
    formatter: Formatter
    messages: LocaleMessages
    locale: Locale
  }
  type MessageFunction = (ctx: MessageContext) => string
  type LocaleMessage = string | MessageFunction | LocaleMessageObject | LocaleMessageArray
  interface LocaleMessageObject {
    [key: string]: LocaleMessage
  }
  interface LocaleMessageArray {
    [index: number]: LocaleMessage
  }
  interface LocaleMessages {
    [key: string]: LocaleMessageObject
  }
  type TranslateResult = string | LocaleMessages

  type LocaleMatcher = 'lookup' | 'best fit'
  type FormatMatcher = 'basic' | 'best fit'

  type DateTimeHumanReadable = 'long' | 'short' | 'narrow'
  type DateTimeDigital = 'numeric' | '2-digit'

  interface SpecificDateTimeFormatOptions extends Intl.DateTimeFormatOptions {
    year?: DateTimeDigital
    month?: DateTimeDigital | DateTimeHumanReadable
    day?: DateTimeDigital
    hour?: DateTimeDigital
    minute?: DateTimeDigital
    second?: DateTimeDigital
    weekday?: DateTimeHumanReadable
    era?: DateTimeHumanReadable
    timeZoneName?: 'long' | 'short'
    localeMatcher?: LocaleMatcher
    formatMatcher?: FormatMatcher
  }

  type DateTimeFormatOptions = Intl.DateTimeFormatOptions | SpecificDateTimeFormatOptions

  interface DateTimeFormat {
    [key: string]: DateTimeFormatOptions
  }
  interface DateTimeFormats {
    [locale: string]: DateTimeFormat
  }
  type DateTimeFormatResult = string

  type CurrencyDisplay = 'symbol' | 'code' | 'name'

  interface SpecificNumberFormatOptions extends Intl.NumberFormatOptions {
    style?: 'decimal' | 'percent'
    currency?: string
    currencyDisplay?: CurrencyDisplay
    localeMatcher?: LocaleMatcher
    formatMatcher?: FormatMatcher
  }

  interface CurrencyNumberFormatOptions extends Intl.NumberFormatOptions {
    style: 'currency'
    currency: string
    currencyDisplay?: CurrencyDisplay
    localeMatcher?: LocaleMatcher
    formatMatcher?: FormatMatcher
  }

  type NumberFormatOptions = Intl.NumberFormatOptions | SpecificNumberFormatOptions | CurrencyNumberFormatOptions

  interface NumberFormat {
    [key: string]: NumberFormatOptions
  }
  interface NumberFormats {
    [locale: string]: NumberFormat
  }
  type NumberFormatResult = string
  type PluralizationRulesMap = {
    [lang: string]: (choice: number, choicesLength: number) => number
  }
  type Modifiers = { [key: string]: (str: string) => string }

  type FormattedNumberPartType =
    | 'currency'
    | 'decimal'
    | 'fraction'
    | 'group'
    | 'infinity'
    | 'integer'
    | 'literal'
    | 'minusSign'
    | 'nan'
    | 'plusSign'
    | 'percentSign'

  type WarnHtmlInMessageLevel = 'off' | 'warn' | 'error'

  interface FormattedNumberPart {
    type: FormattedNumberPartType
    value: string
  }
  interface NumberFormatToPartsResult {
    [index: number]: FormattedNumberPart
  }

  interface Formatter {
    interpolate(message: string, values: Values | undefined, path: string): any[] | null
  }

  type MissingHandler = (locale: Locale, key: Path, vm: Vue | null, values: any) => string | void
  type PostTranslationHandler = (str: string, key?: string) => string
  type ComponentInstanceCreatedListener = (newVm: VueI18n, rootVm: VueI18n) => void

  interface IntlAvailability {
    dateTimeFormat: boolean
    numberFormat: boolean
  }

  interface I18nOptions {
    locale?: Locale
    fallbackLocale?: FallbackLocale
    messages?: LocaleMessages
    dateTimeFormats?: DateTimeFormats
    numberFormats?: NumberFormats
    formatter?: Formatter
    modifiers?: Modifiers
    missing?: MissingHandler
    fallbackRoot?: boolean
    formatFallbackMessages?: boolean
    sync?: boolean
    silentTranslationWarn?: boolean | RegExp
    silentFallbackWarn?: boolean | RegExp
    preserveDirectiveContent?: boolean
    pluralizationRules?: PluralizationRulesMap
    warnHtmlInMessage?: WarnHtmlInMessageLevel
    sharedMessages?: LocaleMessages
    postTranslation?: PostTranslationHandler
    componentInstanceCreatedListener?: ComponentInstanceCreatedListener
    escapeParameterHtml?: boolean
  }
}

// Stub vue-i18n class interfaces
declare interface VueI18n {
  readonly messages: _VueI18n.LocaleMessages
  readonly dateTimeFormats: _VueI18n.DateTimeFormats
  readonly numberFormats: _VueI18n.NumberFormats

  locale: _VueI18n.Locale
  fallbackLocale: _VueI18n.FallbackLocale
  missing: _VueI18n.MissingHandler
  formatter: _VueI18n.Formatter
  formatFallbackMessages: boolean
  silentTranslationWarn: boolean | RegExp
  silentFallbackWarn: boolean | RegExp
  preserveDirectiveContent: boolean
  sync: boolean
  pluralizationRules: _VueI18n.PluralizationRulesMap
  warnHtmlInMessage: _VueI18n.WarnHtmlInMessageLevel
  postTranslation: _VueI18n.PostTranslationHandler
  t(key: _VueI18n.Path, values?: _VueI18n.Values): _VueI18n.TranslateResult
  t(key: _VueI18n.Path, locale: _VueI18n.Locale, values?: _VueI18n.Values): _VueI18n.TranslateResult
  tc(key: _VueI18n.Path, choice?: _VueI18n.Choice, values?: _VueI18n.Values): string
  tc(key: _VueI18n.Path, choice: _VueI18n.Choice, locale: _VueI18n.Locale, values?: _VueI18n.Values): string
  te(key: _VueI18n.Path, locale?: _VueI18n.Locale): boolean
  d(value: number | Date, key?: _VueI18n.Path, locale?: _VueI18n.Locale): _VueI18n.DateTimeFormatResult
  d(value: number | Date, args?: { [key: string]: string }): _VueI18n.DateTimeFormatResult
  d(value: number | Date, options?: _VueI18n.DateTimeFormatOptions): _VueI18n.DateTimeFormatResult
  n(value: number, key?: _VueI18n.Path, locale?: _VueI18n.Locale): _VueI18n.NumberFormatResult
  n(value: number, args?: { [key: string]: string }): _VueI18n.NumberFormatResult
  n(value: number, options?: _VueI18n.NumberFormatOptions, locale?: _VueI18n.Locale): _VueI18n.NumberFormatResult
  getLocaleMessage(locale: _VueI18n.Locale): _VueI18n.LocaleMessageObject
  setLocaleMessage(locale: _VueI18n.Locale, message: _VueI18n.LocaleMessageObject): void
  mergeLocaleMessage(locale: _VueI18n.Locale, message: _VueI18n.LocaleMessageObject): void
  getDateTimeFormat(locale: _VueI18n.Locale): _VueI18n.DateTimeFormat
  setDateTimeFormat(locale: _VueI18n.Locale, format: _VueI18n.DateTimeFormat): void
  mergeDateTimeFormat(locale: _VueI18n.Locale, format: _VueI18n.DateTimeFormat): void
  getNumberFormat(locale: _VueI18n.Locale): _VueI18n.NumberFormat
  setNumberFormat(locale: _VueI18n.Locale, format: _VueI18n.NumberFormat): void
  mergeNumberFormat(locale: _VueI18n.Locale, format: _VueI18n.NumberFormat): void
  getChoiceIndex: (choice: number, choicesLength: number) => number
}
declare const isVueI18n8: boolean
declare const isVueI18n9: boolean

declare class VueI18nLegacy {
  static install()
  static version
}

export * from 'vue-i18n'
export { isVueI18n8, isVueI18n9 }
export default VueI18nLegacy

/* eslint-enable @typescript-eslint/no-explicit-any */
