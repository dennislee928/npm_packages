package kyc

type IDType string

const (
	IDTypePASSPORT       IDType = "PASSPORT"
	IDTypeDrivingLicense IDType = "DRIVING_LICENSE"
	IDTypeIDCard         IDType = "ID_CARD"
)

type AuditStatus string

const (
	AuditStatusPending  AuditStatus = "Pending"
	AuditStatusAccepted AuditStatus = "Accepted"
	AuditStatusRejected AuditStatus = "Rejected"
)

type Workflow int32

const (
	WorkflowIDOnlyCameraAndUpload        Workflow = 100
	WorkflowIDOnlyCameraOnly             Workflow = 101
	WorkflowIDOnlyUploadOnly             Workflow = 102
	WorkflowIDAndIdentityCameraAndUpload Workflow = 200
	WorkflowIDAndIdentityCameraOnly      Workflow = 201
	WorkflowIDAndIdentityUploadOnly      Workflow = 202
)

type IDVState int32

const (
	IDVStateAccept  IDVState = 0
	IDVStateReview  IDVState = 1
	IDVStateReject  IDVState = 2
	IDVStatePending IDVState = 3
	IDVStateInitial IDVState = 4
)

// Ref: https://www.kryptogo.com/docs/kyabc#tag/ID-Verification/paths/~1idv/post
type CreateIDVTaskRequest struct {
	// The country of the ID. Should be ISO-3166 Alpha-3
	Country string `json:"country"`

	// The type of ID. Supports PASSPORT, DRIVING_LICENSE and ID_CARD
	IDType IDType `json:"id_type"`

	// The base64 encoding of ID image (front side). Max. 15MB & <8000 pixels per side & longest side >300 pixels
	IDImage string `json:"id_image"`

	// (Optional) Mime type of ID image. Can be image/jpeg(default) or image/png
	IDMimeType string `json:"id_mime_type,omitempty"`

	// (Optional) The base64 encoding of ID image (back side). Required when id_type is ID_CARD. Max. 15MB & <8000 pixels per side & longest side >300 pixels
	IDBackImage string `json:"id_back_image,omitempty"`

	// (Optional) Mime type of ID image back side. Can be image/jpeg(default) or image/png
	IDBackMimeType string `json:"id_back_mime_type,omitempty"`

	// The base64 encoding of face image. At least one of face_image or id_and_face_image should be provided. Max. 15MB & <8000 pixels per side & longest side >300 pixels
	FaceImage string `json:"face_image"`

	// (Optional) Mime type of face image. Can be image/jpeg(default) or image/png
	FaceMimeType string `json:"face_mime_type,omitempty"`

	// The base64 encoding of image of the person holding ID document. At least one of face_image or id_and_face_image should be provided.
	IDAndFaceImage string `json:"id_and_face_image"`

	// (Optional) Mime type of image of the person holding ID document. Can be image/jpeg(default) or image/png
	IDAndFaceMimeType string `json:"id_and_face_mime_type,omitempty"`

	// (Optional) The expected full name of the ID holder. If it doesn't match the recognized name on the ID, the result state will be review
	ExpectedName string `json:"expected_name,omitempty"`

	// (Optional) The expected birthday of the ID holder in YYYY-MM-DD format. If it doesn't match the recognized birthday on the ID, the result state will be review
	ExpectedBirthday string `json:"expected_birthday,omitempty"`

	// (Optional) The expected ID number of the ID holder. If it doesn't match the recognized ID number on the ID, the result state will be review
	ExpectedIDNumber string `json:"expected_id_number,omitempty"`

	// (Optional) The expected issuing Date of the ID holder. If it doesn't match the recognized ID number on the ID, the result state will be review
	ExpectedIssuingDate string `json:"expected_issuing_date,omitempty"`

	// (Optional) The expected expiry Date of the ID holder. If it doesn't match the recognized ID number on the ID, the result state will be review
	ExpectedExpiryDate string `json:"expected_expiry_date,omitempty"`

	// (Optional) The expected gender of the ID holder. If it doesn't match the recognized ID number on the ID, the result state will be review
	ExpectedGender string `json:"expected_gender,omitempty"`

	// (Optional) The URL which will be called when the IDV task is completed
	CallbackURL string `json:"callback_url,omitempty"`

	// (Optional) The internal customer reference of your system
	CustomerReference string `json:"customer_reference,omitempty"`

	// (Optional) Whether to automatically create DD search task for the target if the result state is accept (default is false). If the result state is review or reject, the DD search task won't be created even this field is true
	AutoCreateDDTask bool `json:"auto_create_dd_task,omitempty"`

	// (Optional) The callback URL to use if a DD search task is created from this IDV task
	DDTaskCallbackURL string `json:"dd_task_callback_url,omitempty"`

	// (Optional) The ID of search setting to use if a DD search task is created from this IDV task. Required when a DD search task is created
	DDTaskSearchSettingID int64 `json:"dd_task_search_setting_id,omitempty"`
}

// Ref: https://www.kryptogo.com/docs/kyabc#tag/ID-Verification/paths/~1idv~1init/post
type InitIDVRequest struct {
	// (Optional) The type of ID. Supports PASSPORT, DRIVING_LICENSE and ID_CARD
	IDType IDType `json:"id_type,omitempty"`

	// (Optional) Renders web page in the specified language. Supported locale values:
	// ar Arabic
	// bg Bulgarian
	// cs Czech
	// da Danish
	// de German
	// el Greek
	// en American English (default)
	// en-GB British English
	// es Spanish
	// es-MX Mexican Spanish
	// et Estonian
	// fi Finnish
	// fr French
	// he Hebrew
	// hr Croatian
	// hu Hungarian
	// hy Armenian
	// id Indonesian
	// it Italian
	// ja Japanese
	// ka Georgian
	// km Khmer
	// ko Korean
	// lt Lithuanian
	// ms Malay
	// nl Dutch
	// no Norwegian
	// pl Polish
	// pt Portuguese
	// pt-BR Brazilian Portuguese
	// ro Romanian
	// ru Russian
	// sk Slovak
	// sv Swedish
	// sw Swahili
	// th Thai
	// tr Turkish
	// vi Vietnamese
	// zh-CN Simplified Chinese
	// zh-HK Traditional Chinese
	Locale string `json:"locale,omitempty"`

	// Applies this acquisition workflow to the transaction. Supported workflowId values:
	// 100 ID only, camera + upload
	// 101 ID only, camera only
	// 102 ID only, upload only
	// 200 ID + Identity, camera + upload
	// 201 ID + Identity, camera only
	// 202 ID + Identity, upload only
	WorkflowID Workflow `json:"workflow_id,omitempty"`

	// (Optional) Redirects to this URL after a successful transaction. Max length is 2047
	SuccessURL string `json:"success_url,omitempty"`

	// (Optional) Redirects to this URL after an unsuccessful transaction. Max length is 255
	ErrorURL string `json:"error_url,omitempty"`

	// The country of the ID. Should be ISO-3166 Alpha-3
	Country string `json:"country,omitempty"`

	// (Optional) The expected full name of the ID holder. If it doesn't match the recognized name on the ID, the result state will be review
	ExpectedName string `json:"expected_name,omitempty"`

	// (Optional) The expected birthday of the ID holder in YYYY-MM-DD format. If it doesn't match the recognized birthday on the ID, the result state will be review
	ExpectedBirthday string `json:"expected_birthday,omitempty"`

	// (Optional) The expected ID number of the ID holder. If it doesn't match the recognized ID number on the ID, the result state will be review
	ExpectedIDNumber string `json:"expected_id_number,omitempty"`

	// (Optional) The URL which will be called when the IDV task is completed
	CallbackURL string `json:"callback_url,omitempty"`

	// The internal customer reference of your system. Max length is 100
	CustomerReference string `json:"customer_reference,omitempty"`

	// (Optional) Whether to automatically create DD search task for the target if the result state is accept (default is false). If the result state is review or reject, the DD search task won't be created even this field is true
	AutoCreateDDTask bool `json:"auto_create_dd_task,omitempty"`

	// (Optional) The callback URL to use if a DD search task is created from this IDV task
	DDTaskCallbackURL string `json:"dd_task_callback_url,omitempty"`

	// (Optional) The ID of search setting to use if a DD search task is created from this IDV task. Required when a DD search task is created
	DDTaskSearchSettingID int64 `json:"dd_task_search_setting_id,omitempty"`

	// (Support when id_type = ID_CARD and country = TWN) Not Issue = 0, Initial Issue = 1, Reissue = 2, Change = 3
	TwIDApplyCode int32 `json:"tw_id_apply_code,omitempty"`

	// (Support when id_type = ID_CARD and country = TWN) Issuing date, Republic of China (ROC) format, 7 digits.
	TwIDApplyDate string `json:"tw_id_apply_date,omitempty"`

	// (Support when id_type = ID_CARD and country = TWN) Issuing location administrative area code, 5 digits, see the Ministry of the Interior API documentation appendix 'Issuing location code comparison' for details.
	// 10001 北縣
	// 10002 宜縣
	// 10003 桃縣
	// 10004 竹縣
	// 10005 苗縣
	// 10006 中縣
	// 10007 彰縣
	// 10008 投縣
	// 10009 雲縣
	// 10010 嘉縣
	// 10011 南縣
	// 10012 高縣
	// 10013 屏縣
	// 10014 東縣
	// 10015 花縣
	// 10016 澎縣
	// 10017 基市
	// 10018 竹市
	// 10020 嘉市
	// 09007 連江
	// 09020 金門
	// 63000 北市
	// 64000 高市
	// 65000 新北市
	// 66000 or 10019 中市
	// 67000 or 10012 南市
	// 68000 桃市
	TwIDIssueSiteID string `json:"tw_id_issue_site_id,omitempty"`
}

// Ref: https://www.kryptogo.com/docs/kyabc#tag/ID-Verification/paths/~1idv~1%7Bidv_id%7D/get
type IDVTaskDetail struct {
	// The unique ID of this IDV task
	IDVTaskID string `json:"idv_task_id,omitempty"`

	// Same as IDV task creation request
	Country string `json:"country,omitempty"`

	// Same as IDV task creation request
	IDType IDType `json:"id_type,omitempty"`

	// Possible id_sub_type if type = ID_CARD
	// 		NATIONAL_ID
	// 		RESIDENT_PERMIT_ID
	// Possible id_sub_type if id_type = DRIVING_LICENSE
	// 		REGULAR_DRIVING_LICENSE
	// 		LEARNING_DRIVING_LICENSE
	// Possible id_sub_type if idType = PASSPORT
	// 		E_PASSPORT
	IDSubType string `json:"id_sub_type,omitempty"`

	// The URL of the uploaded ID image (front side)
	IDImageURL string `json:"id_image_url,omitempty"`

	// The URL of the uploaded ID image (back side)
	IDBackImageURL string `json:"id_back_image_url,omitempty"`

	// The URL of the uploaded face image
	FaceImageURL string `json:"face_image_url,omitempty"`

	// The URL of the uploaded ID with face image
	IDAndFaceImageURL string `json:"id_and_face_image_url,omitempty"`

	// Same as IDV task creation request
	ExpectedName string `json:"expected_name,omitempty"`

	// Same as IDV task creation request
	ExpectedBirthday string `json:"expected_birthday,omitempty"`

	// Same as IDV task creation request
	ExpectedIDNumber string `json:"expected_id_number,omitempty"`

	// Same as IDV task creation request
	ExpectedCountry string `json:"expected_country,omitempty"`

	// Same as IDV task creation request
	ExpectedIssuingDate string `json:"expected_issuing_date,omitempty"`

	// Same as IDV task creation request
	ExpectedExpiryDate string `json:"expected_expiry_date,omitempty"`

	// Same as IDV task creation request
	ExpectedGender string `json:"expected_gender,omitempty"`

	// Same as IDV task creation request
	CustomerReference string `json:"customer_reference,omitempty"`

	// First name recognized from the ID image. Will be null if state is pending.
	FirstName string `json:"first_name,omitempty"`

	// Last name recognized from the ID image. Will be null if state is pending.
	LastName string `json:"last_name,omitempty"`

	// Full name recognized from the ID image. Will be null if state is pending.
	FullName string `json:"full_name,omitempty"`

	// Birthday recognized from the ID image in YYYY-MM-DD format. Will be null if state is pending.
	Birthday string `json:"birthday,omitempty"`

	// ID number recognized from the ID image. Will be null if state is pending.
	IDNumber string `json:"id_number,omitempty"`

	// The issuing date recognized from the ID image in YYYY-MM-DD format. Will be null if state is pending.
	IssuingDate string `json:"issuing_date,omitempty"`

	// The expiry date recognized from the ID image in YYYY-MM-DD format. Will be null if state is pending.
	ExpiryDate string `json:"expiry_date,omitempty"`

	// Enum: "M" "F" "X"
	// The gender recognized from the ID image. Will be null if state is pending.
	Gender string `json:"gender,omitempty"`

	// The result state of this IDV task. Possible values are:
	// 0: accept
	// 1: review
	// 2: reject
	// 3: pending
	// 4: initial
	State IDVState `json:"state,omitempty"`

	// Array of objects (IdvTaskResultReviewReason)
	// The reasons why this IDV task needs manual review. Only exists when state is review
	ReviewReasons []idvTaskResultReviewReason `json:"review_reasons,omitempty"`

	// Array of objects (IdvTaskResultRejectReason)
	// The reasons why this IDV task is rejected. Only exists when state is reject
	RejectReasons []idvTaskResultRejectReason `json:"reject_reasons,omitempty"`

	// The DD search task ID if it's automatically created from this IDV task. Will be null if state is pending or no task is create.
	CreatedSearchTaskID int64 `json:"created_search_task_id,omitempty"`

	// The audit result of this IDV task. Possible values are:
	// Pending
	// Accepted
	// Rejected
	AuditStatus AuditStatus `json:"audit_status,omitempty"`

	// The audit time in unix epoch.
	AuditTimestamp int64 `json:"audit_timestamp,omitempty"`

	// The auditor account.
	AuditorEmployeeAccount string `json:"auditor_employee_account,omitempty"`

	// Array of objects (DDTaskInfo)
	// The DD search tasks created by automatically process from API or manually created from control panel.
	DDTasks []ddTaskInfo `json:"dd_tasks,omitempty"`
}

type idvTaskResultReviewReason struct {
	// See below for possible values
	Code int32 `json:"code,omitempty"`

	// All possible review reason's code and description are:
	// 400: Face mismatch
	// 406: Age difference too big
	// 500: Name doesn't match expected result
	// 501: Birthday doesn't match expected result
	// 502: ID number doesn't match expected result
	// 503: Country doesn't match expected result
	// 504: Issuing date doesn't match expected result
	// 505: Expiry date doesn't match expected result
	// 506: Gender doesn't match expected result
	// 600: Selfie photo not provided
	Description string `json:"description,omitempty"`
}

type idvTaskResultRejectReason struct {
	// See below for possible values
	Code int32 `json:"code,omitempty"`

	// All possible reject reason's code and description are:
	// 100: Manipulated document
	// 102: Photocopy black white
	// 103: Photocopy color
	// 104: Digital copy
	// 105: Fraudster
	// 106: Fake
	// 107: Photo mismatch
	// 108: MRZ check failed
	// 109: Punched document
	// 111: Mismatch printed barcode data
	// 200: Not readable document
	// 201: No document
	// 202: Sample document
	// 206: Missing back
	// 207: Wrong document page
	// 209: Missing signature
	// 211: Different persons shown
	// 213: Invalid watermark
	// 300: Manual rejection
	// 401: Selfie cropped from ID
	// 402: Entire ID used as selfie
	// 403: Multiple people
	// 404: Selfie is screen paper video
	// 405: Selfie is manipulated
	// 407: No face present
	// 408: Face not fully visible
	// 409: Bad quality
	// 410: Black and white
	// 411: Liveness check failed
	// 700: Unsupported ID Type
	// 701: Unsupported ID Country
	// 702: No ID Uploaded
	// 801: Authentication Invalid
	// 802: No more tries remaining for this identity. Try again tomorrow
	// 901: Invalid TwAddonIdentity. 2 more tries remaining
	// 901: Invalid TwAddonIdentity. 1 more try remaining
	// 901: Invalid TwAddonIdentity. No more tries remaining
	// 902: TwAddonIdentity Not in circulation
	// 903: TwAddonIdentity Reported Lost
	Description string `json:"description,omitempty"`
}

type ddTaskInfo struct {
	// The unique ID of the DD search task
	ID int64 `json:"id,omitempty"`

	// The created time in unix epoch
	CreationTimestamp int64 `json:"creation_timestamp,omitempty"`

	// The audit status of the DD search task. Possible values are:
	// 0: Undecided
	// 1: Accepted
	// 2: Rejected
	AuditStatus int32 `json:"audit_status,omitempty"`
}

// Ref: https://www.kryptogo.com/docs/kyabc#tag/ID-Verification/paths/~1idv~1init/post
type InitIDVResponse struct {
	IDVTaskID int64  `json:"idv_task_id"` // The unique ID of this IDV task
	URL       string `json:"url"`         // URL used to load the ID Verification web page
	Timestamp string `json:"timestamp"`   // The created time in unix epoch //! Real Result Format is "2006-01-02T15:04:05Z"
}

// Ref: https://www.kryptogo.com/docs/kyabc#tag/DD-KYB-and-KYC/paths/~1task/post
type CreateTasksRequest []Task

type Task struct {
	SearchSettingID int64               `json:"search_setting_id"` // The ID of search setting to use
	Target          searchTaskKycTarget `json:"target"`            // SearchTaskKycTarget (object) or SearchTaskKybTarget (object)

	// The array of source to search (Default is searching all):
	// 0: Negative news
	// 1: Basic information
	// 2: Taiwan company
	// 3: Taiwan verdict
	// 4: Sanction list
	SearchSource      []int32 `json:"search_source"`
	CallBackURL       string  `json:"callback_url"`       // The URL which will be called when the task is completed, or the task's accepted status is channged
	CustomerReference string  `json:"customer_reference"` // The internal customer reference of your system

	// (Optional) ID of certain IDV task, indicates the source of this search task.
	// If from_idv_id is presented, the following data are forced copy from the IDV task:
	// target.birthday, target.citizenship, customer_reference and callback_url.
	FromIDVID int64 `json:"from_idv_id"`
}

// Ref: https://www.kryptogo.com/docs/kyabc#tag/DD-KYB-and-KYC/paths/~1task/post
type CreateTasksResponse []TaskIDAndTimestamp

type TaskIDAndTimestamp struct {
	TaskID    int64 `json:"task_id"`   // The unique ID of this IDV task
	Timestamp int64 `json:"timestamp"` // The created time in unix epoch
}

type UpdateTaskStatusRequest struct {
	Comment  string `json:"comment"`  // The comment of the search task report (usually the reason of accepting/rejecting this target)
	Accepted bool   `json:"accepted"` // Whether to accept the target of this search task
}

type UpdateTaskStatusResponse struct {
	TaskID   int64  `json:"task_id"`  // The unique ID of this IDV task
	Comment  string `json:"comment"`  // The comment of the search task report (usually the reason of accepting/rejecting this target)
	Accepted bool   `json:"accepted"` // Whether to accept the target of this search task
}

type IDVCallbackRequest struct {
	// The unique ID of this IDV task
	IDVTaskID int64 `json:"idv_task_id"`

	// Same as IDV task creation request
	Country string `json:"country"`

	// Same as IDV task creation request
	IDType string `json:"id_type"`

	// Possible id_sub_type if type = ID_CARD
	// 		NATIONAL_ID
	// 		RESIDENT_PERMIT_ID
	// Possible id_sub_type if id_type = DRIVING_LICENSE
	// 		REGULAR_DRIVING_LICENSE
	// 		LEARNING_DRIVING_LICENSE
	// Possible id_sub_type if idType = PASSPORT
	// 		E_PASSPORT
	IDSubType string `json:"id_sub_type"`

	// The URL of the uploaded ID image (front side)
	IDImageURL string `json:"id_image_url"`

	// The URL of the uploaded ID image (back side)
	IDBackImageURL string `json:"id_back_image_url"`

	// The URL of the uploaded face image
	FaceImageURL string `json:"face_image_url"`

	// The URL of the uploaded ID with face image
	IDAndFaceImageURL string `json:"id_and_face_image_url"`

	// Same as IDV task creation request
	ExpectedName string `json:"expected_name"`

	// Same as IDV task creation request
	ExpectedBirthday string `json:"expected_birthday"`

	// Same as IDV task creation request
	ExpectedIDNumber string `json:"expected_id_number"`

	// Same as IDV task creation request
	ExpectedCountry string `json:"expected_country"`

	// Same as IDV task creation request
	ExpectedIssuingDate string `json:"expected_issuing_date"`

	// Same as IDV task creation request
	ExpectedExpiryDate string `json:"expected_expiry_date"`

	// Same as IDV task creation request
	ExpectedGender string `json:"expected_gender"`

	// Same as IDV task creation request
	CustomerReference string `json:"customer_reference"`

	// First name recognized from the ID image. Will be null if state is pending.
	FirstName string `json:"first_name"`

	// Last name recognized from the ID image. Will be null if state is pending.
	LastName string `json:"last_name"`

	// Full name recognized from the ID image. Will be null if state is pending.
	FullName string `json:"full_name"`

	// Birthday recognized from the ID image in YYYY-MM-DD format. Will be null if state is pending.
	Birthday string `json:"birthday"`

	// ID number recognized from the ID image. Will be null if state is pending.
	IDNumber string `json:"id_number"`

	// The issuing date recognized from the ID image in YYYY-MM-DD format. Will be null if state is pending.
	IssuingDate string `json:"issuing_date"`

	// The expiry date recognized from the ID image in YYYY-MM-DD format. Will be null if state is pending.
	ExpiryDate string `json:"expiry_date"`

	// Enum: "M" "F" "X"
	// The gender recognized from the ID image. Will be null if state is pending.
	Gender string `json:"gender"`

	// The result state of this IDV task. Possible values are:
	// 0: accept
	// 1: review
	// 2: reject
	// 3: pending
	// 4: initial
	State IDVState `json:"state"`

	// Array of objects (IdvTaskResultReviewReason)
	// The reasons why this IDV task needs manual review. Only exists when state is review
	ReviewReasons []idvTaskResultReviewReason `json:"review_reasons"`

	// Array of objects (IdvTaskResultRejectReason)
	// The reasons why this IDV task is rejected. Only exists when state is reject
	RejectReasons []idvTaskResultRejectReason `json:"reject_reasons"`

	// The DD search task ID if it's automatically created from this IDV task. Will be null if state is pending or no task is create.
	CreatedSearchTaskID int64 `json:"created_search_task_id"`

	// The audit result of this IDV task. Possible values are:
	// Pending
	// Accepted
	// Rejected
	AuditStatus AuditStatus `json:"audit_status"`

	// The audit time in unix epoch.
	AuditTimestamp int64 `json:"audit_timestamp"`

	// The auditor account.
	AuditorEmployeeAccount string `json:"auditor_employee_account"`

	// Array of objects (DDTaskInfo)
	// The DD search tasks created by automatically process from API or manually created from control panel.
	DDTasks []ddTaskInfo `json:"dd_tasks"`
}

type TaskSummaryResponse = DDCallbackRequest

type DDCallbackRequest struct {
	TaskID int64 `json:"task_id"`

	SearchSetting searchSetting `json:"search_setting"`

	// SearchTaskKycTarget (object) or SearchTaskKybTarget (object)
	Target searchTaskKycTarget `json:"target"`

	// The internal customer reference of your system
	CustomerReference string `json:"customer_reference"`

	// The progress of the search task. 100 means completed.
	Progress int8 `json:"progress"`

	// The task origin. 0 means single search from UI, 1 means batch search from UI, 2 means search from API
	TaskOrigin int32 `json:"task_origin"`

	// Creation time of the search task in unix epoch
	CreationTime int64 `json:"creation_time"`

	// Update time of the search task in unix epoch.
	UpdateTime int64 `json:"update_time"`

	// Audit time of the search task in unix epoch, i.e. when the accepted status is changed
	AuditTime int64 `json:"audit_time"`

	// (SearchTaskInfoReport (object or null))
	Report searchTaskInfoReport `json:"report"`

	// Metadata of the search task
	Metadata []searchTaskMetadata `json:"metadata"`
}

type searchSetting struct {
	// Search mode. 1 means precise search and 2 means fuzzy search
	Mode int32 `json:"mode"`

	// The range of time (in years) for searching web pages. -1 means no constraint.
	TimeRangeYear int32 `json:"time_range_year"`

	// The language used to search information
	Language string `json:"language"`

	// The negative words used in negative news search
	NegativeWords []string `json:"negative_words"`

	Sites []string `json:"sites"`

	// The number of websites to search for negative news
	NumberOfPages int32 `json:"number_of_pages"`
}

type searchTaskKycTarget struct {
	// Define the task is for KYC or KYB. Should be 1 here for KYC task
	Type int32 `json:"type"`

	Name string `json:"name"`

	// (Optional) Birthday in YYYY-MM-DD format
	Birthday string `json:"birthday"`

	// (Optional) The country code of the citizenship. Could be ISO-3166 Alpha-3 or ISO-3166 Numeric
	Citizenship string `json:"citizenship"`
}

type searchTaskInfoReport struct {
	// Whether the target of this search task matches sanction list
	SanctionMatched bool `json:"sanction_matched"`

	// The potential risk of accepting this target
	PotentialRisk int `json:"potential_risk"`

	// The URL for checking the search results of this task
	URL string `json:"url"`

	// Whether to accept this target. Null value means it's not decided yet
	Accepted bool `json:"accepted"`

	// The comment of the search task report (usually the reason of accepting/rejecting this target)
	Comment string `json:"comment"`
}

type searchTaskMetadata struct {
	// The key(name) of the metadata
	Key string `json:"key"`

	// The value of the metadata
	Value string `json:"value"`
}
