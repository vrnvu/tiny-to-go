package tiny

type Service interface {
	Find(code string) (*Redirect, error)
	Save(redirect *Redirect) error
}
