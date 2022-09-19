
# GoAdmin Instruction

GoAdmin is a golang framework help gopher quickly build a data visualization platform. 

- [github](https://github.com/GoAdminGroup/go-admin)
- [forum](http://discuss.go-admin.com)
- [document](https://book.go-admin.cn)

## Directories Introduction

```
.
├── Makefile            Makefile
├── adm.ini             adm config
├── build               binary build target folder
├── config.yml          config file
├── go.mod              go.mod
├── go.sum              go.sum
├── html                frontend html files
├── logs                logs
├── main.go             program entrance file
├── main_test.go        test file
├── pages               page controllers
├── tables              table models
└── uploads             upload directory
```

## Generate Table Model

### online tool

visit: http://127.0.0.1:8090/admin/info/generate/new

### use adm

```
adm generate
```

