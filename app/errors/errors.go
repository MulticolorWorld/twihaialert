package errors

type AccountAlreadyExistError struct {
}

func (err AccountAlreadyExistError) Error() string {
	return "このアカウントは既に存在します"
}
