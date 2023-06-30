package users

// https://github.com/stretchr/testify
import (
	"testing"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddUser_InvalidName(t *testing.T) {
	assert := assert.New(t)

	udb := NewUserManager()
	_, err := udb.AddUser("", "")
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidUserName))
}

func TestAddUser_InvalidPassword(t *testing.T) {
	assert := assert.New(t)

	udb := NewUserManager()
	_, err := udb.AddUser("foo", "")
	assert.True(errors.IsErrorWithCode(err, errors.ErrInvalidPassword))
}

func TestAddUser(t *testing.T) {
	assert := assert.New(t)

	udb := NewUserManager()
	id, err := udb.AddUser("foo", "haha")
	assert.Nil(err)

	check, err := uuid.Parse(id.String())
	assert.Nil(err)
	assert.Equal(check, id)
}

func TestAddUser_Duplicated(t *testing.T) {
	assert := assert.New(t)

	udb := NewUserManager()
	_, err := udb.AddUser("foo", "haha")
	assert.Nil(err)

	_, err = udb.AddUser("foo", "haha")
	assert.True(errors.IsErrorWithCode(err, errors.ErrUserAlreadyExists))
}

func TestGetUsers(t *testing.T) {
	assert := assert.New(t)

	udb := NewUserManager()

	ids := udb.GetUsers()
	assert.Equal(0, len(ids))

	id, _ := udb.AddUser("foo", "haha")

	ids = udb.GetUsers()
	assert.Equal(1, len(ids))
	assert.Equal(id, ids[0])
}

func TestGetUser_NoUsers(t *testing.T) {
	assert := assert.New(t)

	udb := NewUserManager()

	wrongId := uuid.New()
	_, err := udb.GetUser(wrongId)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}

func TestGetUser_WrongId(t *testing.T) {
	assert := assert.New(t)

	udb := NewUserManager()
	udb.AddUser("foo", "haha")

	wrongId := uuid.New()
	_, err := udb.GetUser(wrongId)
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}

func TestGetUser(t *testing.T) {
	assert := assert.New(t)

	udb := NewUserManager()
	id, _ := udb.AddUser("foo", "haha")

	user, err := udb.GetUser(id)
	assert.Nil(err)

	assert.Equal(id, user.Id)
	assert.Equal("foo", user.Name)
	assert.Equal("haha", user.Password)
}

func TestGetUserFromName_NoUsers(t *testing.T) {
	assert := assert.New(t)

	udb := NewUserManager()

	_, err := udb.GetUserFromName("foo")
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}

func TestGetUserFromname_WrongName(t *testing.T) {
	assert := assert.New(t)

	udb := NewUserManager()
	udb.AddUser("foo", "haha")

	_, err := udb.GetUserFromName("")
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))

	_, err = udb.GetUserFromName("food")
	assert.True(errors.IsErrorWithCode(err, errors.ErrNoSuchUser))
}

func TestGetUserFromName(t *testing.T) {
	assert := assert.New(t)

	udb := NewUserManager()
	id, _ := udb.AddUser("foo", "haha")

	user, err := udb.GetUserFromName("foo")
	assert.Nil(err)

	assert.Equal(id, user.Id)
	assert.Equal("foo", user.Name)
	assert.Equal("haha", user.Password)
}
