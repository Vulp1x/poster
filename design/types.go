// nolint
package design

import (
	. "goa.design/goa/v3/dsl"
)

// Creds defines the credentials to use for authenticating to service methods.
var Creds = Type("Creds", func() {
	Field(1, "jwt", String, "JWT token", func() {
		Example("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
	})
	Required("jwt")
})

// BotAccountRecord defines set of bot's fields and line in input file
var BotAccountRecord = Type("BotAccountRecord", func() {
	Attribute("username", String, "login")
	Attribute("password", String, "login")
	Attribute("user_agent", String, "user agent header", func() {
		Meta("struct:tag:json", "user_agent")
	})
	Attribute("device_id", String, "main id, ex: android-0d735e1f4db26782", func() {
		Meta("struct:tag:json", "device_id")
	})
	Attribute("uuid", String, func() {
		Format(FormatUUID)
	})
	Attribute("phone_id", String, "phone_id", func() {
		Meta("struct:tag:json", "phone_id")
		Format(FormatUUID)
	})
	Attribute("advertising_id", String, "adv id", func() {
		Meta("struct:tag:json", "advertising_id")
		Format(FormatUUID)

	})
	Attribute("family_device_id", String, func() {
		Meta("struct:tag:json", "family_device_id")
		Format(FormatUUID)
	})

	Attribute("headers", MapOf(String, String))

	Required(
		"username", "password", "user_agent", "device_id", "uuid",
		"phone_id", "advertising_id", "family_device_id", "headers",
	)
})

// TargetUser описывает пользователя, которому будет показана реклама
var TargetUser = Type("TargetUser", func() {
	Attribute("username", String, "instagram username")
	Attribute("user_id", Int64, "instagram user id", func() {
		Meta("struct:tag:json", "user_id")
	})

	Required("username", "user_id")
})

// Proxy описывает пользователя, которому будет показана реклама
var Proxy = Type("Proxy", func() {
	Attribute("host", String, "адрес прокси")
	Attribute("port", Int64, "номер порта")

	Attribute("login", String)
	Attribute("password", String)

	Required("host", "port", "login", "password")
})

// UploadError описывает ошибку при чтении файлов
var UploadError = Type("UploadError", func() {
	Attribute("type", Int, "тип файла, в котором была ошибка", func() {
		Enum(1, 2, 3)
		Description(` 1 - список ботов
    2 - список прокси
    3 - список получателей рекламы`)
	})
	Attribute("line", Int)
	Attribute("input", String, "номер порта")

	Attribute("reason", String)

	Required("type", "line", "input", "reason")
})

//
// // DeviceSettings defines device headers
// var DeviceSettings = Type("DeviceSettings", func() {
// 	Attribute("app_version", String, "app version", func() {
// 		Meta("struct:tag:json", "app_version")
// 	})
// 	Attribute("android_version", Int, "android_version", func() {
// 		Meta("struct:tag:json", "android_version")
// 	})
// 	Attribute("android_release", String, "user agent header", func() {
// 		Meta("struct:tag:json", "android_release")
// 	})
// 	Attribute("dpi", String)
// 	Attribute("resolution", String)
// 	Attribute("manufacturer", String)
// 	Attribute("device", String)
// 	Attribute("model", String)
// 	Attribute("cpu", String)
// 	Attribute("version_code", String, "version", func() {
// 		Meta("struct:tag:json", "version_code")
// 	})
//
// 	Required("app_version", "android_version", "android_release", "dpi", "resolution", "manufacturer",
// 		"device", "model", "cpu", "version_code")
// })
