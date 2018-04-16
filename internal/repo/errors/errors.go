package errors

type NothingToUpdate struct{}

func (n NothingToUpdate) Error() string {
	return "no fields to update"
}
