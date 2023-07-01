package users

// https://github.com/stretchr/testify
import (
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var defaultTestUser = User{Id: uuid.New(), Mail: "some@mail", Name: "someName", Password: "somePassword"}

type repositoryGenerator func() Repository

// https://pkg.go.dev/github.com/stretchr/testify/suite
type UsersRepositoryTestSuite struct {
	suite.Suite
	createRepo repositoryGenerator
}

func (suite *UsersRepositoryTestSuite) TestCreateUser_InvalidMail() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()
	_, err := repo.Create(User{})
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidUserMail))
}

func (suite *UsersRepositoryTestSuite) TestCreateUser_InvalidName() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()
	user := User{Mail: "some@mail"}
	_, err := repo.Create(user)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidUserName))
}

func (suite *UsersRepositoryTestSuite) TestCreateUser_InvalidPassword() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()
	user := User{Mail: "some@mail", Name: "someName"}
	_, err := repo.Create(user)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidPassword))
}

func (suite *UsersRepositoryTestSuite) TestCreateUser() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()
	id, err := repo.Create(defaultTestUser)
	assert.Nil(err)

	_, err = uuid.Parse(id.String())
	assert.Nil(err)
	assert.Equal(defaultTestUser.Id, id)
}

func (suite *UsersRepositoryTestSuite) TestCreateUser_Duplicated() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()
	_, err := repo.Create(defaultTestUser)
	assert.Nil(err)

	_, err = repo.Create(defaultTestUser)
	assert.True(errors.IsErrorWithCode(err, errors.ErrUserAlreadyExists))
}

func (suite *UsersRepositoryTestSuite) TestGetUser_NoUsers() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()

	someId := uuid.New()
	_, err := repo.Get(someId)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}

func (suite *UsersRepositoryTestSuite) TestGetUser_WrongId() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()
	repo.Create(defaultTestUser)

	wrongId := uuid.New()
	assert.NotEqual(defaultTestUser.Id, wrongId)
	_, err := repo.Get(wrongId)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}

func (suite *UsersRepositoryTestSuite) TestGetUser() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()
	repo.Create(defaultTestUser)

	userFromRepo, err := repo.Get(defaultTestUser.Id)
	assert.Nil(err)

	assert.Equal(defaultTestUser.Id, userFromRepo.Id)
	assert.Equal(defaultTestUser.Mail, userFromRepo.Mail)
	assert.Equal(defaultTestUser.Name, userFromRepo.Name)
	assert.Equal(defaultTestUser.Password, userFromRepo.Password)
}

func (suite *UsersRepositoryTestSuite) TestGetAll() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()

	ids, err := repo.GetAll()
	assert.Nil(err)
	assert.Equal(0, len(ids))

	repo.Create(defaultTestUser)

	ids, err = repo.GetAll()
	assert.Nil(err)
	assert.Equal(1, len(ids))
	assert.Equal(defaultTestUser.Id, ids[0])
}

func (suite *UsersRepositoryTestSuite) TestPatchUser_Empty() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()
	repo.Create(defaultTestUser)

	patch := User{}
	patched, err := repo.Patch(defaultTestUser.Id, patch)
	assert.Nil(err)

	assert.Equal(defaultTestUser.Id, patched.Id)
	assert.Equal(defaultTestUser.Mail, patched.Mail)
	assert.Equal(defaultTestUser.Name, patched.Name)
	assert.Equal(defaultTestUser.Password, patched.Password)
}

func (suite *UsersRepositoryTestSuite) TestPatchUser_WrongId() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()
	repo.Create(defaultTestUser)

	wrongId := uuid.New()
	assert.NotEqual(defaultTestUser.Id, wrongId)
	patch := User{}
	_, err := repo.Patch(wrongId, patch)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}

func (suite *UsersRepositoryTestSuite) TestPatchUser_Mail() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()
	repo.Create(defaultTestUser)

	patch := User{Mail: "some-other@mail"}
	patched, err := repo.Patch(defaultTestUser.Id, patch)
	assert.Nil(err)

	assert.Equal(defaultTestUser.Id, patched.Id)
	assert.Equal(patch.Mail, patched.Mail)
	assert.Equal(defaultTestUser.Name, patched.Name)
	assert.Equal(defaultTestUser.Password, patched.Password)
}

func (suite *UsersRepositoryTestSuite) TestPatchUser_Name() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()
	repo.Create(defaultTestUser)

	patch := User{Name: "someOtherName"}
	patched, err := repo.Patch(defaultTestUser.Id, patch)
	assert.Nil(err)

	assert.Equal(defaultTestUser.Id, patched.Id)
	assert.Equal(defaultTestUser.Mail, patched.Mail)
	assert.Equal(patch.Name, patched.Name)
	assert.Equal(defaultTestUser.Password, patched.Password)
}

func (suite *UsersRepositoryTestSuite) TestPatchUser_Password() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()
	repo.Create(defaultTestUser)

	patch := User{Password: "somePassword"}
	patched, err := repo.Patch(defaultTestUser.Id, patch)
	assert.Nil(err)

	assert.Equal(defaultTestUser.Id, patched.Id)
	assert.Equal(defaultTestUser.Mail, patched.Mail)
	assert.Equal(defaultTestUser.Name, patched.Name)
	assert.Equal(patch.Password, patched.Password)
}

func (suite *UsersRepositoryTestSuite) TestPatchUser_IdIsNotPatched() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()
	repo.Create(defaultTestUser)

	patch := User{Id: uuid.New()}
	assert.NotEqual(defaultTestUser.Id, patch.Id)
	patched, err := repo.Patch(defaultTestUser.Id, patch)
	assert.Nil(err)

	assert.Equal(defaultTestUser.Id, patched.Id)
	assert.Equal(defaultTestUser.Mail, patched.Mail)
	assert.Equal(defaultTestUser.Name, patched.Name)
	assert.Equal(defaultTestUser.Password, patched.Password)
}

func (suite *UsersRepositoryTestSuite) TestDeleteUser_NoUsers() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()

	someId := uuid.New()
	err := repo.Delete(someId)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}

func (suite *UsersRepositoryTestSuite) TestDeleteUser_WrongId() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()
	repo.Create(defaultTestUser)

	wrongId := uuid.New()
	assert.NotEqual(defaultTestUser.Id, wrongId)
	err := repo.Delete(wrongId)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}

func (suite *UsersRepositoryTestSuite) TestDeleteUser() {
	assert := assert.New(suite.T())

	repo := suite.createRepo()
	repo.Create(defaultTestUser)

	err := repo.Delete(defaultTestUser.Id)
	assert.Nil(err)

	_, err = repo.Get(defaultTestUser.Id)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}
