package usecase

type StatusInteractor struct {
}

func NewStatusInteractor() *StatusInteractor {
	return &StatusInteractor{}
}

func (i *StatusInteractor) Status() error {
	return nil
}
