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
	middlewares []Middleware[T]
}

func NewMiddlewares[T any]() *Middlewares[T] {
	return &Middlewares[T]{}
}

func (mws *Middlewares[T]) Use(mw Middleware[T]) {
	mws.middlewares = append(mws.middlewares, mw)
}

func (mws *Middlewares[T]) Wrap(newFn func() T, handler AppHandler[T]) http.Handler {
	appHandler := handler
	for _, m := range mws.middlewares {
		appHandler = m(appHandler)
	}
	return Wrap[T](newFn, appHandler)
}
