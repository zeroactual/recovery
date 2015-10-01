package recovery

import (
	"net/http"
	"github.com/zenazn/goji/web"
	"fmt"
	"github.com/zeroactual/templates"
)


func Recoverer(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {

			w.Header().Set("Access-Control-Allow-Origin", "*")

			fault := recover()
			if fault != nil {
				w.Header().Set("Content-Type", "text/html")

				err, ok := fault.(*fatal)
				if !ok {
					data := map[string]interface{}{
						"code": 500,
						"message": "Internal Server Error",
						"cause": fmt.Sprint(fault),
					}
					w.WriteHeader(500)
					templates.T.Render(w, "error.html", false, data)
					return
				}

				data := map[string]interface{}{
					"code": err.Code,
					"message": err.Message,
				}
				// Set Headers
				w.WriteHeader(err.Code)
				templates.T.Render(w, "error.html", false, data)
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

type fatal struct {
	Code int
	Message interface{}
}

func Terminate(code int, message interface{}) {
	panic(&fatal{
		code,
		message,
	})
}
