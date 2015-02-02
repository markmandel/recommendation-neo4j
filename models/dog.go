package models

import "fmt"

// A dog that is up for adoption
type Dog struct {
	Id      int `json:-`
	Name    string
	Breed   *Breed
	Adopted bool
	PicURL  string
}

func (d Dog) String() string {
	return fmt.Sprintf("[ Name: %v, Adopted: %v, PicURL: %v Breed: %v ]", d.Name, d.Adopted, d.PicURL, d.Breed)
}
