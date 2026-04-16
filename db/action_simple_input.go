// action_simple_input.go — ActionSimpleInput option struct for InsertActionSimple.
package db

// ActionSimpleInput groups parameters for the InsertActionSimple convenience wrapper.
type ActionSimpleInput struct {
	Snapshot   string
	Detail     string
	BatchID    string
	MediaID    int64
	FileAction FileActionType
}
