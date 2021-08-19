package services

import (
	// Standard Library Imports
	"net/url"
	"strings"
)

type Service string

const (
	Blob  Service = "b"
	Queue Service = "q"
	Table Service = "t"
	File  Service = "f"
)

const paramKey = "ss"

type Services []Service

func (s *Services) SetParam(params *url.Values) {
	if s != nil && len(*s) > 0 {
		ss := ""
		for _, service := range *s {
			ss += string(service)
		}

		params.Add(paramKey, ss)
	}
}

func (s Services) GetParam() (services string) {
	if len(s) > 0 {
		values := &url.Values{}
		s.SetParam(values)

		services = values.Encode()
	}

	return
}

func (s Services) GetURLDecodedParam() (services string) {
	if len(s) > 0 {
		services, _ = url.QueryUnescape(s.GetParam())
	}

	return
}

func Parse(services string) Services {
	vMap := map[Service]struct{}{
		Blob:  {},
		Queue: {},
		Table: {},
		File:  {},
	}

	var ss Services
	splitServices := strings.Split(services, "")
	for _, service := range splitServices {
		check := Service(service)
		if _, ok := vMap[check]; ok {
			ss = append(ss, check)
		}
	}

	return ss
}
