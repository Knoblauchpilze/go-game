package users

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestDbRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &UsersRepositoryTestSuite{
		createRepo: func() Repository {
			return NewDbRepository()
		},
	})
}
