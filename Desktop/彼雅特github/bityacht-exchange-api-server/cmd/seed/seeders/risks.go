package seeders

import (
	"gorm.io/gorm"

	"bityacht-exchange-api-server/cmd/migrate/migrations"
	"bityacht-exchange-api-server/cmd/seed"
)

func init() {
	seed.Register(&SeederRisks{})
}

// SeederUsers is a kind of ISeeder, You can set it as same as model.
type SeederRisks migrations.Migration1690884579

// SeederName for ISeeder
func (*SeederRisks) SeederName() string {
	return "Risks"
}

// TableName for gorm
func (*SeederRisks) TableName() string {
	return (*migrations.Migration1690884579)(nil).TableName()
}

const A = "A.用戶因素"
const Aa = "a.姓名檢核"
const Ab = "b.曾犯特定罪名，經起訴者"
const Ac = "c.行業別"
const Ad = "d.職稱"
const Ae = "e.註冊異常"
const Af = "f.資料異常"
const Ag = "g.電話照會異常"

const B = "B.國家或地區因素"
const Ba = "a.高法遵成本或禁止虛擬通貨交易國家"
const Bb = "b.制裁國家"
const Bc = "c.其他未遵循或未充分遵循國際防制洗錢組織建議之國家或地區"
const Bd = "d.風險國家"
const Be = "e.雙重國籍與單一外國籍"

const C = "C.產品、服務、交易或支付因素"
const Ca = "a.使用產品、服務目的或方式"
const Cb = "b.累計交易額"
const Cc = "c.黑名單相關與用戶交易監控"
const Cd = "d.支付行為"

// Default for ISeeder
func (*SeederRisks) Default(db *gorm.DB) error {
	records := []SeederRisks{
		{Factor: A, SubFactor: Aa, Detail: "姓名檢核命中制裁名單對象", Score: 12},
		{Factor: A, SubFactor: Aa, Detail: "現任國內外PEPs", Score: 8},
		{Factor: A, SubFactor: Aa, Detail: "現任國內外RCA", Score: 8},
		{Factor: A, SubFactor: Aa, Detail: "離任國內外PEPs及RCA，判斷仍具影響力者", Score: 8},
		{Factor: A, SubFactor: Aa, Detail: "離任國內外PEPs及RCA", Score: 4},
		{Factor: A, SubFactor: Aa, Detail: "命中負面新聞名單者，且確定為同一人者", Score: 4},
		{Factor: A, SubFactor: Ab, Detail: "洗錢", Score: 12},
		{Factor: A, SubFactor: Ab, Detail: "資恐犯罪", Score: 12},
		{Factor: A, SubFactor: Ab, Detail: "性剝削（包含兒童）", Score: 12},
		{Factor: A, SubFactor: Ab, Detail: "人口販賣", Score: 12},
		{Factor: A, SubFactor: Ab, Detail: "貪汙賄賂", Score: 12},
		{Factor: A, SubFactor: Ab, Detail: "毒品", Score: 12},
		{Factor: A, SubFactor: Ab, Detail: "詐欺", Score: 12},
		{Factor: A, SubFactor: Ab, Detail: "賭博", Score: 8},
		{Factor: A, SubFactor: Ab, Detail: "組織犯罪", Score: 8},
		{Factor: A, SubFactor: Ab, Detail: "走私", Score: 8},
		{Factor: A, SubFactor: Ab, Detail: "稅務犯罪", Score: 8},
		{Factor: A, SubFactor: Ab, Detail: "證券犯罪", Score: 8},
		{Factor: A, SubFactor: Ab, Detail: "背信掏空犯罪", Score: 8},
		{Factor: A, SubFactor: Ab, Detail: "非法軍火交易", Score: 8},
		{Factor: A, SubFactor: Ab, Detail: "環境犯罪", Score: 8},
		{Factor: A, SubFactor: Ab, Detail: "智慧財產犯罪", Score: 8},
		{Factor: A, SubFactor: Ab, Detail: "無以上罪行，但經評估需納入風險考量", Score: 8},
		{Factor: A, SubFactor: Ac, Detail: "特殊娛樂產業（如舞廳、酒吧、視聽歌唱、夜店及三溫暖，但不包含一般KTV）", Score: 12},
		{Factor: A, SubFactor: Ac, Detail: "博弈產業", Score: 8},
		{Factor: A, SubFactor: Ac, Detail: "國防武器或戰爭設備等軍火製造相關", Score: 8},
		{Factor: A, SubFactor: Ac, Detail: "核子能源", Score: 8},
		{Factor: A, SubFactor: Ac, Detail: "大使館或領事館", Score: 8},
		{Factor: A, SubFactor: Ac, Detail: "非營利組織", Score: 4},
		{Factor: A, SubFactor: Ac, Detail: "金融服務業", Score: 4},
		{Factor: A, SubFactor: Ac, Detail: "虛擬通貨相關業者", Score: 4},
		{Factor: A, SubFactor: Ac, Detail: "珠寶商、銀樓業、古董業", Score: 4},
		{Factor: A, SubFactor: Ac, Detail: "礦業及土石開採業", Score: 4},
		{Factor: A, SubFactor: Ac, Detail: "旅行及相關代訂服務業", Score: 4},
		{Factor: A, SubFactor: Ac, Detail: "汽車貿易及零售業", Score: 4},
		{Factor: A, SubFactor: Ac, Detail: "法律及會計服務業", Score: 4},
		{Factor: A, SubFactor: Ac, Detail: "民間融資性租賃業務", Score: 4},
		{Factor: A, SubFactor: Ac, Detail: "不動產相關行業", Score: 4},
		{Factor: A, SubFactor: Ac, Detail: "公務員（含聘僱人員）", Score: 4},
		{Factor: A, SubFactor: Ac, Detail: "宗教團體相關工作人員", Score: 4},
		{Factor: A, SubFactor: Ac, Detail: "國際貿易業", Score: 4},
		{Factor: A, SubFactor: Ad, Detail: "律師", Score: 4},
		{Factor: A, SubFactor: Ad, Detail: "公證人", Score: 4},
		{Factor: A, SubFactor: Ad, Detail: "會計師", Score: 4},
		{Factor: A, SubFactor: Ad, Detail: "地政士", Score: 4},
		{Factor: A, SubFactor: Ad, Detail: "記帳士", Score: 4},
		{Factor: A, SubFactor: Ad, Detail: "報稅代理人", Score: 4},
		{Factor: A, SubFactor: Ad, Detail: "中介（公司秘書）", Score: 4},
		{Factor: A, SubFactor: Ae, Detail: "足以懷疑使用匿名、假名、人頭、虛設行號或虛設法人、團體建立業務關係", Score: 12},
		{Factor: A, SubFactor: Ae, Detail: "經查證後，持用偽、變造身分證明文件", Score: 12},
		{Factor: A, SubFactor: Ae, Detail: "不願提供佐證資料用以確認用戶身分", Score: 12},
		{Factor: A, SubFactor: Ae, Detail: "審核失敗後重新實名驗證、提供資料不一", Score: 4},
		{Factor: A, SubFactor: Af, Detail: "比對身分證件後，與已婉拒黑名單用戶有親戚關係", Score: 12},
		{Factor: A, SubFactor: Af, Detail: "與已婉拒黑名單用戶相同居住地址或註冊地址", Score: 12},
		{Factor: A, SubFactor: Af, Detail: "基本資料填寫錯誤達3次(含)以上", Score: 8},
		{Factor: A, SubFactor: Af, Detail: "與其他用戶自拍照、背景、字跡相似或相同但無合理原因者", Score: 8},
		{Factor: A, SubFactor: Af, Detail: "60歲以上", Score: 6},
		{Factor: A, SubFactor: Af, Detail: "不同使用者間，使用者名稱/E-mail雷同，但無合理原因者", Score: 6},
		{Factor: A, SubFactor: Af, Detail: "於室外或公共場所進行自拍或證件照拍攝", Score: 6},
		{Factor: A, SubFactor: Af, Detail: "KYC文件拍攝裝置與連線裝置不合，且無合理原因者", Score: 2},
		{Factor: A, SubFactor: Af, Detail: "25歲以下，年收入100萬以上", Score: 2},
		{Factor: A, SubFactor: Af, Detail: "惡意註冊亂填寫資料和上傳照片或是上傳過與KYC無關的相片", Score: 2},
		{Factor: A, SubFactor: Af, Detail: "重複上傳無效或相同的文件達兩次或是補件三次不成功", Score: 2},
		{Factor: A, SubFactor: Af, Detail: "非常見email網域(常見email網域：gmail,yahoo,hotmail)", Score: 2},
		{Factor: A, SubFactor: Ag, Detail: "核資錯誤-姓名、身分證字號、生日", Score: 12},
		{Factor: A, SubFactor: Ag, Detail: "核資錯誤-戶籍地址、綁定銀行、使用者名稱、行業別", Score: 8},
		{Factor: A, SubFactor: Ag, Detail: "照會時聽到有他人提示聲音、非合理理由要求暫緩通話或移轉地點", Score: 8},
		{Factor: A, SubFactor: Ag, Detail: "核資流暢但照本宣科，異常多提供個人資訊（例如：地址鄰里、連續提供非需求資料）", Score: 2},
		{Factor: A, SubFactor: Ag, Detail: "詢問資金來源時不願意提供或表達不願意透漏者", Score: 2},
		{Factor: A, SubFactor: Ag, Detail: "回答內容與事實不符、回答內容前後矛盾或找理由搪塞解釋", Score: 2},
		{Factor: A, SubFactor: Ag, Detail: "照會時聽到翻找資料聲音", Score: 2},
		{Factor: A, SubFactor: Ag, Detail: "不願提供職業資訊、含糊其辭或對職業狀況不清楚者", Score: 2},
		{Factor: A, SubFactor: Ag, Detail: "無投資經驗，且對投資標的不清楚", Score: 2},
		{Factor: A, SubFactor: Ag, Detail: "跟隨朋友/家人投資，且對投資標的不清楚", Score: 2},
		{Factor: A, SubFactor: Ag, Detail: "用戶資金來源來自於朋友，且無法進一步說明朋友的資金來源", Score: 2},
		{Factor: A, SubFactor: Ag, Detail: "表示需要準備照會內容或是找理由晚點照會", Score: 2},
		{Factor: A, SubFactor: Ag, Detail: "說話結巴、回答不出來、不願意回答、刻意跳過問題或回答吞吐、突然咬字模糊、變小聲、惱羞成怒等", Score: 2},
		{Factor: A, SubFactor: Ag, Detail: "表示是第一次接觸虛擬通貨", Score: 2},
		{Factor: A, SubFactor: Ag, Detail: "三次以上致電無接聽", Score: 2},

		{Factor: B, SubFactor: Ba, Detail: "美國", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "中國", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "日本", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "韓國", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "奧地利", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "比利時", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "保加利亞", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "克羅埃西亞", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "賽普勒斯", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "捷克", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "丹麥", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "愛沙尼亞", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "芬蘭", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "法國", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "德國", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "希臘", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "匈牙利", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "愛爾蘭", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "義大利", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "拉脫維亞", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "立陶宛", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "盧森堡", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "馬爾他", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "荷蘭", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "波蘭", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "葡萄牙", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "羅馬尼亞", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "斯洛伐克", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "斯洛維尼亞", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "西班牙", Score: 12},
		{Factor: B, SubFactor: Ba, Detail: "瑞典", Score: 12},
		{Factor: B, SubFactor: Bb, Detail: "北韓", Score: 12},
		{Factor: B, SubFactor: Bb, Detail: "伊朗", Score: 12},
		{Factor: B, SubFactor: Bb, Detail: "緬甸", Score: 12},
		{Factor: B, SubFactor: Bb, Detail: "烏克蘭克里米亞地區", Score: 12},
		{Factor: B, SubFactor: Bb, Detail: "俄羅斯（即時）", Score: 12},
		{Factor: B, SubFactor: Bc, Detail: "阿爾巴尼亞", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "巴貝多", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "布吉納法索", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "開曼群島", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "剛果民主共和國", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "直布羅陀", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "海地", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "牙買加", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "約旦", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "馬利", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "莫三比克", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "奈及利亞", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "巴拿馬", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "菲律賓", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "塞內加爾", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "南非", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "南蘇丹", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "敘利亞", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "坦尚尼亞", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "土耳其", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "烏干達", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "阿拉伯聯合大公國", Score: 8},
		{Factor: B, SubFactor: Bc, Detail: "葉門", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "索馬利亞", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "蘇丹", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "利比亞", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "中非共和國", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "阿富汗", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "伊拉克", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "巴爾幹半島", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "白俄羅斯", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "古巴", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "埃塞爾比亞", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "尼加拉瓜", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "津巴布韋", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "委內瑞拉", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "美屬薩摩亞", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "安奎拉", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "巴哈馬", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "英屬維京群島", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "哥斯達黎加", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "斐濟", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "關島", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "馬紹爾群島", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "帛琉", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "巴拿馬", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "薩摩亞", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "千里達及托巴哥", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "土克斯及開科斯群島", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "美屬維京群島", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "萬那杜", Score: 8},
		{Factor: B, SubFactor: Bd, Detail: "幾內亞", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "幾內亞比索", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "烏克蘭", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "香港", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "澳門", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "越南", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "印尼", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "馬來西亞", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "薩摩亞", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "貝里斯", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "塞席爾", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "杜拜", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "加拿大", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "澳洲", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "辛巴威", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "柬埔寨", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "突尼斯", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "泰國", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "新加坡", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "瑞士", Score: 4},
		{Factor: B, SubFactor: Bd, Detail: "摩洛哥", Score: 4},
		{Factor: B, SubFactor: Be, Detail: "客戶具有雙重國籍", Score: 2},
		{Factor: B, SubFactor: Be, Detail: "客戶僅具有外國籍", Score: 2},

		{Factor: C, SubFactor: Ca, Detail: "曾判斷為協助第三方進行虛擬通貨買賣", Score: 6},
		{Factor: C, SubFactor: Ca, Detail: "不曾匯入虛擬通貨，僅轉入法幣後兌換為虛擬通貨轉出", Score: 2},
		{Factor: C, SubFactor: Cb, Detail: "台幣、虛擬通貨累積金額超過USD$1,000,000元", Score: 2},
		{Factor: C, SubFactor: Cc, Detail: "曾被銀行通報或圈存者", Score: 12},
		{Factor: C, SubFactor: Cc, Detail: "與黑名單用戶使用相同IP位址或GPS位置進行登入使用", Score: 12},
		{Factor: C, SubFactor: Cc, Detail: "與其他用戶使用相同IP，無家庭關係或其他合理緣由", Score: 8},
		{Factor: C, SubFactor: Cc, Detail: "註冊或經通知補件者時於一小時內，與其他客戶使用相同IP位址", Score: 8},
		{Factor: C, SubFactor: Cc, Detail: "曾觸發內部可疑樣態警示", Score: 2},
		{Factor: C, SubFactor: Cc, Detail: "用戶經申報1次以上", Score: 12},
		{Factor: C, SubFactor: Cc, Detail: "用戶經申報1次", Score: 8},
		{Factor: C, SubFactor: Cc, Detail: "曾被聯防165平台要求調閱資料者", Score: 8},
		{Factor: C, SubFactor: Cd, Detail: "轉出虛擬通貨至黑名單錢包地址", Score: 12},
		{Factor: C, SubFactor: Cd, Detail: "接受虛擬通貨錢包地址曾於轉出時或嗣後被認定涉及不法行為", Score: 8},
		{Factor: C, SubFactor: Cd, Detail: "使用銀行帳戶為法幣轉入與轉出媒介", Score: 2},
	}
	return db.Create(&records).Error
}

// Fake for ISeeder
func (s *SeederRisks) Fake(db *gorm.DB) error {
	return seed.ErrNotImplement
}
