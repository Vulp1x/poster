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
	Attribute("record", ArrayOf(String), func() {
		MinLength(4)
		MaxLength(4)
	})
	Attribute("line_number", Int, "номер строки в исходном файле", func() {
		Meta("struct:tag:json", "line_number")
	})
	Required("record", "line_number")
})

// TargetUserRecord описывает пользователя, которому будет показана реклама
var TargetUserRecord = Type("TargetUserRecord", func() {
	Attribute("record", ArrayOf(String), func() {
		MinLength(2)
		MaxLength(2)
	})
	Attribute("line_number", Int, "номер строки в исходном файле", func() {
		Meta("struct:tag:json", "line_number")
	})
	Required("record", "line_number")
})

// ProxyRecord описывает пользователя, которому будет показана реклама
var ProxyRecord = Type("ProxyRecord", func() {
	Attribute("record", ArrayOf(String), func() {
		MinLength(4)
		MaxLength(4)
	})
	Attribute("line_number", Int, "номер строки в исходном файле", func() {
		Meta("struct:tag:json", "line_number")
	})
	Required("record", "line_number")
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

// Task описывает рекламную кампанию
var Task = Type("Task", func() {
	Attribute("id", String, "", func() {
		Format(FormatUUID)
	})

	Attribute("text_template", String, func() {
		Meta("struct:tag:json", "text_template")
		Description("описание под постом")
	})

	Attribute("image", String, "base64 строка картинки")
	Attribute("status", Int)
	Attribute("title", String, "название задачи")

	Attribute("bots_num", Int, "количество ботов в задаче", func() {
		Meta("struct:tag:json", "bots_num")
	})
	Attribute("proxies_num", Int, "количество прокси в задаче", func() {
		Meta("struct:tag:json", "proxies_num")
	})
	Attribute("targets_num", Int, "количество целевых пользователей в задаче", func() {
		Meta("struct:tag:json", "targets_num")
	})

	Attribute("bots_filename", String, "название файла, из которого брали ботов", func() {
		Meta("struct:tag:json", "bots_filename")
	})
	Attribute("proxies_filename", String, "название файла, из которого брали прокси", func() {
		Meta("struct:tag:json", "proxies_filename")
	})
	Attribute("targets_filename", String, "название файла, из которого брали целевых пользователей", func() {
		Meta("struct:tag:json", "targets_filename")
	})

	Required("id", "text_template", "image", "status", "title", "bots_num", "proxies_num", "targets_num")
})

var TaskFilenames = Type("TaskFileNames", func() {
	Attribute("bots_filename", String, "название файла, из которого брали ботов", func() {
		Meta("struct:tag:json", "bots_filename")
	})
	Attribute("proxies_filename", String, "название файла, из которого брали прокси", func() {
		Meta("struct:tag:json", "proxies_filename")
	})
	Attribute("targets_filename", String, "название файла, из которого брали целевых пользователей", func() {
		Meta("struct:tag:json", "targets_filename")
	})

	Required("bots_filename", "proxies_filename", "targets_filename")
})

var BotsProgress = Type("BotsProgress", func() {
	Attribute("user_name", String, "имя пользователя бота", func() {
		Meta("struct:tag:json", "user_name")
	})
	Attribute("posts_count", Int, "количество выложенных постов", func() {
		Meta("struct:tag:json", "posts_count")
	})

	Required("user_name", "posts_count")
})
