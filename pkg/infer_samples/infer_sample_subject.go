package infer_samples

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

// SampleSubject implement Subject interface
type SampleSubject struct {
	// observers
	observers []Observer
}

// NotifyObservers
func (ss *SampleSubject) NotifyObservers() {
	for _, o := range ss.observers {
		o.notify(ss)
	}
}

// AddObserver
func (ss *SampleSubject) AddObserver(observer Observer) {
	ss.observers = append(ss.observers, observer)
}

// AddObservers
func (ss *SampleSubject) AddObservers(observers ...Observer) {
	ss.observers = append(ss.observers, observers...)
}

// RemoveObserver
func (ss *SampleSubject) RemoveObserver(observer Observer) {
	for i := 0; i < len(ss.observers); i++ {
		if ss.observers[i] == observer {
			ss.observers = append(ss.observers[:i], ss.observers[i+1:]...)
		}
	}
}

// RemoveAllObservers
func (ss *SampleSubject) RemoveAllObservers() {
	ss.observers = ss.observers[:0]
}
