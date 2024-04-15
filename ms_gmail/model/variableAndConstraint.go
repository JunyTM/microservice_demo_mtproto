package model

// Lib 'binary' is use int32 type for binary.Write
var MapMessageObject = map[string]interface{}{
	"Login":  LoginPayload{},
	"Regist": RegistPayload{},
}

var MapMessageObjectId = map[string]int32{
	"Login":  1,
	"Regist": 2,
}

// This is a map of message object from Client to Server
var MapObjectDecode = map[int32]interface{}{
	1:   LoginPayload{},
	2:   RegistPayload{},
	3:   AuthResponse{},
}
