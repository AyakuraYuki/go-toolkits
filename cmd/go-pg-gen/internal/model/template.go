package model

type ExistMethod struct {
	ID int64
}

func (m *ExistMethod) IsExist() bool { return m != nil && m.ID > 0 }
