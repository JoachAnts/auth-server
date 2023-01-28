package repo

type Repo interface {
	GetUser(id string) *User
}

type User struct {
	ID   string
	Name string
	Role string
}

func NewRepo(users map[string]User) Repo {
	return &repo{
		users: users,
	}
}

type repo struct {
	users map[string]User
}

func (r *repo) GetUser(id string) *User {
	result, found := r.users[id]
	if !found {
		return nil
	}
	return &result
}
