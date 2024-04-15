package utils

import (
	"ms_auth/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeMessage(t *testing.T) {
	user := model.AuthResponse{
		UserId:        1,
		Email:         "gmaileasdfads",
		ServerMessage: "Authen success",
	}

	buff, err := StructToBuffer(user, 3)
	require.NoError(t, err)
	require.NotNil(t, buff)

	require.NotEmpty(t, buff.Bytes())

	// Decode
	newData := buff.Bytes()
	deBuf, err := BufferToStruct(newData)
	require.NoError(t, err)

	decodeData := deBuf.(model.AuthResponse)
	require.Equal(t, user, decodeData)
}
