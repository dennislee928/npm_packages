package redis

// Define All KeyFormat here, and check it is duplicate or not.

// UsersToken by Users ID
const UsersTokenKeyFormat = "users-token-%d" // #nosec G101

// UserAttemptLogin by Account
const UserAttemptLoginKeyFormat = "users-login-%s"

// ManagersToken by Managers ID
const ManagersTokenKeyFormat = "managers-token-%d" // #nosec G101

// PreverificationKeyFormat by jwt.Type, id and usage
const PreverificationKeyFormat = "preverify-%d-%d-%s"

// VerificationKeyFormat by jwt.Type, id and usage
const VerificationKeyFormat = "verify-%d-%d-%s"
