package GoReactive


type Disposable interface {
  Dispose()
}


type Observer interface {
  Next(value interface{})
  Error(err error)
  Completion()
}


type Observable interface {
  Subscribe(next func(interface{}), completion func(), failure func(error)) Disposable
}


// Start with

type startWithObservable struct {
  initial interface{}
  observable Observable
}

func (observable startWithObservable) Subscribe(next func(interface{}), completion func(), failure func(error))  Disposable {
  next(observable.initial)
  return observable.observable.Subscribe(next, completion, failure)
}

func StartWith(observable Observable, value interface{}) Observable {
  return startWithObservable{ observable: observable, initial: value }
}

// Skip

type skipObservable struct {
  until int
  observable Observable
}


func (observable skipObservable) Subscribe(next func(interface{}), completion func(), failure func(error)) Disposable {
  closure := func(value interface{}) {
    if observable.until > 0 {
      observable.until--
    } else {
      next(value)
    }
  }

  return observable.observable.Subscribe(closure, completion, failure)
}

// Distinct Until Changed

type distrinctUntilChangedObservable struct {
  observable Observable
  previous interface{}
}

func (observable distrinctUntilChangedObservable) Subscribe(next func(interface{}), completion func(), failure func(error)) Disposable {
  closure := func(value interface{}) {
    if observable.previous != value {
      next(value)
      observable.previous = value
    }
  }

  return observable.observable.Subscribe(closure, completion, failure)
}

/// Returns an observable of values which are not equal to the previous value
func DistinctUntilChanged(observable Observable) Observable {
  return distrinctUntilChangedObservable{ observable: observable }
}


/// Returns an Observable that returns the given Observables values after skipping the given until amount of next's
func Skip(observable Observable, until int) Observable {
  return skipObservable{ observable: observable, until: until }
}

// Map

type mappedObservable struct {
  observable Observable
  transform func(interface{}) interface{}
}


func (observable mappedObservable) Subscribe(next func(interface{}), completion func(), failure func(error)) Disposable {
  return observable.observable.Subscribe(func(value interface{}) {
                                    next(observable.transform(value))
                                  },
                                  completion, failure)
}


func Map(observable Observable, transform func(interface{}) interface{}) Observable {
  return mappedObservable{ observable: observable, transform: transform }
}


// Filter

type filterObservable struct {
  observable Observable
  filter func(interface{}) bool
}


func (observable filterObservable) Subscribe(next func(interface{}), completion func(), failure func(error)) Disposable {
  return observable.observable.Subscribe(func(value interface{}) {
                                    if observable.filter(value) {
                                      next(value)
                                    }
                                  },
                                  completion, failure)
}


func Filter(observable Observable, filter func(interface{}) bool) Observable {
  return filterObservable{ observable: observable, filter: filter }
}


func Exclude(observable Observable, exclude func(interface{}) bool) Observable {
  return Filter(observable, func(value interface{}) bool { return !exclude(value) })
}

