package verifications

import (
	"github.com/redis/go-redis/v9"
)

// KEYS[1] = Verification Key, ARGV[1] = Verification Code
var verifyScript = redis.NewScript(`
local code = redis.call("GET", KEYS[1])

if not code then
	return -1
end

if code ~= ARGV[1] then
	return -2
end

redis.call("DEL", KEYS[1])

return 0
`)
