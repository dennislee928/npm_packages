package seeders

import (
	"bityacht-exchange-api-server/cmd/migrate/migrations"
	"bityacht-exchange-api-server/cmd/seed"
	"bityacht-exchange-api-server/internal/pkg/rbac"
	"strconv"

	"gorm.io/gorm"
)

func init() {
	seed.Register(&SeederManagersRoles{})
}

// SeederManagersRoles is a kind of ISeeder, You can set it as same as model.
type SeederManagersRoles migrations.Migration1689660728

// SeederName for ISeeder
func (*SeederManagersRoles) SeederName() string {
	return "ManagersRoles"
}

// TableName for gorm
func (*SeederManagersRoles) TableName() string {
	return (*migrations.Migration1689660728)(nil).TableName()
}

// Default for ISeeder
func (*SeederManagersRoles) Default(db *gorm.DB) error {
	type policy struct {
		Object string
		Action string
	}

	managersRolesAndPolicies := []struct {
		Role     SeederManagersRoles
		Policies []policy
	}{
		{
			Role: SeederManagersRoles{
				Name: "Administrator",
			},
			Policies: []policy{
				// for Member List
				{Object: rbac.ObjectMemberList, Action: rbac.ActionRead},
				{Object: rbac.ObjectMemberList, Action: rbac.ActionCreateJuridicalMember},
				{Object: rbac.ObjectMemberList, Action: rbac.ActionExport},
				{Object: rbac.ObjectMemberList, Action: rbac.ActionStatusChange},
				{Object: rbac.ObjectMemberList, Action: rbac.ActionStatusLogExport},
				{Object: rbac.ObjectMemberList, Action: rbac.ActionLevelChange},

				// for kyc verification review
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionRead},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionResend},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionKryptoReview},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionNameCheckUpload},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionInternalRiskReview},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionComplianceReview},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionFinalReview},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionExport},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionReviewLogExport},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionForeignKryptoIDUpdate},

				{Object: rbac.ObjectTransactions, Action: rbac.ActionRead},
				{Object: rbac.ObjectTransactions, Action: rbac.ActionExport},

				{Object: rbac.ObjectWithdrawalDepositCoinlist, Action: rbac.ActionRead},
				{Object: rbac.ObjectWithdrawalDepositCoinlist, Action: rbac.ActionExport},

				{Object: rbac.ObjectWithdrawalDepositFlatlist, Action: rbac.ActionRead},
				{Object: rbac.ObjectWithdrawalDepositFlatlist, Action: rbac.ActionExport},

				{Object: rbac.ObjectSuspicioustransactionrecords, Action: rbac.ActionRead},

				{Object: rbac.ObjectInvoicemanagement, Action: rbac.ActionRead},
				{Object: rbac.ObjectInvoicemanagement, Action: rbac.ActionExport},
				{Object: rbac.ObjectInvoicemanagement, Action: rbac.ActionInvoiceUpdate},

				{Object: rbac.ObjectTradingpairhandlingfeesetting, Action: rbac.ActionRead},
				{Object: rbac.ObjectAdministratorAccountSettings, Action: rbac.ActionRead},
				{Object: rbac.ObjectBannermanagement, Action: rbac.ActionRead},
			},
		},
		{
			Role: SeederManagersRoles{
				Name: "Compliance",
			},
			Policies: []policy{
				// for Member List
				{Object: rbac.ObjectMemberList, Action: rbac.ActionRead},
				{Object: rbac.ObjectMemberList, Action: rbac.ActionCreateJuridicalMember},
				{Object: rbac.ObjectMemberList, Action: rbac.ActionExport},
				{Object: rbac.ObjectMemberList, Action: rbac.ActionStatusChange},
				{Object: rbac.ObjectMemberList, Action: rbac.ActionStatusLogExport},
				{Object: rbac.ObjectMemberList, Action: rbac.ActionLevelChange},

				// for kyc verification review
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionRead},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionResend},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionKryptoReview},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionNameCheckUpload},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionInternalRiskReview},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionComplianceReview},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionFinalReview},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionExport},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionReviewLogExport},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionForeignKryptoIDUpdate},

				{Object: rbac.ObjectTransactions, Action: rbac.ActionRead},
				{Object: rbac.ObjectTransactions, Action: rbac.ActionExport},

				{Object: rbac.ObjectWithdrawalDepositCoinlist, Action: rbac.ActionRead},
				{Object: rbac.ObjectWithdrawalDepositCoinlist, Action: rbac.ActionExport},

				{Object: rbac.ObjectWithdrawalDepositFlatlist, Action: rbac.ActionRead},
				{Object: rbac.ObjectWithdrawalDepositFlatlist, Action: rbac.ActionExport},

				{Object: rbac.ObjectSuspicioustransactionrecords, Action: rbac.ActionRead},

				{Object: rbac.ObjectInvoicemanagement, Action: rbac.ActionRead},
				{Object: rbac.ObjectInvoicemanagement, Action: rbac.ActionExport},
				{Object: rbac.ObjectInvoicemanagement, Action: rbac.ActionInvoiceUpdate},

				{Object: rbac.ObjectTradingpairhandlingfeesetting, Action: rbac.ActionRead},
				// {Object: rbac.ObjectAdministratorAccountSettings, Action: rbac.ActionRead},
				{Object: rbac.ObjectBannermanagement, Action: rbac.ActionRead},
			},
		},
		{
			Role: SeederManagersRoles{
				Name: "Customer Service",
			},
			Policies: []policy{
				// for Member List
				{Object: rbac.ObjectMemberList, Action: rbac.ActionRead},
				// {Object: rbac.ObjectMemberList, Action: rbac.ActionCreateJuridicalMember},
				// {Object: rbac.ObjectMemberList, Action: rbac.ActionExport},
				// {Object: rbac.ObjectMemberList, Action: rbac.ActionStatusChange},
				// {Object: rbac.ObjectMemberList, Action: rbac.ActionStatusLogExport},
				// {Object: rbac.ObjectMemberList, Action: rbac.ActionLevelChange},

				// for kyc verification review
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionRead},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionResend},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionKryptoReview},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionNameCheckUpload},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionInternalRiskReview},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionComplianceReview},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionFinalReview},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionExport},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionReviewLogExport},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionForeignKryptoIDUpdate},

				{Object: rbac.ObjectTransactions, Action: rbac.ActionRead},
				// {Object: rbac.ObjectTransactions, Action: rbac.ActionExport},

				{Object: rbac.ObjectWithdrawalDepositCoinlist, Action: rbac.ActionRead},
				// {Object: rbac.ObjectWithdrawalDepositCoinlist, Action: rbac.ActionExport},

				{Object: rbac.ObjectWithdrawalDepositFlatlist, Action: rbac.ActionRead},
				// {Object: rbac.ObjectWithdrawalDepositFlatlist, Action: rbac.ActionExport},

				// {Object: rbac.ObjectSuspicioustransactionrecords, Action: rbac.ActionRead},

				// {Object: rbac.ObjectInvoicemanagement, Action: rbac.ActionRead},
				// {Object: rbac.ObjectInvoicemanagement, Action: rbac.ActionExport},
				// {Object: rbac.ObjectInvoicemanagement, Action: rbac.ActionInvoiceUpdate},

				// {Object: rbac.ObjectTradingpairhandlingfeesetting, Action: rbac.ActionRead},
				// {Object: rbac.ObjectAdministratorAccountSettings, Action: rbac.ActionRead},
				{Object: rbac.ObjectBannermanagement, Action: rbac.ActionRead},
			},
		},
		{
			Role: SeederManagersRoles{
				Name: "Risk Management",
			},
			Policies: []policy{
				// for Member List
				{Object: rbac.ObjectMemberList, Action: rbac.ActionRead},
				{Object: rbac.ObjectMemberList, Action: rbac.ActionCreateJuridicalMember},
				{Object: rbac.ObjectMemberList, Action: rbac.ActionExport},
				{Object: rbac.ObjectMemberList, Action: rbac.ActionStatusChange},
				{Object: rbac.ObjectMemberList, Action: rbac.ActionStatusLogExport},
				{Object: rbac.ObjectMemberList, Action: rbac.ActionLevelChange},

				// for kyc verification review
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionRead},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionResend},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionKryptoReview},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionNameCheckUpload},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionInternalRiskReview},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionComplianceReview},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionFinalReview},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionExport},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionReviewLogExport},
				{Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionForeignKryptoIDUpdate},

				{Object: rbac.ObjectTransactions, Action: rbac.ActionRead},
				{Object: rbac.ObjectTransactions, Action: rbac.ActionExport},

				{Object: rbac.ObjectWithdrawalDepositCoinlist, Action: rbac.ActionRead},
				{Object: rbac.ObjectWithdrawalDepositCoinlist, Action: rbac.ActionExport},

				{Object: rbac.ObjectWithdrawalDepositFlatlist, Action: rbac.ActionRead},
				{Object: rbac.ObjectWithdrawalDepositFlatlist, Action: rbac.ActionExport},

				{Object: rbac.ObjectSuspicioustransactionrecords, Action: rbac.ActionRead},

				// {Object: rbac.ObjectInvoicemanagement, Action: rbac.ActionRead},
				// {Object: rbac.ObjectInvoicemanagement, Action: rbac.ActionExport},
				// {Object: rbac.ObjectInvoicemanagement, Action: rbac.ActionInvoiceUpdate},

				// {Object: rbac.ObjectTradingpairhandlingfeesetting, Action: rbac.ActionRead},
				// // {Object: rbac.ObjectAdministratorAccountSettings, Action: rbac.ActionRead},
				// {Object: rbac.ObjectBannermanagement, Action: rbac.ActionRead},
			},
		},
		{
			Role: SeederManagersRoles{
				Name: "Finance",
			},
			Policies: []policy{
				// for Member List
				{Object: rbac.ObjectMemberList, Action: rbac.ActionRead},
				// {Object: rbac.ObjectMemberList, Action: rbac.ActionCreateJuridicalMember},
				// {Object: rbac.ObjectMemberList, Action: rbac.ActionExport},
				// {Object: rbac.ObjectMemberList, Action: rbac.ActionStatusChange},
				// {Object: rbac.ObjectMemberList, Action: rbac.ActionStatusLogExport},
				// {Object: rbac.ObjectMemberList, Action: rbac.ActionLevelChange},

				// for kyc verification review
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionRead},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionResend},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionKryptoReview},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionNameCheckUpload},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionInternalRiskReview},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionComplianceReview},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionFinalReview},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionExport},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionReviewLogExport},
				// {Object: rbac.ObjectKYCVerificationReview, Action: rbac.ActionForeignKryptoIDUpdate},

				{Object: rbac.ObjectTransactions, Action: rbac.ActionRead},
				{Object: rbac.ObjectTransactions, Action: rbac.ActionExport},

				{Object: rbac.ObjectWithdrawalDepositCoinlist, Action: rbac.ActionRead},
				{Object: rbac.ObjectWithdrawalDepositCoinlist, Action: rbac.ActionExport},

				{Object: rbac.ObjectWithdrawalDepositFlatlist, Action: rbac.ActionRead},
				{Object: rbac.ObjectWithdrawalDepositFlatlist, Action: rbac.ActionExport},

				// {Object: rbac.ObjectSuspicioustransactionrecords, Action: rbac.ActionRead},

				{Object: rbac.ObjectInvoicemanagement, Action: rbac.ActionRead},
				{Object: rbac.ObjectInvoicemanagement, Action: rbac.ActionExport},
				{Object: rbac.ObjectInvoicemanagement, Action: rbac.ActionInvoiceUpdate},

				{Object: rbac.ObjectTradingpairhandlingfeesetting, Action: rbac.ActionRead},
				// {Object: rbac.ObjectAdministratorAccountSettings, Action: rbac.ActionRead},
				// {Object: rbac.ObjectBannermanagement, Action: rbac.ActionRead},
			},
		},
	}

	for _, managersRolesAndPolicies := range managersRolesAndPolicies {
		// #nosec G601
		if err := db.Create(&managersRolesAndPolicies.Role).Error; err != nil {
			return err
		}

		subject := strconv.FormatInt(managersRolesAndPolicies.Role.ID, 10)

		policies := make([]migrations.Migration1689902251, 0, len(managersRolesAndPolicies.Policies))
		for _, policy := range managersRolesAndPolicies.Policies {
			policies = append(policies, migrations.Migration1689902251{
				Ptype:   "p",
				Subject: subject,
				Object:  policy.Object,
				Action:  policy.Action,
			})
		}

		if err := db.Create(&policies).Error; err != nil {
			return err
		}
	}

	return nil
}

// Fake for ISeeder
func (*SeederManagersRoles) Fake(db *gorm.DB) error {
	return seed.ErrNotImplement
}
