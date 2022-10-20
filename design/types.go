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

// TaskStatus описывает статус задачи
var TaskStatus = Type("TaskStatus", Int, func() {
	Enum(1, 2, 3, 4, 5, 6)
	Description(`1 - задача только создана, нужно загрузить список ботов, прокси и получателей
	2- в задачу загрузили необходимые списки, нужно присвоить прокси для ботов
	3- задача готова к запуску
	4- задача запущена 
	5 - задача остановлена
	6 - задача завершена`)
})

// Task описывает рекламную кампанию
var Task = Type("Task", func() {
	Attribute("id", String, "", func() {
		Format(FormatUUID)
	})

	Attribute("text_template", String, func() {
		Meta("struct:tag:json", "text_template")
		Description("описание под постом")
	})

	Attribute("post_images", ArrayOf(String), "список base64 строк картинок", func() {
		Meta("struct:tag:json", "post_images")
	})

	Attribute("landing_accounts", ArrayOf(String), func() {
		Description("имена аккаунтов, на которых ведем трафик")
		Meta("struct:tag:json", "landing_accounts")
	})

	Attribute("bot_names", ArrayOf(String), func() {
		Description("имена для аккаунтов-ботов")
		Meta("struct:tag:json", "bot_names")
	})

	Attribute("bot_last_names", ArrayOf(String), func() {
		Description("фамилии для аккаунтов-ботов")
		Meta("struct:tag:json", "bot_last_names")
	})

	Attribute("bot_images", ArrayOf(String), func() {
		Description("аватарки для ботов")
		Meta("struct:tag:json", "bot_images")
	})

	Attribute("bot_urls", ArrayOf(String), func() {
		Description("ссылки для описания у ботов")
		Meta("struct:tag:json", "bot_images")
	})
	Attribute("status", TaskStatus)
	Attribute("title", String, "название задачи")

	Attribute("bots_num", Int, "количество ботов в задаче", func() {
		Meta("struct:tag:json", "bots_num")
	})
	Attribute("residential_proxies_num", Int, "количество резидентских прокси в задаче", func() {
		Meta("struct:tag:json", "residential_proxies_num")
	})
	Attribute("cheap_proxies_num", Int, "количество дешёвых прокси в задаче", func() {
		Meta("struct:tag:json", "cheap_proxies_num")
	})

	Attribute("targets_num", Int, "количество целевых пользователей в задаче", func() {
		Meta("struct:tag:json", "targets_num")
	})

	Attribute("bots_filename", String, "название файла, из которого брали ботов", func() {
		Meta("struct:tag:json", "bots_filename")
	})
	Attribute("residential_proxies_filename", String, "название файла, из которого брали резидентские прокси", func() {
		Meta("struct:tag:json", "residential_proxies_filename")
	})

	Attribute("cheap_proxies_filename", String, "название файла, из которого брали дешёвые прокси", func() {
		Meta("struct:tag:json", "cheap_proxies_filename")
	})
	Attribute("targets_filename", String, "название файла, из которого брали целевых пользователей", func() {
		Meta("struct:tag:json", "targets_filename")
	})

	Required("id", "text_template", "post_images", "status", "title", "bots_num", "residential_proxies_num",
		"cheap_proxies_num", "targets_num", "bot_images", "landing_accounts", "bot_names", "bot_last_names", "bot_urls",
	)
})

var TaskFilenames = Type("TaskFileNames", func() {
	Attribute("bots_filename", String, "название файла, из которого брали ботов", func() {
		Meta("struct:tag:json", "bots_filename")
	})
	Attribute("residential_proxies_filename", String, "название файла, из которого брали резидентские прокси", func() {
		Meta("struct:tag:json", "residential_proxies_filename")
	})
	Attribute("cheap_proxies_filename", String, "название файла, из которого брали дешёвые прокси", func() {
		Meta("struct:tag:json", "cheap_proxies_filename")
	})
	Attribute("targets_filename", String, "название файла, из которого брали целевых пользователей", func() {
		Meta("struct:tag:json", "targets_filename")
	})

	Required("bots_filename", "residential_proxies_filename", "cheap_proxies_filename", "targets_filename")
})

var BotsProgress = Type("BotsProgress", func() {
	Attribute("user_name", String, "имя пользователя бота", func() {
		Meta("struct:tag:json", "user_name")
	})
	Attribute("posts_count", Int, "количество выложенных постов", func() {
		Meta("struct:tag:json", "posts_count")
	})
	Attribute("status", Int, "текущий статус бота, будут ли выкладываться посты")

	Required("user_name", "posts_count", "status")
})
