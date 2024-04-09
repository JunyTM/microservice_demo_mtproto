package model

type MessageSending struct {
	Salt       int64  `json:"salt"`
	SessionId  int64  `json:"session_id"`
	MessageId  int64  `json:"message_id"`
	SeqNo      int32  `json:"seq_no"`
	MessageLen int32  `json:"message_len"`
	Body       []byte `json:"body"`
}
