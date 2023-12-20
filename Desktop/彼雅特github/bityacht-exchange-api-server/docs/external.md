# BitYacht Exchange API External Doc

###### tags: `skycloud`

## Error Code List

| Error Code | Description                         |
| ---------- | ----------------------------------- |
| 4000       | Bind Body Error                     |
| 4001       | Unauthorized                        |
| 4002       | JWT: Bad Signing                    |
| 4003       | JWT: Invalid                        |
| 4004       | JWT: Revoked                        |
| 4005       | Bad Authorization Token             |
| 4006       | Permission Denied                   |
| 4007       | Account(Email) Duplicated           |
| 4008       | Bad URL Parameter                   |
| 4009       | Record Not Found (in db)            |
| 4010       | Record No Change                    |
| 4011       | Bad Body                            |
| 4012       | Account Not Available               |
| 4013       | Insufficient Funds                  |
| 4014       | Bad Query                           |
| 4015       | Fund Deleted                        |
| 4016       | Too Many Requests                   |
| 4017       | Bad Password Strength               |
| 4018       | Over Reset Password Limit           |
| 4019       | Bad Uploaded File Type              |
| 4020       | File Not Found                      |
| 4021       | Email Already Verified              |
| 4022       | Bad Verification Code               |
| 4023       | Verification Code Expired           |
| 4024       | Bad Mobile Barcode Format           |
| 4025       | Bad Invite Code                     |
| 4026       | Image Over Size                     |
| 4027       | Bad Data URI Image Format           |
| 4028       | Bad Image Data                      |
| 4029       | Bad National ID                     |
| 4030       | National ID Duplicated              |
| 4031       | Bad Action                          |
| 4032       | Insufficient Balance                |
| 4033       | Bad Amount                          |
| 4034       | Transaction Pair Not Found          |
| 4035       | Currency Not Found                  |
| 4036       | Bad Calculation Of Transaction      |
| 4037       | Bad Transaction Pair Status         |
| 4038       | Temporary Forbidden                 |
| 4039       | Bad CSV Content                     |
| 4040       | Mobile Barcode Not Exist            |
| 4041       | Bad Cryptocurrency Address          |
| 4042       | Phone Number Duplicated             |
| 4043       | Over Withdrawal Whitelist Limit     |
| 4044       | Too Many Files                      |
| ---------- | -----------------------             |
| 5000       | SQL Error                           |
| 5001       | JWT: Issue Token Error              |
| 5002       | Get Claims Error                    |
| 5003       | Bad Claims Type                     |
| 5004       | Send Email Error                    |
| 5005       | JWT: BadPayload                     |
| 5006       | Encryption                          |
| 5007       | Gorm Scan                           |
| 5008       | Gorm Value                          |
| 5009       | New HTTP Request                    |
| 5010       | Do HTTP Request                     |
| 5011       | Read HTTP Response                  |
| 5012       | Send SMS                            |
| 5013       | JSON Marshal                        |
| 5014       | JSON Unmarshal                      |
| 5015       | Bad Amount Of Funds In Transactions |
| 5016       | Bad Action In Transactions          |
| 5017       | Save Uploaded File                  |
| 5018       | Redis                               |
| 5019       | Redis Bad Script                    |
| 5020       | Call Binance API                    |
| 5021       | Parse URL                           |
| 5022       | Call Max API                        |
| 5023       | JSON Decode                         |
| 5024       | Parse IP                            |
| 5025       | Look Up MMDB                        |
| 5026       | Call EZ Receipt API                 |
| 5027       | Login EZ Receipt API                |
| 5028       | Bad Param EZ Receipt API            |
| 5029       | Write CSV                           |
| 5030       | Call Krypto GO API                  |
| 5031       | Execute Template                    |
| 5032       | Update Cache                        |
| 5033       | Call Cybavo API                     |
| 5034       | Cybavo Wallet Not Found             |
| 5035       | Wallet Address Not Generate         |
| 5036       | Wallet Address Already Set          |
| 5037       | Gen QRCode                          |
| 5038       | ETH address is deploying            |
| 5039       | Memory Cache Error                  |
| 5040       | Bad Base32 String                   |
| 5041       | Generate HOTP                       |
| ---------- | -----------------------             |
| 9993       | Bad Error Type                      |
| 9994       | Failed to Generate ID               |
| 9995       | Schedule Job Key Duplicated         |
| 9996       | Bad Coding                          |
| 9997       | Not Init                            |
| 9998       | Not Implement                       |
| 9999       | Unknown                             |

## JWT

### Common Fields

| Field      | Description                     |
| ---------- | ------------------------------- |
| exp        | Expiration Time                 |
| iat        | Issued At                       |
| jti        | JWT ID                          |
| claimsType | [See Claims Type](#Claims-Type) |
| id         | Managers / Users ID             |

#### Claims Type

| Value | Description |
| ----- | ----------- |
| 1     | Manager     |
| 2     | User        |

### Manager Fields

| Field           | Description      |
| --------------- | ---------------- |
| managersRolesID | Manager Roles ID |
| name            | Name             |

#### Example

Token: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTA1MTEzNzMsImlhdCI6MTY5MDUwNzc3MywianRpIjoiY2I2ZTcxNTYtMWMwNi00ZmRmLTg1ZTQtMzY3NWFkNjM5ODdkIiwiY2xhaW1zVHlwZSI6MSwibWFuYWdlcnNSb2xlc0lEIjoxLCJpZCI6MSwibmFtZSI6IkFkbWluIn0.15ieVSdewP96rJQ29Ftc0fcM020eeEMXH6R6N48C3pg`

Meaning:

```json
{
  "exp": 1690511373,
  "iat": 1690507773,
  "jti": "cb6e7156-1c06-4fdf-85e4-3675ad63987d",
  "claimsType": 1,
  "managersRolesID": 1,
  "id": 1,
  "name": "Admin"
}
```

### Users Fields

| Field         | Description                                                                         |
| ------------- | ----------------------------------------------------------------------------------- |
| account       | Email Account                                                                       |
| countriesCode | Country Code [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) |
| type          | [See Users Type](#Users-Type)                                                       |
| firstName     | First Name                                                                          |
| lastName      | Last Name                                                                           |
| level         | Level                                                                               |
| status        | [See Users Status](#Users-Status)                                                   |

#### Example

Token: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTA1MTE0NzMsImlhdCI6MTY5MDUwNzg3MywianRpIjoiOGY4MTJiYmEtN2U5Ni00YTc4LWI1MDktYjBiNWEyM2MzOGM2IiwiY2xhaW1zVHlwZSI6MiwiaWQiOjMyMjI4MDI2LCJhY2NvdW50IjoiRmFrZVVzZXIwQHRlc3QuY29tLnR3IiwiY291bnRyaWVzQ29kZSI6IiIsInR5cGUiOjEsImZpcnN0TmFtZSI6IiIsImxhc3ROYW1lIjoiIiwibGV2ZWwiOjAsInN0YXR1cyI6MX0.41d2i3VVbs_jM1bEsb8Z2d5nzt_5w_793bs1SIkDdTQ`

Meaning:

```json
{
  "exp": 1690511473,
  "iat": 1690507873,
  "jti": "8f812bba-7e96-4a78-b509-b0b5a23c38c6",
  "claimsType": 2,
  "id": 32228026,
  "account": "FakeUser0@test.com.tw",
  "countriesCode": "",
  "type": 1,
  "firstName": "",
  "lastName": "",
  "level": 0,
  "status": 1
}
```

## Other Types

### Users

#### Users Type

| Value | Description      |
| ----- | ---------------- |
| 1     | Natural Person   |
| 2     | Juridical Person |

#### Users Gender

| Value | Description |
| ----- | ----------- |
| 0     | Unknown     |
| 1     | Male        |
| 2     | Female      |
| 3     | X           |

#### Users Tax Residence

| Value | Description            |
| ----- | ---------------------- |
| 1     | Only Taiwan            |
| 2     | Not Or Not Only Taiwan |

#### Users Status

| Value | Description |
| ----- | ----------- |
| 0     | Unverified  |
| 1     | Enable      |
| 2     | Disable     |
| 3     | Forzen      |

#### Users Login2FAType (Bitwise)

| Value | Description |
| ----- | ----------- |
| 0     | None        |
| 1     | Email       |
| 2     | Phone       |

#### Example

1. Value: 0, Meaning: None
1. Value: 1, Meaning: Email Only
1. Value: 3, Meaning: Email + Phone
