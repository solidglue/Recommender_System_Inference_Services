package infer_features

type InferFeature struct {
	userId          string
	userProfiles    map[string]string
	lastNLong       map[string]string //such as last n item ids ,last n item labels. thousands length.
	lastNShort      map[string]string //such as last n item ids ,last n item labels. hundreds length.
	contextFeatures map[string]string //such as geop, time, os
	itemId          string
	itemProfiles    map[string]string
}
