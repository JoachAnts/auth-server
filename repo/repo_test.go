package repo_test

import (
	"testing"

	"github.com/JoachAnts/auth-server/repo"
	"github.com/stretchr/testify/assert"
)

func TestRepo(t *testing.T) {
	user1 := repo.User{
		ID:   "1",
		Name: "John Smith",
	}
	user2 := repo.User{
		ID:   "2",
		Name: "Oliver Smith",
	}
	r := repo.NewRepo(map[string]repo.User{
		user1.ID: user1,
		user2.ID: user2,
	})

	user1Result := r.GetUser(user1.ID)
	user2Result := r.GetUser(user2.ID)

	assert.NotNil(t, user1Result)
	assert.Equal(t, user1.ID, user1Result.ID)
	assert.Equal(t, user1.Name, user1Result.Name)

	assert.NotNil(t, user2Result)
	assert.Equal(t, user2.ID, user2Result.ID)
	assert.Equal(t, user2.Name, user2Result.Name)
}

func TestUserNotFound(t *testing.T) {
	r := repo.NewRepo(map[string]repo.User{})

	result := r.GetUser("1")

	assert.Nil(t, result)
}
