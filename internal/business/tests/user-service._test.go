package tests

import (
	"testing"

	userService "paywise/internal/business/user"

	"github.com/stretchr/testify/require"
)

func TestPasswordHashing(t *testing.T) {
	passForRegisteration := "123456789"

	hashed, err := userService.HashPassword(passForRegisteration)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	correctLoginPass := passForRegisteration
	isMatching := userService.CheckPassword(correctLoginPass, hashed)
	require.True(t, isMatching)

	wrongLoginPass := "159357852"
	notMatch := userService.CheckPassword(wrongLoginPass, hashed)
	require.False(t, notMatch)
}
