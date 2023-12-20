package seeders

import (
	"bityacht-exchange-api-server/cmd/migrate/migrations"
	"bityacht-exchange-api-server/cmd/seed"

	"gorm.io/gorm"
)

func init() {
	seed.Register(&SeederCountries{})
}

// SeederCountries is a kind of ISeeder, You can set it as same as model.
type SeederCountries migrations.Migration1689659690

// SeederName for ISeeder
func (*SeederCountries) SeederName() string {
	return "Countries"
}

// SeederName for ISeeder
func (s *SeederCountries) TableName() string {
	return (*migrations.Migration1689659690)(nil).TableName()
}

const countryCodeTaiwan = "TWN"

// Default for ISeeder
func (*SeederCountries) Default(db *gorm.DB) error {
	records := []SeederCountries{
		{Code: "AFG", Chinese: "阿富汗", English: "Afghanistan", Locale: "ar"},
		{Code: "ALB", Chinese: "阿爾巴尼亞", English: "Albania"},
		{Code: "DZA", Chinese: "阿爾及利亞", English: "Algeria"},
		{Code: "AND", Chinese: "安道爾", English: "Andorra"},
		{Code: "AGO", Chinese: "安哥拉", English: "Angola"},
		{Code: "ATG", Chinese: "安地卡與巴布達", English: "Antigua and Barbuda"},
		{Code: "ARG", Chinese: "阿根廷", English: "Argentina"},
		{Code: "ARM", Chinese: "亞美尼亞", English: "Armenia"},
		{Code: "AUS", Chinese: "澳大利亞", English: "Australia"},
		{Code: "AZE", Chinese: "亞塞拜然", English: "Azerbaijan"},
		{Code: "BHS", Chinese: "巴哈馬", English: "Bahamas"},
		{Code: "BHR", Chinese: "巴林", English: "Bahrain"},
		{Code: "BGD", Chinese: "孟加拉", English: "Bangladesh"},
		{Code: "BRB", Chinese: "巴貝多", English: "Barbados"},
		{Code: "BLR", Chinese: "白俄羅斯", English: "Belarus"},
		{Code: "BLZ", Chinese: "貝里斯", English: "Belize"},
		{Code: "BEN", Chinese: "貝南", English: "Benin"},
		{Code: "BTN", Chinese: "不丹", English: "Bhutan"},
		{Code: "BOL", Chinese: "玻利維亞", English: "Bolivia"},
		{Code: "BIH", Chinese: "波士尼亞與赫塞哥維納", English: "Bosnia and Herzegovina"},
		{Code: "BWA", Chinese: "波札那", English: "Botswana"},
		{Code: "BRA", Chinese: "巴西", English: "Brazil"},
		{Code: "BRN", Chinese: "汶萊", English: "Brunei"},
		{Code: "BFA", Chinese: "布吉納法索", English: "Burkina Faso"},
		{Code: "BDI", Chinese: "蒲隆地", English: "Burundi"},
		{Code: "KHM", Chinese: "柬埔寨", English: "Cambodia"},
		{Code: "CMR", Chinese: "喀麥隆", English: "Cameroon"},
		{Code: "CAN", Chinese: "加拿大", English: "Canada"},
		{Code: "CPV", Chinese: "維德角", English: "Cape Verde"},
		{Code: "CAF", Chinese: "中非", English: "Central African"},
		{Code: "TCD", Chinese: "查德", English: "Chad"},
		{Code: "CHL", Chinese: "智利", English: "Chile"},
		{Code: "COL", Chinese: "哥倫比亞", English: "Colombia"},
		{Code: "COM", Chinese: "葛摩", English: "Comoros"},
		{Code: "COG", Chinese: "剛果", English: "Congo"},
		{Code: "COK", Chinese: "庫克群島", English: "Cook Islands"},
		{Code: "CRI", Chinese: "哥斯大黎加", English: "Costa Rica"},
		{Code: "CIV", Chinese: "象牙海岸", English: "Cote d'Ivoire"},
		{Code: "CUB", Chinese: "古巴", English: "Cuba"},
		{Code: "DJI", Chinese: "吉布地", English: "Djibouti"},
		{Code: "DMA", Chinese: "多米尼克", English: "Dominica"},
		{Code: "DOM", Chinese: "多明尼加", English: "Dominican Republic"},
		{Code: "COD", Chinese: "民主剛果", English: "DRC、Democratic Congo"},
		{Code: "TLS", Chinese: "東帝汶", English: "East Timor、Timor-Leste"},
		{Code: "ECU", Chinese: "厄瓜多", English: "Ecuador"},
		{Code: "EGY", Chinese: "埃及", English: "Egypt"},
		{Code: "SLV", Chinese: "薩爾瓦多", English: "El Salvador"},
		{Code: "GNQ", Chinese: "赤道幾內亞", English: "Equatorial Guinea"},
		{Code: "ERI", Chinese: "厄利垂亞", English: "Eritrea"},
		{Code: "SWZ", Chinese: "史瓦帝尼", English: "Eswatini"},
		{Code: "ETH", Chinese: "衣索比亞", English: "Ethiopia"},
		{Code: "FJI", Chinese: "斐濟", English: "Fiji"},
		{Code: "GAB", Chinese: "加彭", English: "Gabon"},
		{Code: "GMB", Chinese: "甘比亞", English: "Gambia"},
		{Code: "GEO", Chinese: "喬治亞", English: "Georgia"},
		{Code: "GHA", Chinese: "迦納", English: "Ghana"},
		{Code: "GRD", Chinese: "格瑞那達", English: "Grenada"},
		{Code: "GTM", Chinese: "瓜地馬拉", English: "Guatemala"},
		{Code: "GIN", Chinese: "幾內亞", English: "Guinea"},
		{Code: "GNB", Chinese: "幾內亞比索", English: "Guinea-Bissau"},
		{Code: "GUY", Chinese: "蓋亞那", English: "Guyana"},
		{Code: "HTI", Chinese: "海地", English: "Haiti"},
		{Code: "VAT", Chinese: "梵蒂岡", English: "Holy See"},
		{Code: "HND", Chinese: "宏都拉斯", English: "Honduras"},
		{Code: "ISL", Chinese: "冰島", English: "Iceland"},
		{Code: "IND", Chinese: "印度", English: "India"},
		{Code: "IDN", Chinese: "印尼", English: "Indonesia", Locale: "id"},
		{Code: "IRQ", Chinese: "伊拉克", English: "Iraq"},
		{Code: "ISR", Chinese: "以色列", English: "Israel"},
		{Code: "JAM", Chinese: "牙買加", English: "Jamaica"},
		{Code: "JOR", Chinese: "約旦", English: "Jordan"},
		{Code: "KAZ", Chinese: "哈薩克", English: "Kazakhstan"},
		{Code: "KEN", Chinese: "肯亞", English: "Kenya"},
		{Code: "KIR", Chinese: "吉里巴斯", English: "Kiribati"},
		{Code: "UNK", Chinese: "科索沃", English: "Kosovo"},
		{Code: "KWT", Chinese: "科威特", English: "Kuwait"},
		{Code: "KGZ", Chinese: "吉爾吉斯", English: "Kyrgyzstan"},
		{Code: "LAO", Chinese: "寮國", English: "Laos"},
		{Code: "LBN", Chinese: "黎巴嫩", English: "Lebanon"},
		{Code: "LSO", Chinese: "賴索托", English: "Lesotho"},
		{Code: "LBR", Chinese: "賴比瑞亞", English: "Liberia"},
		{Code: "LBY", Chinese: "利比亞", English: "Libya"},
		{Code: "LIE", Chinese: "列支敦斯登", English: "Liechtenstein"},
		{Code: "MDG", Chinese: "馬達加斯加", English: "Madagascar"},
		{Code: "MWI", Chinese: "馬拉威", English: "Malawi"},
		{Code: "MYS", Chinese: "馬來西亞", English: "Malaysia"},
		{Code: "MDV", Chinese: "馬爾地夫", English: "Maldives"},
		{Code: "MLI", Chinese: "馬利", English: "Mali"},
		{Code: "MHL", Chinese: "馬紹爾", English: "Marshall Islands"},
		{Code: "MRT", Chinese: "茅利塔尼亞", English: "Mauritania"},
		{Code: "MUS", Chinese: "模里西斯", English: "Mauritius"},
		{Code: "MEX", Chinese: "墨西哥", English: "Mexico"},
		{Code: "FSM", Chinese: "密克羅尼西亞", English: "Micronesia"},
		{Code: "MDA", Chinese: "摩爾多瓦", English: "Moldova"},
		{Code: "MCO", Chinese: "摩納哥", English: "Monaco"},
		{Code: "MNG", Chinese: "蒙古", English: "Mongolia"},
		{Code: "MNE", Chinese: "蒙特內哥羅", English: "Montenegro"},
		{Code: "MAR", Chinese: "摩洛哥", English: "Morocco"},
		{Code: "MOZ", Chinese: "莫桑比克", English: "Mozambique"},
		{Code: "NAM", Chinese: "納米比亞", English: "Namibia"},
		{Code: "NRU", Chinese: "諾魯", English: "Nauru"},
		{Code: "NPL", Chinese: "尼泊爾", English: "Nepal"},
		{Code: "NZL", Chinese: "紐西蘭", English: "New Zealand"},
		{Code: "NIC", Chinese: "尼加拉瓜", English: "Nicaragua"},
		{Code: "NER", Chinese: "尼日", English: "Niger"},
		{Code: "NGA", Chinese: "奈及利亞", English: "Nigeria"},
		{Code: "NIU", Chinese: "紐埃", English: "Niue"},
		{Code: "MKD", Chinese: "北馬其頓", English: "North Macedonia"},
		{Code: "NOR", Chinese: "挪威", English: "Norway"},
		{Code: "OMN", Chinese: "阿曼", English: "Oman"},
		{Code: "PAK", Chinese: "巴基斯坦", English: "Pakistan"},
		{Code: "PLW", Chinese: "帛琉", English: "Palau"},
		{Code: "PAN", Chinese: "巴拿馬", English: "Panama"},
		{Code: "PNG", Chinese: "巴布亞紐幾內亞、巴紐", English: "Papua New Guinea、PNG"},
		{Code: "PRY", Chinese: "巴拉圭", English: "Paraguay"},
		{Code: "PER", Chinese: "秘魯", English: "Peru"},
		{Code: "PHL", Chinese: "菲律賓", English: "Philippines"},
		{Code: "QAT", Chinese: "卡達", English: "Qatar"},
		{Code: "RWA", Chinese: "盧安達", English: "Rwanda"},
		{Code: "WSM", Chinese: "薩摩亞", English: "Samoa"},
		{Code: "SMR", Chinese: "聖馬利諾", English: "San Marino"},
		{Code: "STP", Chinese: "聖多美", English: "Sao Tome and Principe"},
		{Code: "SAU", Chinese: "沙烏地阿拉伯、沙烏地", English: "Saudi Arabia"},
		{Code: "SEN", Chinese: "塞內加爾", English: "Senegal"},
		{Code: "SRB", Chinese: "塞爾維亞", English: "Serbia"},
		{Code: "SYC", Chinese: "塞席爾", English: "Seychelles"},
		{Code: "SLE", Chinese: "獅子山", English: "Sierra Leone"},
		{Code: "SGP", Chinese: "新加坡", English: "Singapore"},
		{Code: "SLB", Chinese: "索羅門", English: "Solomon Islands"},
		{Code: "SOM", Chinese: "索馬利亞", English: "Somalia"},
		{Code: "ZAF", Chinese: "南非", English: "South Africa"},
		{Code: "SSD", Chinese: "南蘇丹", English: "South Sudan"},
		{Code: "LKA", Chinese: "斯里蘭卡", English: "Sri Lanka"},
		{Code: "KNA", Chinese: "聖克里斯多福", English: "St. Kitts and Nevis"},
		{Code: "LCA", Chinese: "聖露西亞", English: "St. Lucia"},
		{Code: "VCT", Chinese: "聖文森國", English: "St. Vincent & the Grenadines"},
		{Code: "SDN", Chinese: "蘇丹", English: "Sudan"},
		{Code: "SUR", Chinese: "蘇利南", English: "Suriname"},
		{Code: "CHE", Chinese: "瑞士", English: "Switzerland"},
		{Code: "SYR", Chinese: "敘利亞", English: "Syria"},
		{Code: countryCodeTaiwan, Chinese: "臺灣", English: "Taiwan", Locale: "zh-HK"},
		{Code: "TJK", Chinese: "塔吉克", English: "Tajikistan"},
		{Code: "TZA", Chinese: "坦尚尼亞", English: "Tanzania"},
		{Code: "THA", Chinese: "泰國", English: "Thailand", Locale: "th"},
		{Code: "TGO", Chinese: "多哥", English: "Togo"},
		{Code: "TON", Chinese: "東加", English: "Tonga"},
		{Code: "TTO", Chinese: "千里達", English: "Trinidad and Tobago"},
		{Code: "TUN", Chinese: "突尼西亞", English: "Tunisia"},
		{Code: "TUR", Chinese: "土耳其", English: "Türkiye"},
		{Code: "TKM", Chinese: "土庫曼", English: "Turkmenistan"},
		{Code: "TUV", Chinese: "吐瓦魯", English: "Tuvalu"},
		{Code: "UGA", Chinese: "烏干達", English: "Uganda"},
		{Code: "UKR", Chinese: "烏克蘭", English: "Ukraine"},
		{Code: "ARE", Chinese: "阿聯", English: "United Arab Emirates"},
		{Code: "GBR", Chinese: "英國", English: "United Kingdom、UK", Locale: "en-GB"},
		{Code: "URY", Chinese: "烏拉圭", English: "Uruguay"},
		{Code: "UZB", Chinese: "烏茲別克", English: "Uzbekistan"},
		{Code: "VUT", Chinese: "萬那杜", English: "Vanuatu"},
		{Code: "VEN", Chinese: "委內瑞拉", English: "Venezuela"},
		{Code: "VNM", Chinese: "越南", English: "Vietnam", Locale: "vi"},
		{Code: "YEM", Chinese: "葉門", English: "Yemen"},
		{Code: "ZMB", Chinese: "尚比亞", English: "Zambia"},
		{Code: "ZWE", Chinese: "辛巴威", English: "Zimbabwe"},
	}

	return db.Create(&records).Error
}

// Fake for ISeeder
func (*SeederCountries) Fake(db *gorm.DB) error {
	return seed.ErrNotImplement
}
