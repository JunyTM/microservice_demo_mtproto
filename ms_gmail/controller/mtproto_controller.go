package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ms_gmail/infrastructure"
	"ms_gmail/model"
	"ms_gmail/pb"
	"ms_gmail/utils"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MTProtoController interface {
	Send(w http.ResponseWriter, r *http.Request)
}

type mtProtoController struct{}

func (s *mtProtoController) Send(w http.ResponseWriter, r *http.Request) {
	var res model.Response
	var mtProtoPayload model.MTprotoPayload

	err := json.NewDecoder(r.Body).Decode(&mtProtoPayload)
	if err != nil {
		BadRequest(w, r, err)
		return
	}

	// Connnect to server auth
	conn, err := grpc.Dial(serverHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	mtProto := pb.NewEncryptedServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(),3 * time.Second)
	defer cancel()

	// Required param salt, session_id, message_id, seq_no, message_len > 0 for Trimming zero padding when ecode data
	// Serialize message 
	data := &model.MessageSending{
		Salt:       1,
        SessionId:  1,
        MessageId:  1,
        SeqNo:      1,
        MessageLen: int32(len([]byte(mtProtoPayload.Message))),
        Body:       []byte(mtProtoPayload.Message),
	}
	serializedMessage, err := utils.SerializeMarshal(*data)
	if err != nil {
		InternalServerError(w, r, err)
		return
	}

	// Calculate message AES encryption key and IV
	authKey := infrastructure.GetAuthKey()
	aesKey, aesIV, msgKey, err := utils.ComputeAESKeyIV(authKey, serializedMessage)
	if err != nil {
		InternalServerError(w, r, err)
		return
	}

	// Encrypt message
	cypherText, err := utils.EnscriptAES_IGE(aesKey, aesIV, serializedMessage)
	if err != nil {
		InternalServerError(w, r, err)
		return
	}

	// Send message to server with encrypted
	serverResponse, err := mtProto.Send(ctx, &pb.Message{
		AuthenId: utils.GetAuthKeyId(authKey),
		MessageKey: msgKey,
		SerializeData: cypherText,
	})
	if err != nil {
		log.Fatalf("> gRPC error: %v", err)
	}


    res.Data = serverResponse
    res.Success = true
    res.Message = "Server responded"
    render.JSON(w, r, res)


	// ReCaculate the new AES key and IV by msgKey
	aesKey_, aesIV_, err := utils.ComputeAESKeyIV2([]byte(authKey), serverResponse.GetMessageKey())
	if err != nil {
		InternalServerError(w, r, err)
		return
	}

	// Decrypt message
	decryptedMessage, err := utils.DescriptAES_IGE(aesKey_, aesIV_, serverResponse.GetSerializeData())
	if err != nil {
		InternalServerError(w, r, err)
		return
	}

	// Deserialize message
	newData, err := utils.SerializeUnMarshal(decryptedMessage)
	if err != nil {
		InternalServerError(w, r, err)
		return
	}

	res = model.Response{
		Data:   fmt.Sprintf("gRPC Message reply ==> %s", string(newData.Body)),
		Success: true,
		Message: "Pong",
	}
	render.JSON(w, r, res)
}

func NewMTprotoController() MTProtoController {
	return &mtProtoController{}
}
