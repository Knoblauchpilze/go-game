package users

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestMemoryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &UsersRepositoryTestSuite{
		createRepo: func() Repository {
			return NewMemoryRepository()
		},
	})
}
