package middleware

type middleware struct {
	config Config
}

type Middleware interface {
}

type Config struct {
}

func New(config Config) middleware {
	return middleware{
		config: config,
	}
}
