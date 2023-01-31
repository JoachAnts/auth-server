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
	}, map[string](map[string]repo.Card){})

	user1Result := r.GetUser(user1.ID)
	user2Result := r.GetUser(user2.ID)

	assert.NotNil(t, user1Result)
	assert.Equal(t, user1.ID, user1Result.ID)
	assert.Equal(t, user1.Name, user1Result.Name)

	assert.NotNil(t, user2Result)
	assert.Equal(t, user2.ID, user2Result.ID)
	assert.Equal(t, user2.Name, user2Result.Name)
}

func TestMutipleCards(t *testing.T) {
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
	}, map[string](map[string]repo.Card){
		"1": {
			"1": repo.Card{
				MaskedNumber: "**** **** **** 1111",
				Exp:          "12/23",
				Limit:        10000,
				Balance:      4321,
			},
			"2": repo.Card{
				MaskedNumber: "**** **** **** 2222",
				Exp:          "12/23",
				Limit:        10000,
				Balance:      4321,
			},
		},
		"2": {
			"1": repo.Card{
				MaskedNumber: "**** **** **** 3333",
				Exp:          "12/23",
				Limit:        10000,
				Balance:      4321,
			},
			"2": repo.Card{
				MaskedNumber: "**** **** **** 4444",
				Exp:          "12/23",
				Limit:        10000,
				Balance:      4321,
			},
		},
	})

	user1Result := r.GetCards(user1.ID)

	assert.NotNil(t, user1Result)
	assert.Equal(t, repo.Card{
		MaskedNumber: "**** **** **** 1111",
		Exp:          "12/23",
		Limit:        10000,
		Balance:      4321,
	}, user1Result["1"])
	assert.Equal(t, repo.Card{
		MaskedNumber: "**** **** **** 2222",
		Exp:          "12/23",
		Limit:        10000,
		Balance:      4321,
	}, user1Result["2"])

	user2Result := r.GetCards(user2.ID)

	assert.NotNil(t, user2Result)
	assert.Equal(t, repo.Card{
		MaskedNumber: "**** **** **** 3333",
		Exp:          "12/23",
		Limit:        10000,
		Balance:      4321,
	}, user2Result["1"])
	assert.Equal(t, repo.Card{
		MaskedNumber: "**** **** **** 4444",
		Exp:          "12/23",
		Limit:        10000,
		Balance:      4321,
	}, user2Result["2"])
}

func TestUserNotFound(t *testing.T) {
	r := repo.NewRepo(map[string]repo.User{}, map[string](map[string]repo.Card){})

	result := r.GetUser("1")

	assert.Nil(t, result)
}
