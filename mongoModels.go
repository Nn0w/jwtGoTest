package main

type userLoggedIn struct {
	GUID             string `json:"guid" bson:"guid" binding:"required"`
	RefreshTokenHash string `json:"rtokenid" bson:"rtokenid" binding:"required"`
}

type userDataRecord struct {
	GUID string `json:"guid" bson:"guid" binding:"required"`
	//Login string              `json:"login" bson:"login" binding:"required"`
	//Email           string    `json:"email" bson:"email" binding:"required"`
	//PasswordHash        string    `json:"password" bson:"password" binding:"required,min=8"`
	//Verified        bool      `json:"verified" bson:"verified"`
	//CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	//UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}
