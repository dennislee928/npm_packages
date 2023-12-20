package rbac

import (
	"bityacht-exchange-api-server/internal/database/sql"
	_ "embed"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

//go:embed rbac_model.conf
var modelConf string

var casbinEnforcer *casbin.Enforcer

func Init() {
	model, err := model.NewModelFromString(modelConf)
	if err != nil {
		panic(err)
	}

	adapter, err := gormadapter.NewAdapterByDBUseTableName(sql.DB(), "", "managers_roles_policies")
	if err != nil {
		panic(err)
	}

	if casbinEnforcer, err = casbin.NewEnforcer(model, adapter); err != nil {
		panic(err)
	}
	if err := casbinEnforcer.LoadPolicy(); err != nil {
		panic(err)
	}
}

func AddPolicy(sub, obj, act string) error {
	if _, err := casbinEnforcer.AddPolicy(sub, obj, act); err != nil {
		return err
	}
	if err := casbinEnforcer.SavePolicy(); err != nil {
		return err
	}
	return nil
}

func AddPoliciesEx(data [][]string) error {
	if _, err := casbinEnforcer.AddPoliciesEx(data); err != nil {
		return err
	}
	if err := casbinEnforcer.SavePolicy(); err != nil {
		return err
	}
	return nil
}

func RemoveBySub(sub string) error {
	if _, err := casbinEnforcer.RemovePolicy(sub); err != nil {
		return err
	}
	if err := casbinEnforcer.SavePolicy(); err != nil {
		return err
	}
	return nil
}

func RemoveAPolicy(sub, obj, act string) error {
	if _, err := casbinEnforcer.RemovePolicy(sub, obj, act); err != nil {
		return err
	}
	if err := casbinEnforcer.SavePolicy(); err != nil {
		return err
	}
	return nil
}

func GetAllPolicyBySub(sub string) [][]string {
	return casbinEnforcer.GetPermissionsForUser(sub)
}

func Enforce(sub int64, obj string, act string) (bool, error) {
	b, err := casbinEnforcer.Enforce(strconv.FormatInt(sub, 10), obj, act)
	if err != nil {
		return b, err
	}
	return b, nil
}
