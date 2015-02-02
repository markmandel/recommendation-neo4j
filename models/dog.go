package models

import "fmt"

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
