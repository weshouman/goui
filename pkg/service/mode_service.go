package service

type ModeService struct {
	current int
	modeReg *ModeRegistry
}

func NewModeService(reg *ModeRegistry) *ModeService {
	return &ModeService{modeReg: reg}
}

func (m *ModeService) SetMode(v int) {
	m.current = v
}

func (m *ModeService) CurrentMode() int {
	return m.current
}