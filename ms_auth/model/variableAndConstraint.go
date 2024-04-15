package model

// This is a map of message object from Client to Server
var MapObjectDecode = map[int32]interface{}{
	1: LoginPayload{},
	2: RegistPayload{},
	3: AuthResponse{},
}
