package tests

import (
	"testing"

	authService "paywise/internal/business/auth"

	"github.com/stretchr/testify/require"
)

func TestPasswordHashing(t *testing.T) {
	passForRegisteration := "123456789"

	hashed, err := authService.HashPassword(passForRegisteration)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	correctLoginPass := passForRegisteration
	isMatching := authService.CheckPassword(correctLoginPass, hashed)
	require.True(t, isMatching)

	wrongLoginPass := "159357852"
	notMatch := authService.CheckPassword(wrongLoginPass, hashed)
	require.False(t, notMatch)
}
