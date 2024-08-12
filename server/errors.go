package server

var (
	Err_ReadingRequest        = GeneralResponse{Message: "error reading request"}
	Err_UserAlreadyRegistered = GeneralResponse{Message: "user already registered"}
	Err_UserNotFound          = GeneralResponse{Message: "user not found"}
	Err_SomethingWrong        = GeneralResponse{Message: "something went wrong"}
	Err_AuthHeaderMissing     = GeneralResponse{Message: "'Authorization' header is missing"}
)

var (
	Resp_Success = GeneralResponse{Message: "success"}
)
