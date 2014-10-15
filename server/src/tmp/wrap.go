package msg

const (
	QLOGIN = 1
	RLOGIN = 1
)

func (m *MQLogin) GetType() uint32 {
	return QLOGIN
}

func (m *MRLogin) GetType() uint32 {
	return RLOGIN
}