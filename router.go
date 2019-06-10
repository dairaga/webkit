package webkit

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/dairaga/log"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
)

var (
	_result     = reflect.TypeOf((*Result)(nil)).Elem()
	_respWriter = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	_request    = reflect.TypeOf((*http.Request)(nil))
	_protomsg   = reflect.TypeOf((*proto.Message)(nil)).Elem()
)

func _cast(in string, kind reflect.Kind) reflect.Value {
	if kind == reflect.String {
		return reflect.ValueOf(in)
	}

	if kind >= reflect.Int && kind <= reflect.Int64 {
		a, err := strconv.ParseInt(in, 0, 0)
		if err == nil {
			switch kind {
			case reflect.Int:
				return reflect.ValueOf(int(a))
			case reflect.Int8:
				return reflect.ValueOf(int8(a))
			case reflect.Int16:
				return reflect.ValueOf(int16(a))
			case reflect.Int32:
				return reflect.ValueOf(int32(a))
			case reflect.Int64:
				return reflect.ValueOf(int64(a))
			}
		}
	} else if kind >= reflect.Uint && kind <= reflect.Uint64 {
		a, err := strconv.ParseUint(in, 0, 0)
		if err == nil {
			switch kind {
			case reflect.Uint:
				return reflect.ValueOf(uint(a))
			case reflect.Uint8:
				return reflect.ValueOf(uint8(a))
			case reflect.Uint16:
				return reflect.ValueOf(uint16(a))
			case reflect.Uint32:
				return reflect.ValueOf(uint32(a))
			case reflect.Uint64:
				return reflect.ValueOf(uint64(a))
			}
		}
	} else if kind == reflect.Float32 {
		a, err := strconv.ParseFloat(in, 32)
		if err == nil {
			return reflect.ValueOf(float32(a))
		}
	} else if kind == reflect.Float64 {
		a, err := strconv.ParseFloat(in, 64)
		if err == nil {
			return reflect.ValueOf(a)
		}
	} else if kind == reflect.Bool {
		a, err := strconv.ParseBool(in)
		if err == nil {
			return reflect.ValueOf(a)
		}
	}

	return reflect.Value{}
}

func _isByteSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8
}

func _isPtrOfStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

func _indirectType(t reflect.Type) reflect.Type {
	elm := t
	for elm.Kind() == reflect.Ptr {
		elm = elm.Elem()
	}

	return elm
}

func _readData(r *http.Request, t reflect.Type) (reflect.Value, error) {

	if _isByteSlice(t) {
		bodybytes, err := ioutil.ReadAll(r.Body)
		return reflect.ValueOf(bodybytes), err
	}

	data := reflect.New(t.Elem())

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "application/octet-stream") {
		bodybytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("read body: ", err)
			return reflect.Value{}, err
		}

		if strings.HasPrefix(contentType, "application/json") {
			if err := json.Unmarshal(bodybytes, data.Interface()); err != nil {
				log.Error("json unmarshal: ", err)
				return reflect.Value{}, err
			}
		} else if t.Implements(_protomsg) {
			if err := proto.Unmarshal(bodybytes, data.Interface().(proto.Message)); err != nil {
				log.Error("proto unmarshal: ", err)
				return reflect.Value{}, err
			}
		}
	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(_maxFileSize); err != nil {
			log.Error("parse multipart/form-data: ", err)
			return reflect.Value{}, err
		}
		if err := _schema.Decode(data.Interface(), r.PostForm); err != nil {
			log.Error("schema decode for multipart/form-data: ", err)
			return reflect.Value{}, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			log.Error("parse form: ", err)
			return reflect.Value{}, err
		}
		var formdata url.Values
		if r.Method == "GET" {
			formdata = r.Form
		} else {
			formdata = r.PostForm
		}

		if err := _schema.Decode(data.Interface(), formdata); err != nil {
			log.Error("schema decode: ", err)
			return reflect.Value{}, err
		}
	}

	return data, nil
}

func _mkHandleFunc(f interface{}, params ...string) http.HandlerFunc {

	typ := reflect.TypeOf(f)
	if typ.Kind() != reflect.Func {
		panic("not func type")
	}

	if typ.NumIn() < 2 {
		panic("func must have two parameters as least")
	}

	if !typ.In(0).Implements(_respWriter) || typ.In(0).ConvertibleTo(_request) {
		panic("func must func(http.ResponseWrier, *http.Request, ...)")
	}

	if typ.NumIn() > 2+len(params) {
		if !_isPtrOfStruct(typ.In(typ.NumIn()-1)) && !_isByteSlice(typ.In(typ.NumIn()-1)) {
			panic("last parameter of func must be pointer of struct or slice of byte")
		}
	}

	if typ.NumOut() > 0 && !typ.Out(0).Implements(_result) {
		panic("func return type must be Result interface")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		vals := []reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)}
		for i, x := range params {
			tmp, ok := vars[x]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			val := _cast(tmp, typ.In(i+2).Kind())
			if !val.IsValid() {
				log.Errorf("%d: %v cast to %v error", i, x, typ.In(i).Kind())
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			vals = append(vals, val)
		}

		if typ.NumIn() > 2+len(params) {
			val, err := _readData(r, typ.In(typ.NumIn()-1))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if val.IsValid() {
				vals = append(vals, val)
			}
		}

		outs := reflect.ValueOf(f).Call(vals)

		if len(outs) > 0 && !outs[0].IsNil() {
			Display(w, outs[0].Interface().(Result))
		}
	})
}
