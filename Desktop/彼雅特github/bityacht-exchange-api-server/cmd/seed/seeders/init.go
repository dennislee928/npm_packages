package seeders

import (
	"bityacht-exchange-api-server/cmd/seed"
)

var Offset int

func init() {
	var seederOrder []string

	seederOrder = append(seederOrder, "ManagersRoles", "Managers")
	seederOrder = append(seederOrder, "Currencies", "Mainnets", "TransactionPairs")
	seederOrder = append(seederOrder, "Users", "UsersWallets", "DueDiligences")
	seederOrder = append(seederOrder, "Banks", "BankBranchs")

	seed.SetSeedersOrder(seederOrder)
}
