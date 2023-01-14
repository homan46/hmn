package log

import "go.uber.org/zap"

var ZLog *zap.Logger

func init() {
	var err error

	//ZLog, err = zap.NewProduction()
	ZLog, err = zap.NewDevelopment()

	if err != nil {
		panic("fail to create logger")
	}
}
