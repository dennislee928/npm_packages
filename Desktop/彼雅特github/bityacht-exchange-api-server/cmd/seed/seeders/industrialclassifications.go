package seeders

import (
	"bityacht-exchange-api-server/cmd/migrate/migrations"
	"bityacht-exchange-api-server/cmd/seed"

	"gorm.io/gorm"
)

func init() {
	seed.Register(&SeederIndustrialClassifications{})
}

// SeederIndustrialClassifications is a kind of ISeeder, You can set it as same as model.
type SeederIndustrialClassifications migrations.Migration1689660236

// SeederName for ISeeder
func (*SeederIndustrialClassifications) SeederName() string {
	return "IndustrialClassifications"
}

// TableName for gorm
func (*SeederIndustrialClassifications) TableName() string {
	return (*migrations.Migration1689660236)(nil).TableName()
}

// Default for ISeeder
func (*SeederIndustrialClassifications) Default(db *gorm.DB) error {
	records := []SeederIndustrialClassifications{
		{Code: "P", Chinese: "學術/教育機構", English: "Academy Institutions/Education Institutions"},
		{Code: "A", Chinese: "農林漁牧業", English: "Agriculture, forestry, fishing, and animal husbandry industries"},
		{Code: "G4853", Chinese: "古董與拍賣行業", English: "Antique and Auction Industry"},
		{Code: "C3399", Chinese: "軍火製造業", English: "Arms Industry"},
		{Code: "R", Chinese: "藝術、娛樂及休閒服務", English: "Arts, Entertainment and Recreation"},
		{Code: "", Chinese: "記帳士", English: "Certified Bookkeeper"},
		{Code: "", Chinese: "公務員", English: "Civil Servant"},
		{Code: "F", Chinese: "營建業", English: "Construction Industry"},
		{Code: "", Chinese: "大使館或領事館", English: "Consular Officer"},
		{Code: "", Chinese: "顧問服務業", English: "Consultancy Service"},
		{Code: "", Chinese: "電子支付或第三方支付機構", English: "Eletronic Payment or Third-Party Payment Industries"},
		{Code: "", Chinese: "能源業", English: "Energy Industry"},
		{Code: "", Chinese: "住宿及餐飲產業", English: "Entertainment, accommodation and catering Industries"},
		{Code: "", Chinese: "環境保護工程業", English: "Environmental Protection Works Industry"},
		{Code: "", Chinese: "民間融資性租賃業務", English: "Financial Leasing Industry"},
		{Code: "", Chinese: "金融服務業", English: "Financial Services Industry"},
		{Code: "R92", Chinese: "博弈產業", English: "Gambling Industry"},
		{Code: "Q86", Chinese: "醫療機構", English: "Health Facilities"},
		{Code: "", Chinese: "高科技業", English: "Hightech Industry"},
		{Code: "", Chinese: "國際貿易業", English: "International Trade"},
		{Code: "", Chinese: "珠寶商、銀樓業、貴金屬產業", English: "Jewellery and Precious Metals Industry"},
		{Code: "", Chinese: "地政士", English: "Judicial scrivener"},
		{Code: "M69", Chinese: "法律、會計服務業", English: "Legal and Accounting Services"},
		{Code: "", Chinese: "運輸、倉儲業", English: "Logistics and Warehouse Industry"},
		{Code: "C", Chinese: "製造業", English: "Manufacturing Industry"},
		{Code: "", Chinese: "媒體傳播業", English: "Media and Communication Industry"},
		{Code: "B", Chinese: "礦業及土石開採業", English: "Mining and Quarrying Industry"},
		{Code: "", Chinese: "非營利組織", English: "Non-Profit Organizations"},
		{Code: "", Chinese: "公證人", English: "Notary"},
		{Code: "", Chinese: "核子能源", English: "Nuclear Industry"},
		{Code: "S9491", Chinese: "政治團體工作人員", English: "Personnel of Political Group"},
		{Code: "S941", Chinese: "宗教團體相關工作人員", English: "Personnel of Religion"},
		{Code: "L", Chinese: "不動產相關行業", English: "Real estate related industries"},
		{Code: "", Chinese: "特殊娛樂產業（不包含一般KTV）", English: "Regulated Industries(excluding general KTV)"},
		{Code: "", Chinese: "退休", English: "Retired"},
		{Code: "", Chinese: "學生", English: "Student"},
		{Code: "", Chinese: "報稅代理人", English: "Tax Agent"},
		{Code: "", Chinese: "通訊與資訊安全行業", English: "Telecommunications and Cyber Security Industry"},
		{Code: "", Chinese: "旅行及相關代訂服務業", English: "Travel and Booking Services Industry"},
		{Code: "", Chinese: "未受雇", English: "Unemployed"},
		{Code: "", Chinese: "虛擬通貨相關業者", English: "Virtual Asset Service Provider Industry"},
		{Code: "", Chinese: "批發及零售業", English: "Wholesale or Retail Industries"},
		{Code: "", Chinese: "汽車貿易及零售業", English: "Wholesale or Retail Trade of Automobiles"},
		{Code: "", Chinese: "其他", English: "Others"},
	}

	return db.Create(&records).Error
}

// Fake for ISeeder
func (*SeederIndustrialClassifications) Fake(db *gorm.DB) error {
	return seed.ErrNotImplement
}
