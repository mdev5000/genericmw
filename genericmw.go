package genericmw

import "net/http"

type AppHandler[T any] func(T, http.ResponseWriter, *http.Request)

type Middleware[T any] func(next AppHandler[T]) AppHandler[T]

func Wrap[T any](newFn func() T, next AppHandler[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(newFn(), w, r)
	}
}

type Middlewares[T any] struct {
	newFn       func() T
	middlewares []Middleware[T]
}

func NewMiddlewares[T any](newFn func() T) *Middlewares[T] {
	return &Middlewares[T]{newFn: newFn}
}

func (mws *Middlewares[T]) Use(mw Middleware[T]) {
	mws.middlewares = append(mws.middlewares, mw)
}

func (mws *Middlewares[T]) Wrap(handler AppHandler[T]) http.Handler {
	appHandler := handler
	for _, m := range mws.middlewares {
		appHandler = m(appHandler)
	}
	return Wrap[T](mws.newFn, appHandler)
}
