package types

const (
	DOWNLOADAWAIT = iota //default state
	DOWNLOADING
	DOWNLOADPAUSE
	DOWNLOADFINISH
	DOWNLOADERROR
	DOWNLOADFINISHED
)
