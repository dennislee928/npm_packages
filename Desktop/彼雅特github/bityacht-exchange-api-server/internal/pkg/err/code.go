package errpkg

// Code for error
type Code int64

// Code List for HTTP Code 4xx Series
const (
	CodeBindBody                     Code = 4000
	CodeUnauthorized                 Code = 4001
	CodeJWTBadSigning                Code = 4002
	CodeJWTInvalid                   Code = 4003
	CodeJWTRevoked                   Code = 4004
	CodeBadAuthorizationToken        Code = 4005
	CodePermissionDenied             Code = 4006
	CodeAccountDuplicated            Code = 4007
	CodeBadParam                     Code = 4008
	CodeRecordNotFound               Code = 4009
	CodeRecordNoChange               Code = 4010
	CodeBadBody                      Code = 4011
	CodeAccountNotAvailable          Code = 4012
	CodeInsufficientFunds            Code = 4013
	CodeBadQuery                     Code = 4014
	CodeFundDeleted                  Code = 4015
	CodeTooManyRequests              Code = 4016
	CodeBadPasswordStrength          Code = 4017
	CodeOverResetPasswordLimit       Code = 4018
	CodeBadUploadedFileType          Code = 4019
	CodeFileNotFound                 Code = 4020
	CodeEmailAlreadyVerified         Code = 4021
	CodeBadVerificationCode          Code = 4022
	CodeVerificationCodeExpired      Code = 4023
	CodeBadMobileBarcodeFormat       Code = 4024
	CodeBadInviteCode                Code = 4025
	CodeImageOverSize                Code = 4026
	CodeBadDataURIImageFormat        Code = 4027
	CodeBadImageData                 Code = 4028
	CodeBadNationalID                Code = 4029
	CodeNationalIDDuplicated         Code = 4030
	CodeBadAction                    Code = 4031
	CodeInsufficientBalance          Code = 4032
	CodeBadAmount                    Code = 4033
	CodeTransactionPairNotFound      Code = 4034
	CodeCurrencyNotFound             Code = 4035
	CodeBadCalcOfTransaction         Code = 4036
	CodeBadTransactionPairStatus     Code = 4037
	CodeTemporaryForbidden           Code = 4038
	CodeBadCSVContent                Code = 4039
	CodeMobileBarcodeNotExist        Code = 4040
	CodeBadCryptocurrencyAddress     Code = 4041
	CodePhoneNumberDuplicated        Code = 4042
	CodeOverWithdrawalWhitelistLimit Code = 4043
	CodeTooManyFiles                 Code = 4044
)

// Code List for HTTP Code 5xx Series
const (
	CodeSql                            Code = 5000
	CodeJWTIssueToken                  Code = 5001
	CodeGetClaims                      Code = 5002
	CodeBadClaimsType                  Code = 5003
	CodeSendEmail                      Code = 5004
	CodeJWTBadPayload                  Code = 5005
	CodeEncryption                     Code = 5006
	CodeGormScan                       Code = 5007
	CodeGormValue                      Code = 5008
	CodeNewHTTPRequest                 Code = 5009
	CodeDoHTTPRequest                  Code = 5010
	CodeReadHTTPResponse               Code = 5011
	CodeSendSMS                        Code = 5012
	CodeJSONMarshal                    Code = 5013
	CodeJSONUnmarshal                  Code = 5014
	CodeBadAmountOfFundsInTransactions Code = 5015
	CodeBadRecord                      Code = 5016
	CodeSaveUploadedFile               Code = 5017
	CodeRedis                          Code = 5018
	CodeRedisBadScript                 Code = 5019
	CodeCallBinanceAPI                 Code = 5020
	CodeParseURL                       Code = 5021
	CodeCallMaxAPI                     Code = 5022
	CodeJSONDecode                     Code = 5023
	CodeParseIP                        Code = 5024
	CodeLookUpMMDB                     Code = 5025
	CodeCallEZReceiptAPI               Code = 5026
	CodeLoginEZReceiptAPI              Code = 5027
	CodeBadParamEZReceiptAPI           Code = 5028
	CodeWriteCSV                       Code = 5029
	CodeCallKryptoGOAPI                Code = 5030
	CodeExecuteTemplate                Code = 5031
	CodeUpdateCache                    Code = 5032
	CodeCallCybavoAPI                  Code = 5033
	CodeCybavoWalletNotFound           Code = 5034
	CodeWalletAddressNotGen            Code = 5035
	CodeWalletAddressAlreadySet        Code = 5036
	CodeGenQRCode                      Code = 5037
	CodeAddressDeploying               Code = 5038
	CodeMemoryCacheError               Code = 5039
	CodeBadBase32String                Code = 5040
	CodeGenerateHOTP                   Code = 5041
)

// Others
const (
	CodeBadErrorType             Code = 9993 // Use general error instead of *Error
	CodeFailedToGenerateID       Code = 9994
	CodeScheduleJobKeyDuplicated Code = 9995
	CodeBadCoding                Code = 9996
	CodeNotInit                  Code = 9997
	CodeNotImplement             Code = 9998
	CodeUnknown                  Code = 9999
)
