package kyc

var (
	KryptoGO *kryptoGO
)

func Init() {
	KryptoGO = newKryptoGO()
}
