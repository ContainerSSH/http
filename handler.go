package http

type RequestHandler interface {
	OnRequest(request ServerRequest, response ServerResponse) error
}

type ServerRequest interface {
	Decode(target interface{}) error
}

type ServerResponse interface {
	SetStatus(statusCode uint16)
	SetBody(interface{})
}


