package main

//
// type APIResponse map[string]interface{}
//
// type Status string
//
// const (
// 	UnknownStatus = ""
// 	OKStatus      = "ok"
// 	FailedStatus  = "fail"
// )
//
// func (r APIResponse) GetStatus() Status {
// 	val, ok := r["status"]
// 	if !ok {
// 		return UnknownStatus
// 	}
//
// 	st, ok := val.(string)
// 	if !ok {
// 		return UnknownStatus
// 	}
//
// 	return Status(st)
// }
//
// func (r APIResponse) GetMessage() string {
//
// }
