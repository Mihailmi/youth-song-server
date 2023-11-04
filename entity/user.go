package entity

type User struct {
	Name     string `bson:"name"`
	Surname  string `bson:"surname"`
	Email    string `bson:"email"`
	Login    string `bson:"login"`
	Password string `bson:"password"`
	Role     string `bson:"role"`
	JobTitle string `bson:"jobTilte"`
	IsActive bool   `bson:"isActive"`
}
