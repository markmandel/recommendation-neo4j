package main

import "log"

/*
Feature Flags
*/
var (
	flagQueryParam              = "r"
	showLookedAtDogs            = true
	showPersonalRecommendations = true
)

//processFlags processes URL variables and adjusts the feature flags accordingly.
func processFlags(r string) {

	switch r {
	case "none": //r=none
		log.Printf("Disabling all features")
		showLookedAtDogs = false
		showPersonalRecommendations = false

	case "looked": //r=looked
		log.Printf("Enabling Looked At Dogs")
		showLookedAtDogs = true

	case "slope": //r=slope
		log.Printf("Enabling Slope One Recommendations")
		showPersonalRecommendations = true
	}
}
