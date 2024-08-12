package server

var (
	Err_ReadingRequest        = GeneralResponse{Message: "error reading request"}
	Err_UserAlreadyRegistered = GeneralResponse{Message: "user already registered"}
	Err_UserNotFound          = GeneralResponse{Message: "user not found"}
	Err_SomethingWrong        = GeneralResponse{Message: "something went wrong"}
)

var (
	Resp_Success = GeneralResponse{Message: "success"}
)
