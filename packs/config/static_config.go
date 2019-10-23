package config

import "time"

const (
	HTTP_READTIMEOUT		=	10 * time.Second
	HTTP_WRETETIMEOUT		=	10 * time.Second
	HTTP_MAXHEADERSBYTES	=	1 << 20
)
