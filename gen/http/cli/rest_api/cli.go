// Code generated by goa v3.8.5, DO NOT EDIT.
//
// rest-api HTTP client CLI support package
//
// Command:
// $ goa gen github.com/inst-api/poster/design

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	adminservicec "github.com/inst-api/poster/gen/http/admin_service/client"
	authservicec "github.com/inst-api/poster/gen/http/auth_service/client"
	tasksservicec "github.com/inst-api/poster/gen/http/tasks_service/client"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//	command (subcommand1|subcommand2|...)
func UsageCommands() string {
	return `admin-service (add-manager|drop-manager)
auth-service (signin|profile)
tasks-service (create-task-draft|upload-file|start-task|stop-task|get-task|list-tasks)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` admin-service add-manager --body '{
      "login": "Quia voluptatem impedit mollitia.",
      "password": "y3x"
   }' --token "Perferendis odit."` + "\n" +
		os.Args[0] + ` auth-service signin --login "user@test.ru" --password "password"` + "\n" +
		os.Args[0] + ` tasks-service create-task-draft --body '{
      "post_image": "Voluptatibus autem.",
      "text_template": "Dolor nesciunt nisi aperiam deserunt aut.",
      "title": "Excepturi eos aspernatur."
   }' --token "Aperiam quo mollitia necessitatibus voluptatibus porro voluptatum."` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	scheme, host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restore bool,
	tasksServiceUploadFileEncoderFn tasksservicec.TasksServiceUploadFileEncoderFunc,
) (goa.Endpoint, interface{}, error) {
	var (
		adminServiceFlags = flag.NewFlagSet("admin-service", flag.ContinueOnError)

		adminServiceAddManagerFlags     = flag.NewFlagSet("add-manager", flag.ExitOnError)
		adminServiceAddManagerBodyFlag  = adminServiceAddManagerFlags.String("body", "REQUIRED", "")
		adminServiceAddManagerTokenFlag = adminServiceAddManagerFlags.String("token", "", "")

		adminServiceDropManagerFlags         = flag.NewFlagSet("drop-manager", flag.ExitOnError)
		adminServiceDropManagerManagerIDFlag = adminServiceDropManagerFlags.String("manager-id", "REQUIRED", "id менеджера, которого необходимо удалить")
		adminServiceDropManagerTokenFlag     = adminServiceDropManagerFlags.String("token", "", "")

		authServiceFlags = flag.NewFlagSet("auth-service", flag.ContinueOnError)

		authServiceSigninFlags        = flag.NewFlagSet("signin", flag.ExitOnError)
		authServiceSigninLoginFlag    = authServiceSigninFlags.String("login", "REQUIRED", "login used to perform signin")
		authServiceSigninPasswordFlag = authServiceSigninFlags.String("password", "REQUIRED", "Password used to perform signin")

		authServiceProfileFlags     = flag.NewFlagSet("profile", flag.ExitOnError)
		authServiceProfileTokenFlag = authServiceProfileFlags.String("token", "REQUIRED", "")

		tasksServiceFlags = flag.NewFlagSet("tasks-service", flag.ContinueOnError)

		tasksServiceCreateTaskDraftFlags     = flag.NewFlagSet("create-task-draft", flag.ExitOnError)
		tasksServiceCreateTaskDraftBodyFlag  = tasksServiceCreateTaskDraftFlags.String("body", "REQUIRED", "")
		tasksServiceCreateTaskDraftTokenFlag = tasksServiceCreateTaskDraftFlags.String("token", "REQUIRED", "")

		tasksServiceUploadFileFlags      = flag.NewFlagSet("upload-file", flag.ExitOnError)
		tasksServiceUploadFileBodyFlag   = tasksServiceUploadFileFlags.String("body", "REQUIRED", "")
		tasksServiceUploadFileTaskIDFlag = tasksServiceUploadFileFlags.String("task-id", "REQUIRED", "id задачи, в которую загружаем пользователей/прокси")
		tasksServiceUploadFileTokenFlag  = tasksServiceUploadFileFlags.String("token", "REQUIRED", "")

		tasksServiceStartTaskFlags      = flag.NewFlagSet("start-task", flag.ExitOnError)
		tasksServiceStartTaskTaskIDFlag = tasksServiceStartTaskFlags.String("task-id", "REQUIRED", "id задачи")
		tasksServiceStartTaskTokenFlag  = tasksServiceStartTaskFlags.String("token", "REQUIRED", "")

		tasksServiceStopTaskFlags      = flag.NewFlagSet("stop-task", flag.ExitOnError)
		tasksServiceStopTaskTaskIDFlag = tasksServiceStopTaskFlags.String("task-id", "REQUIRED", "id задачи")
		tasksServiceStopTaskTokenFlag  = tasksServiceStopTaskFlags.String("token", "REQUIRED", "")

		tasksServiceGetTaskFlags      = flag.NewFlagSet("get-task", flag.ExitOnError)
		tasksServiceGetTaskTaskIDFlag = tasksServiceGetTaskFlags.String("task-id", "REQUIRED", "id задачи")
		tasksServiceGetTaskTokenFlag  = tasksServiceGetTaskFlags.String("token", "REQUIRED", "")

		tasksServiceListTasksFlags     = flag.NewFlagSet("list-tasks", flag.ExitOnError)
		tasksServiceListTasksTokenFlag = tasksServiceListTasksFlags.String("token", "REQUIRED", "")
	)
	adminServiceFlags.Usage = adminServiceUsage
	adminServiceAddManagerFlags.Usage = adminServiceAddManagerUsage
	adminServiceDropManagerFlags.Usage = adminServiceDropManagerUsage

	authServiceFlags.Usage = authServiceUsage
	authServiceSigninFlags.Usage = authServiceSigninUsage
	authServiceProfileFlags.Usage = authServiceProfileUsage

	tasksServiceFlags.Usage = tasksServiceUsage
	tasksServiceCreateTaskDraftFlags.Usage = tasksServiceCreateTaskDraftUsage
	tasksServiceUploadFileFlags.Usage = tasksServiceUploadFileUsage
	tasksServiceStartTaskFlags.Usage = tasksServiceStartTaskUsage
	tasksServiceStopTaskFlags.Usage = tasksServiceStopTaskUsage
	tasksServiceGetTaskFlags.Usage = tasksServiceGetTaskUsage
	tasksServiceListTasksFlags.Usage = tasksServiceListTasksUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "admin-service":
			svcf = adminServiceFlags
		case "auth-service":
			svcf = authServiceFlags
		case "tasks-service":
			svcf = tasksServiceFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "admin-service":
			switch epn {
			case "add-manager":
				epf = adminServiceAddManagerFlags

			case "drop-manager":
				epf = adminServiceDropManagerFlags

			}

		case "auth-service":
			switch epn {
			case "signin":
				epf = authServiceSigninFlags

			case "profile":
				epf = authServiceProfileFlags

			}

		case "tasks-service":
			switch epn {
			case "create-task-draft":
				epf = tasksServiceCreateTaskDraftFlags

			case "upload-file":
				epf = tasksServiceUploadFileFlags

			case "start-task":
				epf = tasksServiceStartTaskFlags

			case "stop-task":
				epf = tasksServiceStopTaskFlags

			case "get-task":
				epf = tasksServiceGetTaskFlags

			case "list-tasks":
				epf = tasksServiceListTasksFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     interface{}
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "admin-service":
			c := adminservicec.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "add-manager":
				endpoint = c.AddManager()
				data, err = adminservicec.BuildAddManagerPayload(*adminServiceAddManagerBodyFlag, *adminServiceAddManagerTokenFlag)
			case "drop-manager":
				endpoint = c.DropManager()
				data, err = adminservicec.BuildDropManagerPayload(*adminServiceDropManagerManagerIDFlag, *adminServiceDropManagerTokenFlag)
			}
		case "auth-service":
			c := authservicec.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "signin":
				endpoint = c.Signin()
				data, err = authservicec.BuildSigninPayload(*authServiceSigninLoginFlag, *authServiceSigninPasswordFlag)
			case "profile":
				endpoint = c.Profile()
				data, err = authservicec.BuildProfilePayload(*authServiceProfileTokenFlag)
			}
		case "tasks-service":
			c := tasksservicec.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "create-task-draft":
				endpoint = c.CreateTaskDraft()
				data, err = tasksservicec.BuildCreateTaskDraftPayload(*tasksServiceCreateTaskDraftBodyFlag, *tasksServiceCreateTaskDraftTokenFlag)
			case "upload-file":
				endpoint = c.UploadFile(tasksServiceUploadFileEncoderFn)
				data, err = tasksservicec.BuildUploadFilePayload(*tasksServiceUploadFileBodyFlag, *tasksServiceUploadFileTaskIDFlag, *tasksServiceUploadFileTokenFlag)
			case "start-task":
				endpoint = c.StartTask()
				data, err = tasksservicec.BuildStartTaskPayload(*tasksServiceStartTaskTaskIDFlag, *tasksServiceStartTaskTokenFlag)
			case "stop-task":
				endpoint = c.StopTask()
				data, err = tasksservicec.BuildStopTaskPayload(*tasksServiceStopTaskTaskIDFlag, *tasksServiceStopTaskTokenFlag)
			case "get-task":
				endpoint = c.GetTask()
				data, err = tasksservicec.BuildGetTaskPayload(*tasksServiceGetTaskTaskIDFlag, *tasksServiceGetTaskTokenFlag)
			case "list-tasks":
				endpoint = c.ListTasks()
				data, err = tasksservicec.BuildListTasksPayload(*tasksServiceListTasksTokenFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// admin-serviceUsage displays the usage of the admin-service command and its
// subcommands.
func adminServiceUsage() {
	fmt.Fprintf(os.Stderr, `The secured service exposes endpoints that require valid authorization credentials.
Usage:
    %[1]s [globalflags] admin-service COMMAND [flags]

COMMAND:
    add-manager: admins could add drivers from main system
    drop-manager: admins could delete managers from main system

Additional help:
    %[1]s admin-service COMMAND --help
`, os.Args[0])
}
func adminServiceAddManagerUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] admin-service add-manager -body JSON -token STRING

admins could add drivers from main system
    -body JSON: 
    -token STRING: 

Example:
    %[1]s admin-service add-manager --body '{
      "login": "Quia voluptatem impedit mollitia.",
      "password": "y3x"
   }' --token "Perferendis odit."
`, os.Args[0])
}

func adminServiceDropManagerUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] admin-service drop-manager -manager-id STRING -token STRING

admins could delete managers from main system
    -manager-id STRING: id менеджера, которого необходимо удалить
    -token STRING: 

Example:
    %[1]s admin-service drop-manager --manager-id "FB1E8AC6-4FA4-C883-ED5A-54960E88F5FE" --token "Omnis odit eligendi temporibus in."
`, os.Args[0])
}

// auth-serviceUsage displays the usage of the auth-service command and its
// subcommands.
func authServiceUsage() {
	fmt.Fprintf(os.Stderr, `The secured service exposes endpoints that require valid authorization credentials.
Usage:
    %[1]s [globalflags] auth-service COMMAND [flags]

COMMAND:
    signin: Creates a valid JWT
    profile: get user profile

Additional help:
    %[1]s auth-service COMMAND --help
`, os.Args[0])
}
func authServiceSigninUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] auth-service signin -login STRING -password STRING

Creates a valid JWT
    -login STRING: login used to perform signin
    -password STRING: Password used to perform signin

Example:
    %[1]s auth-service signin --login "user@test.ru" --password "password"
`, os.Args[0])
}

func authServiceProfileUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] auth-service profile -token STRING

get user profile
    -token STRING: 

Example:
    %[1]s auth-service profile --token "Rerum ea est facere itaque repellendus."
`, os.Args[0])
}

// tasks-serviceUsage displays the usage of the tasks-service command and its
// subcommands.
func tasksServiceUsage() {
	fmt.Fprintf(os.Stderr, `сервис для создания, редактирования и работы с задачами (рекламными компаниями)
Usage:
    %[1]s [globalflags] tasks-service COMMAND [flags]

COMMAND:
    create-task-draft: создать драфт задачи
    upload-file: загрузить файл с пользователями, прокси
    start-task: начать выполнение задачи 
    stop-task: остановить выполнение задачи 
    get-task: получить задачу по id
    list-tasks: получить все задачи для текущего пользователя

Additional help:
    %[1]s tasks-service COMMAND --help
`, os.Args[0])
}
func tasksServiceCreateTaskDraftUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] tasks-service create-task-draft -body JSON -token STRING

создать драфт задачи
    -body JSON: 
    -token STRING: 

Example:
    %[1]s tasks-service create-task-draft --body '{
      "post_image": "Voluptatibus autem.",
      "text_template": "Dolor nesciunt nisi aperiam deserunt aut.",
      "title": "Excepturi eos aspernatur."
   }' --token "Aperiam quo mollitia necessitatibus voluptatibus porro voluptatum."
`, os.Args[0])
}

func tasksServiceUploadFileUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] tasks-service upload-file -body JSON -task-id STRING -token STRING

загрузить файл с пользователями, прокси
    -body JSON: 
    -task-id STRING: id задачи, в которую загружаем пользователей/прокси
    -token STRING: 

Example:
    %[1]s tasks-service upload-file --body '{
      "bots": [
         {
            "advertising_id": "7D80331A-7620-D09D-7CCB-2EF87B797732",
            "device_id": "Tenetur sed ea excepturi delectus quasi.",
            "family_device_id": "76FB876C-96AC-91E7-BD21-B0C2988DDF65",
            "headers": {
               "Ducimus voluptate hic dicta impedit.": "Totam labore amet ut iure praesentium."
            },
            "password": "Quidem quis earum maxime omnis reiciendis adipisci.",
            "phone_id": "1265498D-5A84-134A-1C7A-ED5B4B92788E",
            "user_agent": "Perspiciatis autem quo.",
            "username": "Autem est quia.",
            "uuid": "5E3B665E-1239-9C12-9643-FFC1E6C04697"
         },
         {
            "advertising_id": "7D80331A-7620-D09D-7CCB-2EF87B797732",
            "device_id": "Tenetur sed ea excepturi delectus quasi.",
            "family_device_id": "76FB876C-96AC-91E7-BD21-B0C2988DDF65",
            "headers": {
               "Ducimus voluptate hic dicta impedit.": "Totam labore amet ut iure praesentium."
            },
            "password": "Quidem quis earum maxime omnis reiciendis adipisci.",
            "phone_id": "1265498D-5A84-134A-1C7A-ED5B4B92788E",
            "user_agent": "Perspiciatis autem quo.",
            "username": "Autem est quia.",
            "uuid": "5E3B665E-1239-9C12-9643-FFC1E6C04697"
         },
         {
            "advertising_id": "7D80331A-7620-D09D-7CCB-2EF87B797732",
            "device_id": "Tenetur sed ea excepturi delectus quasi.",
            "family_device_id": "76FB876C-96AC-91E7-BD21-B0C2988DDF65",
            "headers": {
               "Ducimus voluptate hic dicta impedit.": "Totam labore amet ut iure praesentium."
            },
            "password": "Quidem quis earum maxime omnis reiciendis adipisci.",
            "phone_id": "1265498D-5A84-134A-1C7A-ED5B4B92788E",
            "user_agent": "Perspiciatis autem quo.",
            "username": "Autem est quia.",
            "uuid": "5E3B665E-1239-9C12-9643-FFC1E6C04697"
         }
      ],
      "proxies": [
         {
            "host": "Eligendi accusantium enim.",
            "login": "In et assumenda voluptate deleniti ut aut.",
            "password": "Quod tenetur.",
            "port": 4440880418486248045
         },
         {
            "host": "Eligendi accusantium enim.",
            "login": "In et assumenda voluptate deleniti ut aut.",
            "password": "Quod tenetur.",
            "port": 4440880418486248045
         }
      ],
      "targets": [
         {
            "user_id": 372918204077540257,
            "username": "Labore impedit repellat."
         },
         {
            "user_id": 372918204077540257,
            "username": "Labore impedit repellat."
         }
      ]
   }' --task-id "Dolor rerum." --token "Velit optio quod rerum aut blanditiis."
`, os.Args[0])
}

func tasksServiceStartTaskUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] tasks-service start-task -task-id STRING -token STRING

начать выполнение задачи 
    -task-id STRING: id задачи
    -token STRING: 

Example:
    %[1]s tasks-service start-task --task-id "Maiores dolorem quia." --token "Maxime et."
`, os.Args[0])
}

func tasksServiceStopTaskUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] tasks-service stop-task -task-id STRING -token STRING

остановить выполнение задачи 
    -task-id STRING: id задачи
    -token STRING: 

Example:
    %[1]s tasks-service stop-task --task-id "Perferendis quia rem nostrum sint dolorum." --token "Blanditiis ex alias."
`, os.Args[0])
}

func tasksServiceGetTaskUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] tasks-service get-task -task-id STRING -token STRING

получить задачу по id
    -task-id STRING: id задачи
    -token STRING: 

Example:
    %[1]s tasks-service get-task --task-id "Cupiditate eaque dolorem quos." --token "Facilis fugiat similique ab sint voluptatum aspernatur."
`, os.Args[0])
}

func tasksServiceListTasksUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] tasks-service list-tasks -token STRING

получить все задачи для текущего пользователя
    -token STRING: 

Example:
    %[1]s tasks-service list-tasks --token "Veniam ipsum iusto necessitatibus."
`, os.Args[0])
}
