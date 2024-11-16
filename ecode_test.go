package ecode

import (
	"errors"
	"log"
	"net/http"
	"testing"
)

var (
	ParameterErr = NewErrNo(1000400, "request param error")
)

func Test(t *testing.T) {
	e := FromError(ParameterErr)
	log.Println(e.Error()) // error: code = 1000400 msg = equest param error metadata = map[] cause = <nil>
	log.Println(e.Code())  // 1000400
	log.Println(e.Msg())   // request param error
	log.Println("============================")
}

func TestWithSuccess(t *testing.T) {
	Success = NewErrNo(1, "success")

	e2 := FromError(nil)
	log.Println(e2.Error()) // error: code = 1 msg = success metadata = map[] cause = <nil>
	log.Println(e2.Code())  // 1
	log.Println(e2.Msg())   // success
	log.Println("============================")
}

func TestWithMetadata(t *testing.T) {
	sms := NewErrNo(10000, "中国电信").WithMetadata(map[string]string{
		"name": "jerry",
	})
	log.Println(sms.Error())  // error: code = 10000 msg = 中国电信 metadata = map[name:jerry] cause = <nil>
	log.Println(sms.Code())   // 10000
	log.Println(sms.Msg())    // 中国电信
	log.Println(sms.Metadata) // map[name:jerry]
	log.Println("============================")
}

func TestWithCause(t *testing.T) {
	mms := NewErrNo(10086, "中国移动").WithCause(errors.New("我是原因"))
	log.Println(mms.Error())  // error: code = 10086 msg = 中国移动 metadata = map[] cause = 我是原因
	log.Println(mms.Code())   // 10086
	log.Println(mms.Msg())    // 中国电信
	log.Println(mms.Unwrap()) // 我是原因
	log.Println("============================")
}

func TestIs(t *testing.T) {
	tests := []struct {
		name string
		e    *ErrorNo
		err  error
		want bool
	}{
		{
			name: "true",
			e:    NewErrNo(404, ""),
			err:  NewErrNo(http.StatusNotFound, ""),
			want: true,
		},
		{
			name: "false",
			e:    NewErrNo(0, ""),
			err:  errors.New("test"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ok := tt.e.Is(tt.err); ok != tt.want {
				t.Errorf("ErrorNo.ErrorNo() = %v, want %v", ok, tt.want)
			}
		})
	}
}
