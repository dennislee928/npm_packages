package rbac

import (
	"strconv"
	"testing"

	"github.com/spf13/viper"

	"bityacht-exchange-api-server/configs"
)

func TestRBAC(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	Init()

	rolesID := int64(1)
	// Check the permission.
	ans := casbinEnforcer.GetPermissionsForUser(strconv.FormatInt(rolesID, 10))
	t.Log(ans)

	// Check the permission.
	b, ok := Enforce(rolesID, ObjectMemberList, ActionRead)
	t.Log(b, ok)

	// Check the permission.
	b, ok = Enforce(rolesID, "data1", "read")
	t.Log(b, ok)
}
