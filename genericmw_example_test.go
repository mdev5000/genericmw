package genericmw_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/mdev5000/genericmw"
)

type UserSetter interface {
	SetUserID(userId int)
}

func userMiddleware[T UserSetter](next genericmw.AppHandler[T]) genericmw.AppHandler[T] {
	return func(values T, w http.ResponseWriter, r *http.Request) {
		values.SetUserID(5)
		next(values, w, r)
	}
}

type Auth struct {
	LoggedIn bool
}

type AuthSetter interface {
	SetAuth(auth Auth)
}

func authMiddleware[T AuthSetter](next genericmw.AppHandler[T]) genericmw.AppHandler[T] {
	return func(values T, w http.ResponseWriter, r *http.Request) {
		values.SetAuth(Auth{LoggedIn: true})
		next(values, w, r)
	}
}

type Root struct {
	SomeRootThing string
	UserID        int
	Auth          Auth
}

func NewRoot() *Root { return &Root{} }

func (r *Root) SetAuth(auth Auth) {
	r.Auth = auth
}

func (r *Root) SetUserID(userId int) {
	r.UserID = userId
}

func myHandler(values *Root, w http.ResponseWriter, r *http.Request) {
	printType(values)
}

func printType(values interface{}) {
	fmt.Printf("%T: %+v\n", values, values)
}

func rootMiddleware(next genericmw.AppHandler[*Root]) genericmw.AppHandler[*Root] {
	return func(values *Root, w http.ResponseWriter, r *http.Request) {
		values.SomeRootThing = "value"
		next(values, w, r)
	}
}

type DifferentRoot struct {
	UserID int
	Auth   Auth
}

func NewDifferentRoot() *DifferentRoot { return &DifferentRoot{} }

func (r *DifferentRoot) SetUserID(userId int) {
	r.UserID = userId
}

func (r *DifferentRoot) SetAuth(auth Auth) {
	r.Auth = auth
}

func myHandler2(values *DifferentRoot, w http.ResponseWriter, r *http.Request) {
	printType(values)
}

func ExampleMiddlewares() {
	mw := genericmw.NewMiddlewares[*Root]()
	mw.Use(authMiddleware[*Root])
	mw.Use(userMiddleware[*Root])
	mw.Use(rootMiddleware)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/somepath", nil)
	mw.Wrap(NewRoot, myHandler).ServeHTTP(w, r)

	// Output: *genericmw_test.Root: &{SomeRootThing:value UserID:5 Auth:{LoggedIn:true}}
}

func ExampleWrap() {
	{
		h := genericmw.Wrap[*Root](NewRoot, authMiddleware(userMiddleware(myHandler)))
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/somepath", nil)
		h.ServeHTTP(w, r)
	}

	{
		h := genericmw.Wrap[*DifferentRoot](NewDifferentRoot, authMiddleware(userMiddleware(myHandler2)))
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/somepath", nil)
		h.ServeHTTP(w, r)
	}

	//Output: *genericmw_test.Root: &{SomeRootThing: UserID:5 Auth:{LoggedIn:true}}
	//*genericmw_test.DifferentRoot: &{UserID:5 Auth:{LoggedIn:true}}
}
