package rest

import (
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tracy/log"
	"tracy/proxy"
	"tracy/tracer/common"
	"tracy/tracer/types"
)

/*AddTracer decodes an HTTP request to add a new tracer(s) to the database. */
func AddTracers(w http.ResponseWriter, r *http.Request) {
	ret := []byte("{}")
	status := http.StatusInternalServerError
	in := types.Request{}
	json.NewDecoder(r.Body).Decode(&in)

	var err error
	if ret, err = common.AddTracer(in); err == nil {
		ret = ServerError(err)
		log.Error.Printf(err.Error())
	} else {
		status = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(ret)
}

/*GetTracers gets all the tracers. */
func GetTracers(w http.ResponseWriter, r *http.Request) {
	ret := []byte("{}")
	status := http.StatusInternalServerError

	var err error
	if ret, err = common.GetTracers(); err != nil {
		ret = ServerError(err)
		log.Error.Printf(err.Error())
	} else {
		status = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(ret)
}

/*GetTracer gets the tracer data structure belonging to the ID in the URL. */
func GetTracer(w http.ResponseWriter, r *http.Request) {
	ret := []byte("{}")
	status := http.StatusInternalServerError

	vars := mux.Vars(r)
	if tracerID, ok := vars["tracerID"]; ok {
		id, err := strconv.ParseInt(tracerID, 10, 32)
		if err != nil {
			ret = ServerError(err)
			log.Error.Printf(err.Error())
		} else {
			if ret, err = common.GetTracer(int(id)); err != nil {
				ret = ServerError(err)
				log.Error.Printf(err.Error())
			} else {
				status = http.StatusOK
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(ret)
}

/*GenerateTracer generates a new tracer and stored it in the database. Often used for
 * frontend heavy applications that might start using the input right away before
 * sending a request to the Tracy proxy. */
func GenerateTracer(w http.ResponseWriter, r *http.Request) {
	ret := []byte("{}")
	status := http.StatusInternalServerError
	var err error

	r.ParseForm()
	tracerString := r.Form.Get("tracer_string")
	if len(tracerString) != 0 {
		requestURL := r.Form.Get("url")
		if len(requestURL) != 0 {
			_, payload := proxy.GenerateTracerFromTag(tracerString)
			if payload != nil {
				//TODO: should collect more information about the location of where
				// it was generated. generating a tracer like this loses information
				// about inputs without being obvious about it. if we wanted to do
				// reproduction steps, how would we do that here?
				genTracer := types.Request{
					RawRequest:    "GENERATED", // For generated tracers, there won't be a request
					RequestMethod: "GENERATED", // For generated tracers, there won't be a request method
					RequestURL:    requestURL,
					Tracers: []types.Tracer{
						types.Tracer{
							TracerString: tracerString,
						},
					},
				}

				if ret, err = common.AddTracer(genTracer); err != nil {
					ret = ServerError(err)
					log.Error.Printf(err.Error())
				} else {
					status = http.StatusOK
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(ret)
}
