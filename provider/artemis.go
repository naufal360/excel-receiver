package provider

import (
	"excel-receiver/config"
	"fmt"

	"github.com/go-stomp/stomp/v3"
)

func NewArtemis() (conn *stomp.Conn, err error) {
	cfg := config.Configuration.Artemis

	conn, err = stomp.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		stomp.ConnOpt.Login(cfg.Username, cfg.Password),
	)

	if err != nil {
		return nil, err
	}

	return conn, err
}
