package basemodel

// Observer
type Observer interface {
	//notify
	notify(sub Subject)
}
