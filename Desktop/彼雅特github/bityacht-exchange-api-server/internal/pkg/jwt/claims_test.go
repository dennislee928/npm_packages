package jwt

import (
	"testing"

	"bityacht-exchange-api-server/configs"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"

	"github.com/spf13/viper"
)

func TestJWT(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()

	managerPayloads := []ManagerPayload{
		{ManagersRolesID: 1, ID: 1, Name: "Admin"},
		{ManagersRolesID: 2, ID: 2, Name: "Manager"},
		{ManagersRolesID: 3, ID: 3, Name: "Reviewer"},
		{ManagersRolesID: 4, ID: 4, Name: "Operator"},
	}

	for _, payload := range managerPayloads {
		if accessToken, _, err := IssueManagerToken(payload); err != nil {
			t.Error(err)
		} else if claims, err := ValidateManager(accessToken); err != nil {
			t.Error(err)
		} else if claims.ManagerPayload != payload {
			t.Errorf("claims.ManagerPayload: '%+v' != '%+v'\n", claims.ManagerPayload, payload)
		}
	}

	userPayloads := []UserPayload{
		{ID: 1, Account: "account1", CountriesCode: "TWN", Type: users.TypeNaturalPerson, FirstName: "大明", LastName: "王", Level: 0, Status: usersmodifylogs.SLStatusEnable},
		{ID: 2, Account: "account2", CountriesCode: "ATG", Type: users.TypeNaturalPerson, FirstName: "FirstName", LastName: "LastName", Level: 1, Status: usersmodifylogs.SLStatusEnable},
		{ID: 3, Account: "account3", CountriesCode: "TWN", Type: users.TypeJuridicalPerson, FirstName: "小明", LastName: "李", Level: 2, Status: usersmodifylogs.SLStatusEnable},
		{ID: 4, Account: "account4", CountriesCode: "UGA", Type: users.TypeJuridicalPerson, FirstName: "FirstName1", LastName: "LastName1", Level: 3, Status: usersmodifylogs.SLStatusForzen},
	}

	for _, payload := range userPayloads {
		if accessToken, _, err := IssueUserToken(payload); err != nil {
			t.Error(err)
		} else if claims, err := ValidateUser(accessToken); err != nil {
			t.Error(err)
		} else if claims.UserPayload != payload {
			t.Errorf("claims.UserPayload: '%+v' != '%+v'\n", claims.UserPayload, payload)
		}
	}
}
