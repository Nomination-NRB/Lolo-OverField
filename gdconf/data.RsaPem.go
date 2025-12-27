package gdconf

type RsaPem struct {
	PrivatePem []byte
}

func (g *GameConfig) loadRsaPem() {
	g.Data.RsaPem = &RsaPem{
		PrivatePem: make([]byte, 0),
	}
	ReadFile(&g.Data.RsaPem.PrivatePem, g.dataPath+"PrivateKey.pem")
}

func GetPrivatePem() []byte {
	return cc.Data.RsaPem.PrivatePem
}
