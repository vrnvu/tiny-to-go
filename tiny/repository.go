package tiny

type Repository interface {
	Find(code string) (*Redirect, error)
	Save(redirect *Redirect) error
}
