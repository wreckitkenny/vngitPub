package model

// User structs user fields
type User struct {
	Username 	string `bson:"username,omitempty" json:"username,omitempty"`
	Password 	string `bson:"password,omitempty" json:"password,omitempty"`
	Department 	string `bson:"department,omitempty" json:"department,omitempty"`
	FullName	string `bson:"fullname,omitempty" json:"fullname,omitempty"`
	Email    	string `bson:"email,omitempty" json:"email,omitempty"`
	Role		string `bson:"role,omitempty" json:"role,omitempty"`
}