package users

import (
	"fmt"
	"github.com/google/uuid"
)

type UserInserter interface {
	Insert(user *User) error
}

type UserUpdater interface {
	UpdateEmail(user *User) error
}

type UserDeleter interface {
	Delete(id uuid.UUID, delEmail string) error
}

type UserGetter interface {
	GetById(id uuid.UUID) (*User, error)
}

type userRepository interface {
	UserInserter
	UserUpdater
	UserDeleter
	UserGetter
}

// UserService aggregates the methods a user may need to operate over the usersrepository.
// It expects an implementation of the following interfaces:
//
// # UserInserter which enforces a user insertion into the system
//
// # UserUpdater for the purpose of updating user's email
//
// UserDeleter to delete an existing user from the database.
type UserService struct {
	userRepository userRepository
}

func NewUserService(userRepository userRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

// newDefaultUser creates a new user with the provided email and id.
// Since we use Supabase for auth, and we want to control other fields in our db,
// we do not include a username or other fields here.
// The user is created with the RegularRole by default.
func newDefaultUser(email string, id uuid.UUID) *User {
	return &User{
		ID:        id,
		Email:     email,
		Role:      RegularRole,
		IsDeleted: false,
	}
}

// InsertDefaultUser creates a new default user with basic permissions and
// inserts it into the database using Inserter. The user is defined by the specified email and id.
// Returns an error if the user could not be inserted into the repository.
// The error wraps the underlying error returned by the Inserter.Insert method.
func (u *UserService) InsertDefaultUser(email string, id uuid.UUID) error {
	// Create a new default user with basic role that has limited permissions
	user := newDefaultUser(email, id)

	// Insert this role into the db
	err := u.userRepository.Insert(user)
	if err != nil {
		return fmt.Errorf("could not insert user: %w", err)
	}

	return nil
}

// UpdateUserEmail updates the email of a user in the repository using the UserUpdater.
// The email is identified by the uuid provided. The new email for the user is also provided.
// An error is returned if the email could not be updated in the repository.
// The error wraps the underlying error returned by the UserUpdater.UpdateEmail method.
func (u *UserService) UpdateUserEmail(userId uuid.UUID, email string) error {
	user := User{ID: userId, Email: email}

	err := u.userRepository.UpdateEmail(&user)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser removes a user from the repository. The user to be removed is identified
// by the provided UUID. If there is a problem removing the user, an error will be returned
// which wraps the underlying error returned by the UserDeleter's Delete method.
func (u *UserService) DeleteUser(userId uuid.UUID, oldEmail string) error {
	delEmail := fmt.Sprintf("%v+%v", oldEmail, userId)

	err := u.userRepository.Delete(userId, delEmail)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}

func (u *UserService) GetById(userId uuid.UUID) (*User, error) {
	user, err := u.userRepository.GetById(userId)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}
