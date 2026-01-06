// Common cities data
export interface CityOption {
  code: string
  name: string
  pinyin?: string
}

export const CITIES: CityOption[] = [
  { code: 'beijing', name: '北京', pinyin: 'beijing' },
  { code: 'shanghai', name: '上海', pinyin: 'shanghai' },
  { code: 'guangzhou', name: '广州', pinyin: 'guangzhou' },
  { code: 'shenzhen', name: '深圳', pinyin: 'shenzhen' },
  { code: 'hangzhou', name: '杭州', pinyin: 'hangzhou' },
  { code: 'chengdu', name: '成都', pinyin: 'chengdu' },
  { code: 'wuhan', name: '武汉', pinyin: 'wuhan' },
  { code: 'nanjing', name: '南京', pinyin: 'nanjing' },
  { code: 'xian', name: '西安', pinyin: 'xian' },
  { code: 'chongqing', name: '重庆', pinyin: 'chongqing' },
]

export const getCityName = (code: string): string => {
  const city = CITIES.find((c) => c.code === code)
  return city?.name || code
}

