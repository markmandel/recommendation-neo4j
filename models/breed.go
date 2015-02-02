package models

import "net/url"

type Breed struct {
	Name        string
	Description string
	WikiURL     *url.URL
}

func (b *Breed) String() string {
	return "[Breed: " + b.Name + ", Description: " + b.Description + ", Wiki: " + b.WikiURL.String() + " ]"
}
