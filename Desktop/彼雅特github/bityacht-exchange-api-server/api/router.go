package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"bityacht-exchange-api-server/api/controllers/v1/admin"
	"bityacht-exchange-api-server/api/controllers/v1/admin/kycs"
	"bityacht-exchange-api-server/api/controllers/v1/admin/kycs/risks"
	"bityacht-exchange-api-server/api/controllers/v1/admin/receipts"
	"bityacht-exchange-api-server/api/controllers/v1/admin/settings/banners"
	"bityacht-exchange-api-server/api/controllers/v1/admin/settings/mainnets"
	"bityacht-exchange-api-server/api/controllers/v1/admin/settings/managers"
	"bityacht-exchange-api-server/api/controllers/v1/admin/settings/transactionpairs"
	"bityacht-exchange-api-server/api/controllers/v1/admin/spots"
	"bityacht-exchange-api-server/api/controllers/v1/admin/suspicioustxs"
	"bityacht-exchange-api-server/api/controllers/v1/admin/transactions"
	"bityacht-exchange-api-server/api/controllers/v1/admin/users"
	"bityacht-exchange-api-server/api/controllers/v1/admin/users/bank"
	"bityacht-exchange-api-server/api/controllers/v1/admin/users/invite"
	userskycs "bityacht-exchange-api-server/api/controllers/v1/admin/users/kycs"
	usersrisks "bityacht-exchange-api-server/api/controllers/v1/admin/users/kycs/risks"
	loginlogs "bityacht-exchange-api-server/api/controllers/v1/admin/users/loginlogs"
	"bityacht-exchange-api-server/api/controllers/v1/admin/users/reviewslogs"
	"bityacht-exchange-api-server/api/controllers/v1/admin/users/statuslogs"
	"bityacht-exchange-api-server/api/controllers/v1/assets"
	"bityacht-exchange-api-server/api/controllers/v1/cybavo"
	"bityacht-exchange-api-server/api/controllers/v1/kryptogo"
	"bityacht-exchange-api-server/api/controllers/v1/user"
	userassets "bityacht-exchange-api-server/api/controllers/v1/user/assets"
	userspot "bityacht-exchange-api-server/api/controllers/v1/user/assets/spot"
	"bityacht-exchange-api-server/api/controllers/v1/user/banks"
	"bityacht-exchange-api-server/api/controllers/v1/user/commissions"
	"bityacht-exchange-api-server/api/controllers/v1/user/idv"
	usersettings "bityacht-exchange-api-server/api/controllers/v1/user/settings"
	"bityacht-exchange-api-server/api/controllers/v1/user/wallets"
	"bityacht-exchange-api-server/api/middleware"
	"bityacht-exchange-api-server/configs"
	docs "bityacht-exchange-api-server/docs"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	"bityacht-exchange-api-server/internal/pkg/rbac"
)

func NewRouter() *gin.Engine {
	switch configs.Config.Log.Level {
	case zerolog.WarnLevel, zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel, zerolog.Disabled:
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(gin.Recovery(), middleware.CORS(), middleware.RequestID(), middleware.Logger())

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/health", middleware.SetSkipLog(), healthHandler)

			assetsGroup := v1.Group("/assets", middleware.SetSkipLog())
			{
				assetsGroup.GET("/banners/:Filename", assets.GetBannerHandler)
			}

			adminGroup := v1.Group("/admin")
			{
				roleCheck := middleware.ManagerRoleCheck
				adminGroup.GET("/auth", admin.AuthHandler)
				adminGroup.POST("/login", middleware.SetSensitive(), admin.LoginHandler)
				adminGroup.POST("/forgot-password", admin.ForgotPasswordHandler)

				adminGroup.Use(middleware.Authorization(jwt.TypeManager))
				adminGroup.PATCH("/password", middleware.SetSensitive(), admin.UpdatePasswordHandler)
				adminGroup.POST("/logout", admin.LogoutHandler)

				settingsGroup := adminGroup.Group("/settings")
				{
					managersGroup := settingsGroup.Group("/managers", roleCheck(rbac.ObjectAdministratorAccountSettings, rbac.ActionRead))
					{
						managersGroup.POST("", managers.CreateHandler)
						managersGroup.GET("", managers.GetHandler)
						managersGroup.PATCH("/:ID", middleware.SetSensitive(), managers.UpdateHandler)
						managersGroup.DELETE("/:ID", managers.DeleteHandler)
					}

					transactionPairsGroup := settingsGroup.Group("/transactionpairs", roleCheck(rbac.ObjectTradingpairhandlingfeesetting, rbac.ActionRead))
					{
						transactionPairsGroup.GET("", transactionpairs.GetHandler)
						transactionPairsGroup.PATCH("", transactionpairs.UpdateHandler)
					}

					bannersGroup := settingsGroup.Group("/banners", roleCheck(rbac.ObjectBannermanagement, rbac.ActionRead))
					{
						bannersGroup.POST("", middleware.SetSensitive(), banners.CreateHandler)
						bannersGroup.GET("", banners.GetListHandler)
						bannersGroup.PATCH("/:ID", middleware.SetSensitive(), banners.UpdateHandler)
						bannersGroup.DELETE("/:ID", banners.DeleteHandler)
						bannersGroup.POST("/priority", banners.PriorityHandler)
					}

					mainnetsGroup := settingsGroup.Group("/mainnets")
					{
						mainnetsGroup.GET("", mainnets.GetHandler)
						mainnetsGroup.PATCH("/:Currency/:Mainnet", mainnets.UpdateHandler)
					}
				}

				usersGroup := adminGroup.Group("/users", roleCheck(rbac.ObjectMemberList, rbac.ActionRead))
				{
					usersGroup.GET("", users.GetListHandler)
					usersGroup.POST("", users.CreateHandler)
					usersGroup.GET("/export", roleCheck(rbac.ObjectMemberList, rbac.ActionExport), users.ExportHandler)
					usersGroup.GET("/options", middleware.SetSkipLog(), users.GetOptionsHandler)

					usersGroup.GET("/:ID", users.GetHandler)
					usersGroup.GET("/:ID/withdrawal-whitelist", users.GetWithdrawalWhitelistHandler)
					usersGroup.GET("/:ID/withdrawal-whitelist/export", roleCheck(rbac.ObjectMemberList, rbac.ActionExport), users.ExportWithdrawalWhitelistHandler)
					usersGroup.DELETE("/:ID/withdrawal-whitelist/:WhitelistID", users.DeleteWithdrawalWhitelistHandler)
					usersGroup.PATCH("/:ID/level", roleCheck(rbac.ObjectMemberList, rbac.ActionLevelChange), users.UpdateLevelHandler)
					usersGroup.PATCH("/:ID/status", roleCheck(rbac.ObjectMemberList, rbac.ActionStatusChange), users.UpdateStatusHandler)

					usersGroup.GET("/:ID/invite-info", invite.GetInfoHandler)
					usersGroup.GET("/:ID/invite", invite.GetHandler)
					usersGroup.GET("/:ID/invite-rewards", invite.GetRewardsHandler)
					usersGroup.GET("/:ID/invite-rewards/export", invite.ExportRewardsHandler)

					loginlogGroup := usersGroup.Group("/:ID/login-logs")
					{
						loginlogGroup.GET("", loginlogs.GetListHandler)
						loginlogGroup.GET("/export", roleCheck(rbac.ObjectMemberList, rbac.ActionExport), loginlogs.GetExportHandler)
					}

					statuslogGroup := usersGroup.Group("/:ID/statuslogs")
					{
						statuslogGroup.GET("", statuslogs.GetHandler)
						statuslogGroup.GET("/export", roleCheck(rbac.ObjectMemberList, rbac.ActionStatusLogExport), statuslogs.ExportHandler)
					}

					userKycsGroup := usersGroup.Group("/:ID/kycs", roleCheck(rbac.ObjectKYCVerificationReview, rbac.ActionRead))
					{
						userKycsGroup.GET("", userskycs.GetWithDDHandler)

						// TODO: make sure the action of roleCheck
						// userKycsGroup.PUT("/result-image", roleCheck(rbac.ObjectKYCVerificationReview, rbac.ActionKryptoReview), middleware.SetSensitive(), userskycs.UploadResultImageHandler) //! Deprecated (Meeting at 2023/11/01)
						userKycsGroup.PATCH("/idv-audit-status", roleCheck(rbac.ObjectKYCVerificationReview, rbac.ActionKryptoReview), middleware.SetSensitive(), userskycs.UpdateIDVAuditStatusHandler)
						userKycsGroup.PATCH("/task-id", roleCheck(rbac.ObjectKYCVerificationReview, rbac.ActionKryptoReview), userskycs.UpdateKryptoGOTaskIDHandler)

						userKycsGroup.POST("/resent-krypto", roleCheck(rbac.ObjectKYCVerificationReview, rbac.ActionResend), userskycs.ResentToKryptoGoHandler)
						userKycsGroup.PATCH("/krypto-review", roleCheck(rbac.ObjectKYCVerificationReview, rbac.ActionKryptoReview), userskycs.UpdateKryptoReviewHandler)
						userKycsGroup.PATCH("/name-check", roleCheck(rbac.ObjectKYCVerificationReview, rbac.ActionNameCheckUpload), userskycs.UpdateNameCheckHandler)
						userKycsGroup.GET("/name-check/pdf", userskycs.GetNameCheckPdfHandler)
						userKycsGroup.PATCH("/compliance-review", roleCheck(rbac.ObjectKYCVerificationReview, rbac.ActionComplianceReview), userskycs.UpdateComplianceReviewHandler)
						userKycsGroup.PATCH("/final-review", roleCheck(rbac.ObjectKYCVerificationReview, rbac.ActionFinalReview), userskycs.UpdateFinalReviewHandler)

						reviewslogGroup := userKycsGroup.Group("/reviewslogs")
						{
							reviewslogGroup.GET("", reviewslogs.GetHandler)
							reviewslogGroup.GET("/export", roleCheck(rbac.ObjectKYCVerificationReview, rbac.ActionReviewLogExport), reviewslogs.ExportHandler)
						}

						risksGroup := userKycsGroup.Group("/risks")
						{
							risksGroup.GET("", usersrisks.GetHandler)
							risksGroup.POST("", usersrisks.UpdateHandler)
						}
					}

					// TODO: make sure the action of roleCheck
					userBankGroup := usersGroup.Group("/:ID/bank", roleCheck(rbac.ObjectWithdrawalDepositFlatlist, rbac.ActionRead))
					{
						userBankGroup.GET("", bank.GetHandler)
						userBankGroup.PATCH("", bank.PatchHandler)

						bankLogGroup := userBankGroup.Group("/logs")
						{
							bankLogGroup.GET("", bank.GetLogHandler)
							bankLogGroup.GET("/export", bank.ExportLogHandler)
						}
					}
				}

				kycsGroup := adminGroup.Group("/kycs", roleCheck(rbac.ObjectKYCVerificationReview, rbac.ActionRead))
				{
					kycsGroup.GET("", kycs.GetListHandler)
					kycsGroup.GET("/export", roleCheck(rbac.ObjectKYCVerificationReview, rbac.ActionExport), kycs.ExportHandler)

					annualKycsGroup := kycsGroup.Group("/annual")
					{
						annualKycsGroup.GET("", kycs.GetAnnualListHandler)
						annualKycsGroup.GET("/export", roleCheck(rbac.ObjectKYCVerificationReview, rbac.ActionExport), kycs.ExportAnnualHandler)
					}

					risksGroup := kycsGroup.Group("/risks")
					{
						risksGroup.GET("", risks.GetListHandler)
						risksGroup.POST("", risks.CreateHandler)
						risksGroup.PATCH("/:ID", risks.UpdateHandler)
						risksGroup.DELETE("/:ID", risks.DeleteHandler)
					}
				}

				spotsGroup := adminGroup.Group("/spots", roleCheck(rbac.ObjectWithdrawalDepositCoinlist, rbac.ActionRead))
				{
					spotsGroup.GET("", spots.GetListHandler)
					spotsGroup.GET("/export", roleCheck(rbac.ObjectWithdrawalDepositCoinlist, rbac.ActionExport), spots.ExportHandler)
					spotsGroup.GET("/aegis-export", roleCheck(rbac.ObjectWithdrawalDepositCoinlist, rbac.ActionExport), spots.AegisExportHandler)
					spotsGroup.POST("/aegis-import", roleCheck(rbac.ObjectWithdrawalDepositCoinlist, rbac.ActionExport), middleware.SetSensitive(), spots.AegisImportHandler) // TODO: make sure the action of roleCheck
				}

				transactionsGroup := adminGroup.Group("/transactions", roleCheck(rbac.ObjectTransactions, rbac.ActionRead))
				{
					transactionsGroup.GET("", transactions.GetListHandler)
					transactionsGroup.GET("/export", roleCheck(rbac.ObjectWithdrawalDepositCoinlist, rbac.ActionExport), transactions.ExportHandler)
				}

				receiptsGroup := adminGroup.Group("/receipts", roleCheck(rbac.ObjectInvoicemanagement, rbac.ActionRead))
				{
					receiptsGroup.GET("", receipts.GetListHandler)
					receiptsGroup.GET("/:ID", receipts.GetDetailHandler)
					receiptsGroup.POST("/issue", receipts.IssueHandler)
					receiptsGroup.GET("/export", roleCheck(rbac.ObjectInvoicemanagement, rbac.ActionExport), receipts.ExportCSVHandler)
				}

				suspiciousTXs := adminGroup.Group("/suspicious-txs", roleCheck(rbac.ObjectSuspicioustransactionrecords, rbac.ActionRead))
				{
					suspiciousTXs.GET("", suspicioustxs.GetListHandler)
					suspiciousTXs.GET("/export-csv", suspicioustxs.ExportCSVHandler)

					suspiciousTXs.GET("/:ID", suspicioustxs.GetDetailHandler)
					suspiciousTXs.PATCH("/:ID", suspicioustxs.PatchHandler)
					suspiciousTXs.POST("/:ID/file", middleware.SetSensitive(), suspicioustxs.UploadFileHandler)
					suspiciousTXs.GET("/:ID/file", suspicioustxs.DownloadFileHandler)
					suspiciousTXs.DELETE("/:ID/file", suspicioustxs.DeleteFileHandler)
				}
			}

			userGroup := v1.Group("/user")
			{
				permCheck := middleware.UserPermissionCheck
				userGroup.POST("/register", middleware.SetSensitive(), user.RegisterHandler)
				userGroup.POST("/resend-verify", user.ResendEmailVerificationCodeHandler)
				userGroup.POST("/verify", user.VerifyEmailHandler)

				userGroup.POST("/login", middleware.SetSensitive(), user.LoginHandler)
				userGroup.POST("/2fa-login", user.TwoFactorLoginHandler)
				userGroup.POST("/forgot-password", user.ForgotPasswordHandler)
				userGroup.POST("/verify-reset-password", user.VerifyResetPasswordHandler)
				userGroup.POST("/reset-password", middleware.SetSensitive(), user.ResetPasswordHandler)
				userGroup.POST("/token", user.RefreshTokenHandler)

				userGroup.GET("/spot-trend", middleware.SetSkipLog(), user.GetSpotTrendHandler)
				userGroup.GET("/banners", middleware.SetSkipLog(), user.GetBannersHandler)

				userGroup.Use(middleware.Authorization(jwt.TypeUser))

				userGroup.POST("/logout", user.LogoutHandler)

				userGroup.GET("/info", user.GetUserInfoHandler)
				userGroup.GET("/login-logs", user.GetUserLoginLogsHandler)

				commissionsGroup := userGroup.Group("/commissions")
				{
					commissionsGroup.GET("", commissions.GetHandler)
					commissionsGroup.POST("/withdraw", permCheck(0, 1, true), commissions.WithdrawHandler)
				}

				settingsGroup := userGroup.Group("/settings")
				{
					settingsGroup.PATCH("/password", middleware.SetSensitive(), usersettings.UpdatePasswordHandler)
					// settingsGroup.PATCH("/2fa", usersettings.UpdateLogin2FATypeHandler) //! Deprecated (Meeting at 2023/10/2)
					settingsGroup.POST("/issue-withdraw-2fa-verify", usersettings.IssueWithdraw2FAVerifyHandler)
					settingsGroup.GET("/withdraw-2fa-info", usersettings.GetWithdraw2FAInfoHandler)
					settingsGroup.PATCH("/withdraw-2fa", usersettings.UpdateWithdraw2FAHandler)

					settingsGroup.PATCH("/mobile-barcode", usersettings.UpdateMobileBarcodeHandler)
				}

				banksGroup := userGroup.Group("/banks")
				{
					banksGroup.GET("/options", middleware.SetSkipLog(), banks.GetOptionsHandler)

					banksGroup.PUT("/account", middleware.SetSensitive(), banks.UpsertAccountHandler)
					banksGroup.DELETE("/account", banks.DeleteAccountHandler)
				}

				idvGroup := userGroup.Group("/idv")
				{
					idvGroup.GET("/options", middleware.SetSkipLog(), idv.GetOptionsHandler)

					idvGroup.Use(permCheck(0, -1, true))
					idvGroup.POST("", middleware.SetSensitive(), idv.CreateIDVerificationHandler)
					idvGroup.POST("/check-phone", idv.CheckPhoneHandler)
					idvGroup.POST("/issue-phone-verify", middleware.SetSensitive(), idv.IssuePhoneVerificationCodeHandler)
					idvGroup.POST("/verify-phone", idv.VerifyPhoneHandler)
					idvGroup.GET("/krypto-go-url", idv.GetKryptoGoURLHandler)
					idvGroup.PATCH("/image", idv.UpdateIDVImageHandler)
				}

				assetsGroup := userGroup.Group("/assets")
				{
					assetsGroup.GET("", userassets.GetHandler)

					spotGroup := assetsGroup.Group("/spot")
					{
						spotGroup.GET("/histories", userspot.GetHistoriesHandler)
						spotGroup.GET("/transactions", userspot.GetTransactionsHandler)
						spotGroup.GET("/options", userspot.GetOptionsHandler)
						spotGroup.GET("/trend", middleware.SetSkipLog(), userspot.GetTrendHandler)

						spotGroup.POST("/transactions", permCheck(1, 1, true), userspot.CreateTransactionHandler)
					}
				}

				walletsGroup := userGroup.Group("/wallets")
				{
					walletsGroup.GET("/deposit/address", permCheck(1, 1, true), wallets.GetDepositAddressHandler)
					walletsGroup.POST("/deposit/address", permCheck(1, 1, true), wallets.GenDepositAddressHandler)
					walletsGroup.POST("/withdraw", permCheck(2, 1, true), wallets.WithdrawHandler)
					walletsGroup.POST("/2fa-withdraw", permCheck(2, 1, true), wallets.TwoFactorWithdrawHandler)
					walletsGroup.GET("/withdraw-info", wallets.GetWithdrawInfoHandler)

					withdrawalWhitelistGroup := walletsGroup.Group("/withdrawal-whitelist")
					{
						withdrawalWhitelistGroup.GET("", wallets.GetWithdrawalWhitelistHandler)
						withdrawalWhitelistGroup.POST("", wallets.CreateWithdrawalWhitelistHandler)
						withdrawalWhitelistGroup.DELETE("/:ID", wallets.DeleteWithdrawalWhitelistHandler)
					}
				}
			}

			kryptoGoGroup := v1.Group("/krypto-go")
			{
				// TODO: Check the request is from KryptoGO
				kryptoGoGroup.POST("/idv-callback/:UsersID/:IDVsID", kryptogo.IDVCallbackHandler)
				kryptoGoGroup.POST("/dd-callback/:UsersID/:DDsID", kryptogo.DDCallbackHandler)
			}

			cybavoGroup := v1.Group("/cybavo")
			{
				cybavoGroup.POST("/callback", cybavo.CallbackHandler)
				cybavoGroup.POST("/withdrawal-callback", cybavo.WithdrawalCallbackHandler)
			}
		}
	}

	return r
}

func healthHandler(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
