package model

type Room struct {
	Id          string `json:"id" bson:"_id"`
	Users       []User `json:"users" bson:"users"`
	HiddenVotes bool   `json:"hiddentVotes" bson:"hiddenVotes"`
}

type User struct {
	Id          string `json:"id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	LifeTimeEnd string `json:"lifeTimeEnd" bson:"lifeTimeEnd"`
	Vote        string `json:"vote" bson:"vote"`
}

type VoteDTO struct {
	Vote string `json:"vote"`
}
