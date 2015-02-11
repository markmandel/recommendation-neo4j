package main

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
		showLookedAtDogs = false
		showPersonalRecommendations = false

	case "looked": //r=looked
		showLookedAtDogs = true

	case "slope": //r=slope
		showPersonalRecommendations = true
	}

}
