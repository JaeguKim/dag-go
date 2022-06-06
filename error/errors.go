package error

type (
	ErrorType   int
	systemError struct {
		errorType ErrorType
		reason    string
	}
)

const (
	resourceErrorTypeIllegalName ErrorType = iota
	resourceErrorTypeIllegalType
	resourceErrorTypeUnsupportedType
	resourceErrorTypeIllegalScalar
	resourceErrorTypeIllegalRanges
	resourceErrorTypeIllegalSet
	resourceErrorTypeIllegalDisk
	resourceErrorTypeIllegalReservation
	resourceErrorTypeIllegalShare

	noReason = "" // make error generation code more readable
)

var (
	errorMessages = map[ErrorType]string{
		resourceErrorTypeIllegalName:        "missing or illegal resource name",
		resourceErrorTypeIllegalType:        "missing or illegal resource type",
		resourceErrorTypeUnsupportedType:    "unsupported resource type",
		resourceErrorTypeIllegalScalar:      "illegal scalar resource",
		resourceErrorTypeIllegalRanges:      "illegal ranges resource",
		resourceErrorTypeIllegalSet:         "illegal set resource",
		resourceErrorTypeIllegalDisk:        "illegal disk resource",
		resourceErrorTypeIllegalReservation: "illegal resource reservation",
		resourceErrorTypeIllegalShare:       "illegal shared resource",
	}
)

func (t ErrorType) Generate(reason string) error {
	msg := errorMessages[t]
	if reason != noReason {
		if msg != "" {
			msg += ": " + reason
		} else {
			msg = reason
		}
	}
	return &systemError{errorType: t, reason: msg}
}

func (err *systemError) Error() string {
	if err.reason != "" {
		return "resource error: " + err.reason
	}
	return "resource error"
}

