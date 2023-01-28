package repo

type Repo interface {
	GetUser(id string) *User
	GetCard(userID string) *Card
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
