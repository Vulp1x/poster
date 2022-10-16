// nolint
package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("rest-api", func() {
	Title("REST api for simple route app")
})

// JWTAuth defines a security scheme that uses JWT tokens.
var JWTAuth = JWTSecurity("jwt", func() {
	Description(`Secures endpoint by requiring a valid JWT token retrieved via the signin endpoint. Supports scopes "api:read" and "api:write".`)
	Scope("driver", "Read-only access")
	Scope("admin", "Read and write access")
})

// BasicAuth defines a security scheme using basic authentication. The scheme
// protects the "signin" action used to create JWTs.
var BasicAuth = BasicAuthSecurity("basic", func() {
	Description("Basic authentication used to authenticate security principal during signin")
	Scope("driver", "Read-only access")
})

var _ = Service("auth_service", func() {
	Description("The secured service exposes endpoints that require valid authorization credentials.")

	Error("unauthorized", String, "Credentials are invalid")
	Error("bad request", String, "Invalid request")
	Error("internal error", String, "internal error")
	Error("user not found", String, "Not found")

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
		Response("bad request", StatusBadRequest)
		Response("user not found", StatusNotFound)
		Response("internal error", StatusInternalServerError)
	})

	Method("signin", func() {
		Description("Creates a valid JWT")

		// The signin endpoint is secured via basic auth
		Security(BasicAuth)

		Payload(func() {

			Description("Credentials used to authenticate to retrieve JWT token")
			UsernameField(1, "login", String, "login used to perform signin", func() {
				Example("user@test.ru")
			})
			PasswordField(2, "password", String, "Password used to perform signin", func() {
				Example("password")
			})
			Required("login", "password")
		})

		Result(Creds)

		HTTP(func() {
			POST("/api/signin")
			// Use Authorization header to provide basic auth value.
			Response(StatusOK)
		})
	})

	Method("profile", func() {
		Description("get user profile")

		Security(JWTAuth)

		Payload(func() {
			TokenField(1, "token", String, func() {
				Description("JWT used for authentication")
			})
			Required("token")
		})

		HTTP(func() {
			GET("/api/profile")
			// Use Authorization header to provide basic auth value.
			Response(StatusOK)
			Response(StatusUnauthorized)
		})
	})
})

var _ = Service("tasks_service", func() {
	Description("сервис для создания, редактирования и работы с задачами (рекламными компаниями)")

	Error("unauthorized", String, "Credentials are invalid")
	Error("bad request", String, "Invalid request")
	Error("internal error", String, "internal error")
	Error("task not found", String, "Not found")

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
		Response("bad request", StatusBadRequest)
		Response("task not found", StatusNotFound)
		Response("internal error", StatusInternalServerError)
	})

	Method("create task draft", func() {
		Description("создать драфт задачи")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("title", String, func() {
				Description("название задачи")
			})

			Attribute("text_template", String, func() {
				Description("шаблон для подписи под постом")
				Meta("struct:tag:json", "text_template")
			})

			Attribute("post_images", ArrayOf(String), func() {
				Description("список фотографий для постов")
				Meta("struct:tag:json", "post_images")
			})

			Required("token", "title", "text_template", "post_images")
		})

		Result(String, func() {
			Description("task_id для созданной задачи")
			Format(FormatUUID)
		})

		HTTP(func() {
			POST("/api/tasks/draft")
			// Use Authorization header to provide basic auth value.
			Response(StatusOK)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("update task", func() {
		Description(`обновить информацию о задаче. Не меняет статус задачи, можно вызывать сколько угодно раз.
			Нельзя вызвать для задачи, которая уже выполняется, для этого надо сначала остановить выполнение.`)

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("task_id", String, func() {
				Description("id задачи, которую хотим обновить")
				Meta("struct:tag:json", "task_id")
			})

			Attribute("title", String, func() {
				Description("название задачи")
			})

			Attribute("text_template", String, func() {
				Description("шаблон для подписи под постом")
				Meta("struct:tag:json", "text_template")
			})

			Attribute("post_images", ArrayOf(String), func() {
				Description("фотография для постов")
				Meta("struct:tag:json", "post_images")
			})

			Required("token", "task_id")
		})

		Result(Task)

		HTTP(func() {
			PUT("/api/tasks/{task_id}/")
			// Use Authorization header to provide basic auth value.
			Response(StatusOK)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("upload files", func() {
		Description("загрузить файл с пользователями, прокси")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("task_id", String, func() {
				Description("id задачи, в которую загружаем пользователей/прокси")
				Meta("struct:tag:json", "task_id")
			})

			Attribute("filenames", TaskFilenames)

			Attribute("bots", ArrayOf(BotAccountRecord), "список ботов")
			Attribute("residential_proxies", ArrayOf(ProxyRecord), "список проксей для использования")
			Attribute("cheap_proxies", ArrayOf(ProxyRecord), "список дешёвых проксей для загрузки фото")
			Attribute("targets", ArrayOf(TargetUserRecord), "список аккаунтов, которым показать надо рекламу")

			Required("token", "task_id", "bots", "residential_proxies", "cheap_proxies", "targets", "filenames")
		})

		Result(func() {
			Attribute("upload_errors", ArrayOf(UploadError), func() {
				Description("ошибки, которые возникли при загрузке файлов")
				Meta("struct:tag:json", "upload_errors")
			})

			Attribute("status", TaskStatus)

			Required("status", "upload_errors")
		})

		HTTP(func() {
			POST("/api/tasks/{task_id}/upload")
			MultipartRequest()
			// Use Authorization header to provide basic auth value.
			Response(StatusOK)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("assign proxies", func() {
		Description("присвоить ботам прокси")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("task_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "task_id")
			})

			Required("token", "task_id")
		})

		Result(func() {
			Attribute("bots_number", Int, func() {
				Description("количество аккаунтов с проксями, которые будут использованы для текущей задачи")
				Meta("struct:tag:json", "bots_number")
			})

			Attribute("status", TaskStatus)
			Attribute("task_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "task_id")
			})

			Required("bots_number", "task_id", "status")
		})

		HTTP(func() {
			POST("/api/tasks/{task_id}/assign")
			Response(StatusOK)
			Response(StatusBadRequest)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("force delete", func() {
		Description("удалить задачу и все связанные с ней сущности. Использовать только для тестов")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("task_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "task_id")
			})

			Required("token", "task_id")
		})

		HTTP(func() {
			DELETE("/api/tasks/{task_id}/force")
			Response(StatusOK)
			Response(StatusBadRequest)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
			Response(StatusInternalServerError)

		})
	})

	Method("start task", func() {
		Description("начать выполнение задачи ")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("task_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "task_id")
			})

			Required("token", "task_id")
		})

		Result(func() {
			Attribute("status", TaskStatus)
			Attribute("task_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "task_id")
			})

			Required("task_id", "status")
		})

		HTTP(func() {
			POST("/api/tasks/{task_id}/start")
			Response(StatusOK)
			Response(StatusBadRequest)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("stop task", func() {
		Description("остановить выполнение задачи ")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("task_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "task_id")
			})

			Required("token", "task_id")
		})

		Result(func() {
			Attribute("status", TaskStatus)
			Attribute("task_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "task_id")
			})

			Required("task_id", "status")
		})

		HTTP(func() {
			POST("/api/tasks/{task_id}/stop")
			Response(StatusOK)
			Response(StatusBadRequest)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("get task", func() {
		Description("получить задачу по id")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("task_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "task_id")
			})

			Required("token", "task_id")
		})

		Result(Task)

		HTTP(func() {
			GET("/api/tasks/{task_id}/")
			Response(StatusOK)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("get progress", func() {
		Description("получить статус выполнения задачи по id")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Attribute("task_id", String, func() {
				Description("id задачи")
				Meta("struct:tag:json", "task_id")
			})

			Required("token", "task_id")
		})

		Result(ArrayOf(BotsProgress))

		HTTP(func() {
			GET("/api/tasks/{task_id}/progress")
			Response(StatusOK)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
		})
	})

	Method("list tasks", func() {
		Description("получить все задачи для текущего пользователя")

		Security(JWTAuth)

		Payload(func() {
			Token("token", String, func() {
				Description("JWT used for authentication")
			})

			Required("token")
		})

		Result(ArrayOf(Task))

		HTTP(func() {
			GET("/api/tasks/")
			Response(StatusOK)
			Response(StatusNotFound)
			Response(StatusUnauthorized)
			Response(StatusInternalServerError)
		})
	})
})

var _ = Service("admin_service", func() {
	Description("The secured service exposes endpoints that require valid authorization credentials.")

	Error("unauthorized", String, "Credentials are invalid")
	Error("user not found", String)
	Error("bad request", String, "Invalid request")
	Error("internal error", String, "internal error")

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
		Response("user not found", StatusNotFound)
		Response("bad request", StatusBadRequest)
		Response("bad request", StatusInternalServerError)
	})

	Method("add_manager", func() {
		Description("admins could add drivers from main system")

		Security(JWTAuth, func() { // Use JWT to auth requests to this endpoint.
			Scope("admin") // Enforce presence of "api:write" scope in JWT claims.
		})

		Payload(func() {
			TokenField(1, "token", String, func() {
				Description("JWT used for authentication")
			})
			Field(2, "login", String)
			Field(3, "password", String, func() {
				MinLength(4)
			})

			Required("login", "password")
		})

		Error("invalid-scopes", String, "Token scopes are invalid")

		HTTP(func() {
			POST("api/admin/driver")
			Response(StatusOK)
			Response("invalid-scopes", StatusForbidden)
		})
	})

	Method("drop_manager", func() {
		Description("admins could delete managers from main system")

		Security(JWTAuth, func() { // Use JWT to auth requests to this endpoint.
			Scope("admin") // Enforce presence of "api:write" scope in JWT claims.
		})

		Payload(func() {
			TokenField(1, "token", String, func() {
				Description("JWT used for authentication")
			})

			Field(2, "manager_id", String, func() {
				Description("id менеджера, которого необходимо удалить")
				Format(FormatUUID)
				Meta("struct:tag:json", "manager_id")
			})

			Required("manager_id")
		})

		Error("invalid-scopes", String, "Token scopes are invalid")

		HTTP(func() {
			DELETE("api/admin/driver/{manager_id}")
			Response(StatusOK)
			Response(StatusNotFound)
			Response("invalid-scopes", StatusForbidden)
			Response(StatusInternalServerError)
		})
	})

})
