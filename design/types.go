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
	Enum(1, 2, 3, 4, 5, 6, 7, 8)
	Description(`1 - задача только создана, нужно загрузить список ботов, прокси и получателей
	2- в задачу загрузили необходимые списки, нужно присвоить прокси для ботов
	3- задача готова к запуску
	4- задача запущена 
	5 - задача остановлена
	6 - выкладывание постов закончено
	7 - меняем описания у выложенных постов
	8 - задача завершена полностью, терминальный статус`)
})

// TaskType описывает тип задачи (фотографии или рилсы)
var TaskType = Type("TaskType", Int, func() {
	Enum(1, 2)
	Description(`1 - загружаем изображения
	2- загружаем видео в рилсы`)
})

// Task описывает рекламную кампанию
var Task = Type("Task", func() {
	Attribute("id", String, "", func() {
		Format(FormatUUID)
	})

	Attribute("type", TaskType)

	Attribute("text_template", String, func() {
		Meta("struct:tag:json", "text_template")
		Description("описание под постом")
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

	Attribute("bot_urls", ArrayOf(String), func() {
		Description("ссылки для описания у ботов")
		Meta("struct:tag:json", "bot_urls")
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

	Attribute("video_filename", String, "название файла с видео, если тип задачи - рилсы", func() {
		Meta("struct:tag:json", "video_filename")
	})

	Attribute("follow_targets", Boolean, func() {
		Description("нужно ли подписываться на аккаунты")
		Meta("struct:tag:json", "follow_targets")
	})

	Attribute("need_photo_tags", Boolean, func() {
		Description("делать отметки на фотографии")
		Meta("struct:tag:json", "need_photo_tags")
	})

	Attribute("per_post_sleep_seconds", UInt, func() {
		Description("задержка между постами")
		Meta("struct:tag:json", "per_post_sleep_seconds")
	})

	Attribute("photo_tags_delay_seconds", UInt, func() {
		Description("задержка между загрузкой фотографии и проставлением отметок (в секундах)")
		Meta("struct:tag:json", "photo_tags_delay_seconds")
	})

	Attribute("posts_per_bot", UInt, func() {
		Description("количество постов для каждого бота")
		Meta("struct:tag:json", "posts_per_bot")
	})

	Attribute("photo_tags_posts_per_bot", UInt, func() {
		Description("количество постов с отметками на фото для каждого бота")
		Meta("struct:tag:json", "photo_tags_posts_per_bot")
	})

	Attribute("targets_per_post", UInt, func() {
		Description("количество упоминаний под каждым постом")
		Meta("struct:tag:json", "targets_per_post")
	})

	Attribute("photo_targets_per_post", UInt, func() {
		Description("количество упоминаний на фото у каждого поста")
		Meta("struct:tag:json", "photo_targets_per_post")
	})

	Attribute("post_images", ArrayOf(String), "список base64 строк картинок", func() {
		Meta("struct:tag:json", "post_images")
	})

	Attribute("bot_images", ArrayOf(String), func() {
		Description("аватарки для ботов")
		Meta("struct:tag:json", "bot_images")
	})

	Attribute("testing_tag_username", String, func() {
		Description("username пользователя в Instagram, без @. Фиксированная отметка для каждого поста, чтобы проверить работу ботов")
		Meta("struct:tag:json", "testing_tag_username")
	})

	Attribute("testing_tag_user_id", Int64, func() {
		Description("username пользователя в Instagram, без @. Фиксированная отметка для каждого поста, чтобы проверить работу ботов")
		Meta("struct:tag:json", "testing_tag_user_id")
	})

	Required("id", "type", "text_template", "post_images", "status", "title", "bots_num", "residential_proxies_num",
		"cheap_proxies_num", "targets_num", "bot_images", "landing_accounts", "bot_names", "bot_last_names", "bot_urls",
		"targets_per_post", "posts_per_bot", "photo_tags_delay_seconds", "per_post_sleep_seconds", "need_photo_tags", "follow_targets",
		"photo_targets_per_post", "photo_tags_posts_per_bot",
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
	Attribute("username", String, "имя пользователя бота", func() {
	})
	Attribute("posts_count", Int32, "количество выложенных постов", func() {
		Meta("struct:tag:json", "posts_count")
	})
	Attribute("status", Int32, "текущий статус бота, будут ли выкладываться посты")
	Attribute("description_targets_notified", Int32, "количество аккаунтов, которых упомянули в постах", func() {
		Meta("struct:tag:json", "description_targets_notified")
	})
	Attribute("photo_targets_notified", Int32, "количество аккаунтов, которых упомянули в постах на фото", func() {
		Meta("struct:tag:json", "photo_targets_notified")
	})

	Attribute("file_order", Int32, "номер бота в загруженном файле", func() {
		Meta("struct:tag:json", "file_order")
	})

	Required("username", "posts_count", "status", "description_targets_notified", "photo_targets_notified", "file_order")
})

var TaskProgress = Type("TaskProgress", func() {
	Attribute("bots_progresses", ArrayOf(BotsProgress), func() {
		Description("результат работы по каждому боту")
		Meta("struct:tag:json", "bots_progresses")
	})

	Attribute("targets_notified", Int, "количество аккаунтов, которых упомянули в постах", func() {
		Meta("struct:tag:json", "targets_notified")
	})
	Attribute("photo_targets_notified", Int, "количество аккаунтов, которых упомянули в постах на фото", func() {
		Meta("struct:tag:json", "photo_targets_notified")
	})
	Attribute("targets_failed", Int, "количество аккаунтов, которых не получилось упомянуть, при перезапуске задачи будут использованы заново", func() {
		Meta("struct:tag:json", "targets_failed")
	})

	Attribute("targets_waiting", Int, "количество аккаунтов, которых не выбрали для постов", func() {
		Meta("struct:tag:json", "targets_waiting")
	})

	Attribute("done", Boolean, "закончена ли задача")

	Attribute("bots_total", Int, "общее количество ботов в задаче", func() {
		Meta("struct:tag:json", "bots_total")
	})

	Required("bots_progresses", "targets_notified", "photo_targets_notified", "targets_failed",
		"targets_waiting", "done", "bots_total")
})
