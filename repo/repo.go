package repo

type Repo interface {
	GetUser(id string) *User
	GetCard(userID string) *Card
	SetCardLimit(userID string, newLimit int) *Card
}

type User struct {
	ID   string
	Name string
	Role string
}

type Card struct {
	MaskedNumber string
	Limit        int
	Balance      int
	Exp          string
}

func NewRepo(users map[string]User, cards map[string]Card) Repo {
	return &repo{
		users: users,
		cards: cards,
	}
}

type repo struct {
	users map[string]User
	cards map[string]Card
}

func (r *repo) GetUser(id string) *User {
	result, found := r.users[id]
	if !found {
		return nil
	}
	return &result
}

func (r *repo) GetCard(id string) *Card {
	result, found := r.cards[id]
	if !found {
		return nil
	}
	return &result
}

func (r *repo) SetCardLimit(userID string, newLimit int) *Card {
	result, found := r.cards[userID]
	if !found {
		return nil
	}
	result.Limit = newLimit
	r.cards[userID] = result
	return &result
}
