package token

import (
	"github.com/redis/go-redis/v9"
)

// KEYS[1] = User or Manager's Attempt Login Key, ARGV[1] = loggedIn(bool), ARGV[2] = attemptLoginExp, ARGV[3] = lockLoginExp
var updateAttemptLoginScript = redis.NewScript(`
if ARGV[1] == "1" then
	redis.call("DEL", KEYS[1])
	return 0
end

local record = redis.call("GET", KEYS[1])
local ex = tonumber(ARGV[2])
local output = 0

if not record then
	record = { c = 1, l = false }
else
	record = cjson.decode(record)

	record["c"] = record["c"] + 1
	if record["c"] >= 5 then
		record["l"] = true
		ex = tonumber(ARGV[3])
		output = 1
	end
end
redis.call("SET", KEYS[1], cjson.encode(record), "EX", ex)

return output
`)

// KEYS[1] = User or Manager's Token Key, ARGV[1] = JWT ID (jti)
var logoutScript = redis.NewScript(`
local record = redis.call("GET", KEYS[1])

if not record then
	return -1
end

record = cjson.decode(record)
if record["Token"] ~= ARGV[1] then
	return -2
end

redis.call("DEL", KEYS[1])

return 0
`)

// KEYS[1] = User or Manager's Token Key, ARGV[1] = Old Refresh Token, ARGV[2] = New Record
var refreshTokenScript = redis.NewScript(`
local record = redis.call("GET", KEYS[1])

if not record then
	return -1
end

record = cjson.decode(record)

if record["RT"] ~= ARGV[1] then
	return -2
end

local newRecord = cjson.decode(ARGV[2])
newRecord["LA"] = record["LA"]

newRecord = cjson.encode(newRecord)
redis.call("SET", KEYS[1], newRecord)

return newRecord
`)
