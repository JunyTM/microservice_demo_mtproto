package utils

import (
	"log"
	"ms_gmail/model"
	"testing"

	"github.com/stretchr/testify/require"
)

var auth_key = "ACA2547BAA3B4B3E6AD9129ADEE796042FC4ACAAABA40DCE7A1196640595829E2BA85C5E68E5D2B7B713B28720CCD0E65DA3B80CD8A9282BDE895755B924F4CBBCED119DEE057D3066D7095C25771A53D9AD6EC3464EE0E7FE1A4D9851F17B6CC7D6E636E0BDEA7D66153CF37005835A528E65673C1F5ADAA27465964BA89AA6FD045108B25B1DC9AB2D979864858C9EACBFA4A399CDBF2154D2CC6743C5CAA40683343D80A36C7F1289DC5AC5EE585DEA1E1CD296444EA6CD8F6C3564B2D355B310A3C02608EEA5EC36D38E941BF811489AB7A24C069E44FF54714CB4A45EEE2F310FDDC0A1B08BD3DFF8801E27C8508EFD7AD51EF3C13EF1D9A02EE4601741"
var messageBody ="Test MTproto system heloooo"

func TestSerialize(t *testing.T) {
	message := model.MessageSending{
		Salt:       123456789,
		SessionId:  123456789,
		MessageId:  123456789,
		SeqNo:      123456789,
		MessageLen: int32(len([]byte("The fk systemd lor"))),
		Body:       []byte("The fk systemd lor"),
	}

	serializedMessage, err := SerializeMarshal(message)
	require.NoError(t, err)

	deserializedMessage, err := SerializeUnMarshal(serializedMessage)
	require.NoError(t, err)
	log.Println("=> deserializedMessage: ", deserializedMessage)
	require.Equal(t, message.Body, deserializedMessage.Body)
}

func TestIGE(t *testing.T) {
	message := "The fk systemd, !!!!!!!!!!"

	key, iv, _, err := ComputeAESKeyIV(auth_key, []byte(message))
	require.NoError(t, err)

	ciphertext, err := EnscriptAES_IGE(key, iv, []byte(message))
	require.NoError(t, err)
	plaintext, err := DescriptAES_IGE(key, iv, ciphertext)
	require.NoError(t, err)
	require.Equal(t, message, string(plaintext))
}

// Required param salt, session_id, message_id, seq_no, message_len > 0 for Trimming zero padding when ecode data
func TestSendMessage(t *testing.T) {
	message := model.MessageSending{
		Salt:       123456789,
		SessionId:  123456789,
		MessageId:  123456789,
		SeqNo:      123456789,
		MessageLen: int32(len([]byte(messageBody))),
		Body:       []byte(messageBody),
	}

	// Serialize message
	serializedMessage, err := SerializeMarshal(message)
	require.NoError(t, err)

	// Calculate message AES encryption key and IV
	key, iv, msgKey, err := ComputeAESKeyIV(auth_key, serializedMessage)
	require.NoError(t, err)

	// Encrypt message
	dataEncrypted, err := EnscriptAES_IGE(key, iv, serializedMessage)
	require.NoError(t, err)

	log.Println("=> serverResponse: ", string(dataEncrypted))

	//Caculate the new AES key and IV by msgKey
	key_, iv_, err := ComputeAESKeyIV2([]byte(auth_key), msgKey)
	require.NoError(t, err)

	// Decrypt message
	decryptedMessage, err := DescriptAES_IGE(key_, iv_, dataEncrypted)
	require.NoError(t, err)

	// Deserialize message
	deserializedMessage, err := SerializeUnMarshal(decryptedMessage)
	require.NoError(t, err)
	require.NotEmpty(t, string(deserializedMessage.Body))
	require.Equal(t, string(message.Body), string(deserializedMessage.Body))
	log.Println("=> deserializedMessage: ", deserializedMessage)
}
