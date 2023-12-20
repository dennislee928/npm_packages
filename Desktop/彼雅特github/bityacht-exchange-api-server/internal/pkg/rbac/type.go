package rbac

const (
	ObjectMemberList                    = "member_list"
	ObjectKYCVerificationReview         = "kyc_verification_review"
	ObjectTransactions                  = "transactions"
	ObjectWithdrawalDepositCoinlist     = "withdrawal_deposit_coin_list"
	ObjectWithdrawalDepositFlatlist     = "withdrawal_deposit_flat_list"
	ObjectSuspicioustransactionrecords  = "suspicious_transaction_records"
	ObjectInvoicemanagement             = "invoice_management"
	ObjectTradingpairhandlingfeesetting = "trading_pair_handling_fee_setting"
	ObjectAdministratorAccountSettings  = "administrator_account_settings"
	ObjectBannermanagement              = "banner_management"
)

const (
	// for all objects
	ActionRead   = "read"
	ActionExport = "export"

	// for member list
	ActionCreateJuridicalMember = "create_juridical_member"
	ActionStatusChange          = "status_change"
	ActionLevelChange           = "level_change"
	ActionStatusLogExport       = "status_log_export"

	// for kyc verification review
	ActionResend                = "resend"
	ActionKryptoReview          = "krypto_review"
	ActionNameCheckUpload       = "name_check_upload"
	ActionInternalRiskReview    = "internal_risk_review"
	ActionComplianceReview      = "compliance_review"
	ActionFinalReview           = "final_review"
	ActionReviewLogExport       = "review_log_export"
	ActionForeignKryptoIDUpdate = "foreign_krypto_id_update" // #nosec G101

	// for invoic
	ActionInvoiceUpdate = "invoice_update"
)
