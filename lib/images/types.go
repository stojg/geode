package images

import "fmt"

const (
	readError = iota
	writeError
	formatError
	memoryError
)

func newDecoderErr(code int, text string, args ...interface{}) error {
	return &DecoderError{
		Type: code,
		Err:  fmt.Errorf(text, args...),
	}
}

type DecoderError struct {
	Type int
	Err  error
}

func (d *DecoderError) Error() string {
	switch d.Type {
	case readError:
		return fmt.Sprintf("RGBE decoder read error: %s", d.Err.Error())
	case writeError:
		return fmt.Sprintf("RGBE decoder write error: %s", d.Err.Error())
	case formatError:
		return fmt.Sprintf("RGBE decoder format error: %s", d.Err.Error())
	case memoryError:
		fallthrough
	default:
		return fmt.Sprintf("RGBE decoder error: %s", d.Err.Error())
	}
}

func (d *DecoderError) Is(target error) bool {
	t, ok := target.(*DecoderError)
	if !ok {
		return false
	}
	return t.Type == d.Type && t.Err.Error() == d.Err.Error()
}
