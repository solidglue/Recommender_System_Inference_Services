package basemodel

//Subject
type Subject interface {
	//notice
	NotifyObservers()

	//add observer
	AddObserver(observer Observer)

	//add observers
	AddObservers(observers ...Observer)

	//remove observer
	RemoveObserver(observer Observer)

	//clean observer
	RemoveAllObservers()
}

// modelSubject implement Subject interface
type modelSubject struct {
	// observers
	observers []Observer
}

// NotifyObservers
func (ss *modelSubject) NotifyObservers() {
	for _, o := range ss.observers {
		o.notify(ss)
	}
}

// AddObserver
func (ss *modelSubject) AddObserver(observer Observer) {
	ss.observers = append(ss.observers, observer)
}

// AddObservers
func (ss *modelSubject) AddObservers(observers ...Observer) {
	ss.observers = append(ss.observers, observers...)
}

// RemoveObserver
func (ss *modelSubject) RemoveObserver(observer Observer) {
	for i := 0; i < len(ss.observers); i++ {
		if ss.observers[i] == observer {
			ss.observers = append(ss.observers[:i], ss.observers[i+1:]...)
		}
	}
}

// RemoveAllObservers
func (ss *modelSubject) RemoveAllObservers() {
	ss.observers = ss.observers[:0]
}
