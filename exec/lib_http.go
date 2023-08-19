package exec

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"sync"
)

var lock = sync.Mutex{}

func listenAndServe(args []*OwlObj) (*OwlObj, bool) {
	this := args[0]
	port := args[1].TrueStr()

	routes, ok := this.GetAttr("routes")

	if !ok {
		return NewString("Attribute 'routes' not found"), false
	}

	urls := routes.Attr

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lock.Lock()
		defer lock.Unlock()

		res := transformResponse(w)
		req := transformRequest(r)

		url := r.URL.String()

		for matcher, handler := range urls {
			matched, err := regexp.MatchString(matcher, url)

			if !matched || err != nil {
				continue
			}

			arg := NewList([]*OwlObj{res, req})
			handler.Call(arg)
			break
		}
	})

	err := http.ListenAndServe(port, nil)

	return NewString(err.Error()), true
}

func HttpLibExport() *OwlObj {
	o := NewOwlObj()

	o.SetAttr("ListenAndServe", NewCallBridge(listenAndServe))

	return o
}

func transformRequest(r *http.Request) *OwlObj {
	o := NewOwlObj()

	o.SetAttr("Method", NewString(r.Method))
	o.SetAttr("URL", NewString(r.URL.String()))
	o.SetAttr("Header", transformHeader(r.Header))
	o.SetAttr("Body", transformBody(r.Body))
	o.SetAttr("ContentLength", NewInt(r.ContentLength))

	return o
}

func transformHeader(h http.Header) *OwlObj {
	o := NewOwlObj()

	for k, v := range h {
		o.SetAttr(k, NewString(v[0]))
	}

	return o
}

func transformBody(f io.ReadCloser) *OwlObj {
	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	newStr := buf.String()

	return NewString(newStr)
}

func transformResponse(w http.ResponseWriter) *OwlObj {
	o := NewOwlObj()

	o.SetAttr("SetHeader", NewCallBridge(func(args []*OwlObj) (*OwlObj, bool) {
		name := args[1].TrueStr()
		value := args[2].TrueStr()

		w.Header().Set(name, value)

		return nil, true
	}))

	o.SetAttr("SetStatus", NewCallBridge(func(args []*OwlObj) (*OwlObj, bool) {
		code, ok := args[1].TrueInt()

		if !ok {
			return NewString("Status is not an int."), false
		}

		w.WriteHeader(int(code))
		return nil, true
	}))

	o.SetAttr("SetBody", NewCallBridge(func(args []*OwlObj) (*OwlObj, bool) {
		w.Write([]byte(args[1].TrueStr()))
		return nil, true
	}))

	return o
}
