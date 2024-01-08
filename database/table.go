package database

type Table string

func (t Table) String() string {
	return string(t)
}
