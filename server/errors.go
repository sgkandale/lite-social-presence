package server

var (
	Err_ReadingRequest            = GeneralResponse{Message: "error reading request"}
	Err_UserAlreadyRegistered     = GeneralResponse{Message: "user already registered"}
	Err_UserNotFound              = GeneralResponse{Message: "user not found"}
	Err_UserIdMissing             = GeneralResponse{Message: "user_id is missing"}
	Err_FriendshipNotFound        = GeneralResponse{Message: "friendship not found"}
	Err_FriendshipAlreadyExists   = GeneralResponse{Message: "friendship already exists"}
	Err_RequestIdMissing          = GeneralResponse{Message: "request_id is missing"}
	Err_RequestIdInvalid          = GeneralResponse{Message: "request_id is invalid"}
	Err_FriendshipRequestNotFound = GeneralResponse{Message: "friendship request not found"}
	Err_SomethingWrong            = GeneralResponse{Message: "something went wrong"}
	Err_AuthHeaderMissing         = GeneralResponse{Message: "'Authorization' header is missing"}
)

var (
	Resp_Success = GeneralResponse{Message: "success"}
)
