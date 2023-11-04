package entity

type Event struct {
	Name         string `bson:"name"`
	Data         string `bson:"data"`
	Description  string `bson:"description"`
	Registration string `bson:"registration"`
}
