package lib

// GinLogger logger for gin framework [subbed from main logger]
type GinLogger struct {
	*Logger
}

// Writer interface implementation for gin-framework
func (l GinLogger) Write(p []byte) (n int, err error) {
	str := string(p)
	size := len(p)
	l.Info(str)
	return size, nil
}
