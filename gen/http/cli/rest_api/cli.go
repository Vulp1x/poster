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
	return `tasks-service (create-task|upload-file|start-task|get-task|list-tasks)
admin-service (add-manager|drop-manager)
auth-service (signin|profile)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` tasks-service create-task --body '{
      "description": "Ipsam expedita libero eum et.",
      "tittle": "Dolor culpa temporibus sit."
   }' --token "Ut quibusdam sequi sed unde sed."` + "\n" +
		os.Args[0] + ` admin-service add-manager --body '{
      "login": "Illo aut non sint alias.",
      "password": "roc"
   }' --token "Rerum nam optio animi magnam."` + "\n" +
		os.Args[0] + ` auth-service signin --login "user@test.ru" --password "password"` + "\n" +
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
		tasksServiceFlags = flag.NewFlagSet("tasks-service", flag.ContinueOnError)

		tasksServiceCreateTaskFlags     = flag.NewFlagSet("create-task", flag.ExitOnError)
		tasksServiceCreateTaskBodyFlag  = tasksServiceCreateTaskFlags.String("body", "REQUIRED", "")
		tasksServiceCreateTaskTokenFlag = tasksServiceCreateTaskFlags.String("token", "REQUIRED", "")

		tasksServiceUploadFileFlags      = flag.NewFlagSet("upload-file", flag.ExitOnError)
		tasksServiceUploadFileBodyFlag   = tasksServiceUploadFileFlags.String("body", "REQUIRED", "")
		tasksServiceUploadFileTaskIDFlag = tasksServiceUploadFileFlags.String("task-id", "REQUIRED", "id задачи, в которую загружаем пользователей/прокси")
		tasksServiceUploadFileTokenFlag  = tasksServiceUploadFileFlags.String("token", "REQUIRED", "")

		tasksServiceStartTaskFlags      = flag.NewFlagSet("start-task", flag.ExitOnError)
		tasksServiceStartTaskTaskIDFlag = tasksServiceStartTaskFlags.String("task-id", "REQUIRED", "id задачи")
		tasksServiceStartTaskTokenFlag  = tasksServiceStartTaskFlags.String("token", "REQUIRED", "")

		tasksServiceGetTaskFlags      = flag.NewFlagSet("get-task", flag.ExitOnError)
		tasksServiceGetTaskTaskIDFlag = tasksServiceGetTaskFlags.String("task-id", "REQUIRED", "id задачи")
		tasksServiceGetTaskTokenFlag  = tasksServiceGetTaskFlags.String("token", "REQUIRED", "")

		tasksServiceListTasksFlags     = flag.NewFlagSet("list-tasks", flag.ExitOnError)
		tasksServiceListTasksTokenFlag = tasksServiceListTasksFlags.String("token", "REQUIRED", "")

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
	)
	tasksServiceFlags.Usage = tasksServiceUsage
	tasksServiceCreateTaskFlags.Usage = tasksServiceCreateTaskUsage
	tasksServiceUploadFileFlags.Usage = tasksServiceUploadFileUsage
	tasksServiceStartTaskFlags.Usage = tasksServiceStartTaskUsage
	tasksServiceGetTaskFlags.Usage = tasksServiceGetTaskUsage
	tasksServiceListTasksFlags.Usage = tasksServiceListTasksUsage

	adminServiceFlags.Usage = adminServiceUsage
	adminServiceAddManagerFlags.Usage = adminServiceAddManagerUsage
	adminServiceDropManagerFlags.Usage = adminServiceDropManagerUsage

	authServiceFlags.Usage = authServiceUsage
	authServiceSigninFlags.Usage = authServiceSigninUsage
	authServiceProfileFlags.Usage = authServiceProfileUsage

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
		case "tasks-service":
			svcf = tasksServiceFlags
		case "admin-service":
			svcf = adminServiceFlags
		case "auth-service":
			svcf = authServiceFlags
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
		case "tasks-service":
			switch epn {
			case "create-task":
				epf = tasksServiceCreateTaskFlags

			case "upload-file":
				epf = tasksServiceUploadFileFlags

			case "start-task":
				epf = tasksServiceStartTaskFlags

			case "get-task":
				epf = tasksServiceGetTaskFlags

			case "list-tasks":
				epf = tasksServiceListTasksFlags

			}

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
		case "tasks-service":
			c := tasksservicec.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "create-task":
				endpoint = c.CreateTask()
				data, err = tasksservicec.BuildCreateTaskPayload(*tasksServiceCreateTaskBodyFlag, *tasksServiceCreateTaskTokenFlag)
			case "upload-file":
				endpoint = c.UploadFile(tasksServiceUploadFileEncoderFn)
				data, err = tasksservicec.BuildUploadFilePayload(*tasksServiceUploadFileBodyFlag, *tasksServiceUploadFileTaskIDFlag, *tasksServiceUploadFileTokenFlag)
			case "start-task":
				endpoint = c.StartTask()
				data, err = tasksservicec.BuildStartTaskPayload(*tasksServiceStartTaskTaskIDFlag, *tasksServiceStartTaskTokenFlag)
			case "get-task":
				endpoint = c.GetTask()
				data, err = tasksservicec.BuildGetTaskPayload(*tasksServiceGetTaskTaskIDFlag, *tasksServiceGetTaskTokenFlag)
			case "list-tasks":
				endpoint = c.ListTasks()
				data, err = tasksservicec.BuildListTasksPayload(*tasksServiceListTasksTokenFlag)
			}
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
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// tasks-serviceUsage displays the usage of the tasks-service command and its
// subcommands.
func tasksServiceUsage() {
	fmt.Fprintf(os.Stderr, `сервис для создания, редактирования и работы с задачами (рекламными компаниями)
Usage:
    %[1]s [globalflags] tasks-service COMMAND [flags]

COMMAND:
    create-task: создать драфт задачи
    upload-file: загрузить файл с пользователями, прокси
    start-task: начать выполнение задачи 
    get-task: получить задачу по id
    list-tasks: получить все задачи для текущего пользователя

Additional help:
    %[1]s tasks-service COMMAND --help
`, os.Args[0])
}
func tasksServiceCreateTaskUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] tasks-service create-task -body JSON -token STRING

создать драфт задачи
    -body JSON: 
    -token STRING: 

Example:
    %[1]s tasks-service create-task --body '{
      "description": "Ipsam expedita libero eum et.",
      "tittle": "Dolor culpa temporibus sit."
   }' --token "Ut quibusdam sequi sed unde sed."
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
      "bytes": "RXZlbmlldCBtb2xlc3RpYWUgc2ludCByZXJ1bSBldCBvZGl0Lg==",
      "file_type": 1
   }' --task-id "Saepe excepturi." --token "Ut expedita officiis."
`, os.Args[0])
}

func tasksServiceStartTaskUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] tasks-service start-task -task-id STRING -token STRING

начать выполнение задачи 
    -task-id STRING: id задачи
    -token STRING: 

Example:
    %[1]s tasks-service start-task --task-id "Est error." --token "Ut voluptas sit ut quo."
`, os.Args[0])
}

func tasksServiceGetTaskUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] tasks-service get-task -task-id STRING -token STRING

получить задачу по id
    -task-id STRING: id задачи
    -token STRING: 

Example:
    %[1]s tasks-service get-task --task-id "Optio molestiae eius consequuntur nemo nulla placeat." --token "At odio accusantium."
`, os.Args[0])
}

func tasksServiceListTasksUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] tasks-service list-tasks -token STRING

получить все задачи для текущего пользователя
    -token STRING: 

Example:
    %[1]s tasks-service list-tasks --token "Doloremque fugit aspernatur sed culpa."
`, os.Args[0])
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
      "login": "Illo aut non sint alias.",
      "password": "roc"
   }' --token "Rerum nam optio animi magnam."
`, os.Args[0])
}

func adminServiceDropManagerUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] admin-service drop-manager -manager-id STRING -token STRING

admins could delete managers from main system
    -manager-id STRING: id менеджера, которого необходимо удалить
    -token STRING: 

Example:
    %[1]s admin-service drop-manager --manager-id "77EB7E77-465C-FCC6-CEC6-11F6C8938D24" --token "Sed voluptatum dolores corrupti expedita."
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
    %[1]s auth-service profile --token "Ipsa non nobis commodi iure unde fugiat."
`, os.Args[0])
}
