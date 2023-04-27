package main

type config struct {
	Env             string `required:"true" default:"development" desc:"production, development"`
	DSN             string `required:"true" desc:"DSN для соединения с базой данных"`
	BindAddr        string `required:"true" default:":8080" split_words:"true" desc:"Адрес и порт входящих соединений"`
	ReadTimeout     int    `required:"true" default:"10" split_words:"true" desc:"Таймаут на чтение запроса"`
	WriteTimeout    int    `required:"true" default:"10" split_words:"true" desc:"Таймаут на запись ответа"`
	ShutdownTimeout int    `required:"true" default:"30" split_words:"true" desc:"Время до принудительного завершения сервиса после получения сигнала выхода (s)"`
}
