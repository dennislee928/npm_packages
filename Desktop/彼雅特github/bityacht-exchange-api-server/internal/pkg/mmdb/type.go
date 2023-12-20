package mmdb

type I18n struct {
	EN   string `maxminddb:"en"`
	ZhCN string `maxminddb:"zh-CN"`
}

type City struct {
	Names I18n `maxminddb:"names"`
}

type Continent struct {
	Code  string `maxminddb:"code"`
	Names I18n   `maxminddb:"names"`
}

type Country struct {
	ISOCode string `maxminddb:"iso_code"`
	Names   I18n   `maxminddb:"names"`
}

type ISP struct {
	Names I18n `maxminddb:"names"`
}

type Province struct {
	Names I18n `maxminddb:"names"`
}

type CityResult struct {
	City      City      `maxminddb:"city"`
	Continent Continent `maxminddb:"continent"`
	Country   Country   `maxminddb:"country"`
	ISP       ISP       `maxminddb:"isp"`
	Province  Province  `maxminddb:"province"`
}

func (cr CityResult) String() string {
	if cr.City.Names.EN != "" {
		return cr.Country.Names.EN + " (" + cr.City.Names.EN + ")"
	}
	if cr.Continent.Names.EN != "" {
		return cr.Continent.Names.EN
	}

	return "Unknown"
}
