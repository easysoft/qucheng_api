package middleware

type Manager struct {
}

func New() *Manager {
	return &Manager{}
}

func (m *Manager) Mysql() *MysqlManager {
	return &MysqlManager{}
}
