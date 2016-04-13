package mangolog

import (
	"log"
	"time"

	. "github.com/paulbellamy/mango"
	"github.com/theplant/qortex/models/members"
)

// Started GET "/look/women" for 111.86.201.25 at Sun Apr 22 11:21:28 +0000 2012
// Completed 200 OK in 932ms (Views: 738.4ms | ActiveRecord: 183.9ms)

func MakeangLogger(debug bool) Middleware {
	return func(env Env, app App) (status Status, headers Headers, body Body) {

		startTime := time.Now()
		r := env.Request()
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			r.ParseForm()
		}

		log.Printf("Started %s \"%s\" for %s at %s\n", r.Method, r.RequestURI, r.RemoteAddr, startTime.String())

		if id, ok := r.Form["Id"]; ok {
			log.Printf("The Id is: %s \n", id)
		}

		// Log post form value when enabling the debug mode.
		if debug {
			log.Printf(" ##### \n Form %v; \n Accept: %+v \n\n", r.Form, r.Header["Accept"])
		}

		status, headers, body = app(env)

		member, _ := env["LOGGED_IN_MEMBER_KEY"].(*members.Member)
		if member != nil {
			execution := time.Now().Sub(startTime) / time.Millisecond
			log.Printf("%s Completed \"%s\" %d in %dms\n\n", member.Name()+" "+member.Id.Hex(), r.RequestURI, status, execution)
			return
		}

		execution := time.Now().Sub(startTime) / time.Millisecond
		log.Printf("Completed \"%s\" %d in %dms\n\n", r.RequestURI, status, execution)

		return
	}
}
