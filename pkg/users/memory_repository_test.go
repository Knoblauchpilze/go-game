package users

import (
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultTestUser = User{Id: uuid.New(), Mail: "some@mail", Name: "someName", Password: "somePassword"}

func TestMemoryRepository_CreateUser_InvalidMail(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()
	_, err := repo.Create(User{})
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidUserMail))
}

func TestMemoryRepository_CreateUser_InvalidName(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()
	user := User{Mail: "some@mail"}
	_, err := repo.Create(user)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidUserName))
}

func TestMemoryRepository_CreateUser_InvalidPassword(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()
	user := User{Mail: "some@mail", Name: "someName"}
	_, err := repo.Create(user)
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidPassword))
}

func TestMemoryRepository_CreateUser(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()
	id, err := repo.Create(defaultTestUser)
	assert.Nil(err)

	_, err = uuid.Parse(id.String())
	assert.Nil(err)
	assert.Equal(defaultTestUser.Id, id)
}

func TestMemoryRepository_CreateUser_Duplicated(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()
	_, err := repo.Create(defaultTestUser)
	assert.Nil(err)

	_, err = repo.Create(defaultTestUser)
	assert.True(errors.IsErrorWithCode(err, errors.ErrUserAlreadyExists))
}

func TestMemoryRepository_GetUser_NoUsers(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()

	someId := uuid.New()
	_, err := repo.Get(someId)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}

func TestMemoryRepository_GetUser_WrongId(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()
	repo.Create(defaultTestUser)

	wrongId := uuid.New()
	assert.NotEqual(defaultTestUser.Id, wrongId)
	_, err := repo.Get(wrongId)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}

func TestMemoryRepository_GetUser(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()
	repo.Create(defaultTestUser)

	userFromRepo, err := repo.Get(defaultTestUser.Id)
	assert.Nil(err)

	assert.Equal(defaultTestUser.Id, userFromRepo.Id)
	assert.Equal(defaultTestUser.Mail, userFromRepo.Mail)
	assert.Equal(defaultTestUser.Name, userFromRepo.Name)
	assert.Equal(defaultTestUser.Password, userFromRepo.Password)
}

func TestMemoryRepository_GetAll(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()

	ids, err := repo.GetAll()
	assert.Nil(err)
	assert.Equal(0, len(ids))

	repo.Create(defaultTestUser)

	ids, err = repo.GetAll()
	assert.Nil(err)
	assert.Equal(1, len(ids))
	assert.Equal(defaultTestUser.Id, ids[0])
}

func TestMemoryRepository_PatchUser_Empty(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()
	repo.Create(defaultTestUser)

	patch := User{}
	patched, err := repo.Patch(defaultTestUser.Id, patch)
	assert.Nil(err)

	assert.Equal(defaultTestUser.Id, patched.Id)
	assert.Equal(defaultTestUser.Mail, patched.Mail)
	assert.Equal(defaultTestUser.Name, patched.Name)
	assert.Equal(defaultTestUser.Password, patched.Password)
}

func TestMemoryRepository_PatchUser_WrongId(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()
	repo.Create(defaultTestUser)

	wrongId := uuid.New()
	assert.NotEqual(defaultTestUser.Id, wrongId)
	patch := User{}
	_, err := repo.Patch(wrongId, patch)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}

func TestMemoryRepository_PatchUser_Mail(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()
	repo.Create(defaultTestUser)

	patch := User{Mail: "some-other@mail"}
	patched, err := repo.Patch(defaultTestUser.Id, patch)
	assert.Nil(err)

	assert.Equal(defaultTestUser.Id, patched.Id)
	assert.Equal(patch.Mail, patched.Mail)
	assert.Equal(defaultTestUser.Name, patched.Name)
	assert.Equal(defaultTestUser.Password, patched.Password)
}

func TestMemoryRepository_PatchUser_Name(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()
	repo.Create(defaultTestUser)

	patch := User{Name: "someOtherName"}
	patched, err := repo.Patch(defaultTestUser.Id, patch)
	assert.Nil(err)

	assert.Equal(defaultTestUser.Id, patched.Id)
	assert.Equal(defaultTestUser.Mail, patched.Mail)
	assert.Equal(patch.Name, patched.Name)
	assert.Equal(defaultTestUser.Password, patched.Password)
}

func TestMemoryRepository_PatchUser_Password(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()
	repo.Create(defaultTestUser)

	patch := User{Password: "somePassword"}
	patched, err := repo.Patch(defaultTestUser.Id, patch)
	assert.Nil(err)

	assert.Equal(defaultTestUser.Id, patched.Id)
	assert.Equal(defaultTestUser.Mail, patched.Mail)
	assert.Equal(defaultTestUser.Name, patched.Name)
	assert.Equal(patch.Password, patched.Password)
}

func TestMemoryRepository_PatchUser_IdIsNotPatched(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()
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

func TestMemoryRepository_DeleteUser_NoUsers(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()

	someId := uuid.New()
	err := repo.Delete(someId)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}

func TestMemoryRepository_DeleteUser_WrongId(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()
	repo.Create(defaultTestUser)

	wrongId := uuid.New()
	assert.NotEqual(defaultTestUser.Id, wrongId)
	err := repo.Delete(wrongId)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}

func TestMemoryRepository_DeleteUser(t *testing.T) {
	assert := assert.New(t)

	repo := NewMemoryRepository()
	repo.Create(defaultTestUser)

	err := repo.Delete(defaultTestUser.Id)
	assert.Nil(err)

	_, err = repo.Get(defaultTestUser.Id)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}
