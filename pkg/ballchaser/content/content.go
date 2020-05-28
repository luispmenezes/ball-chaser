package content

type Content struct {
	Length        uint
	CRC           uint
	Levels        []string
	KeyFrames     []KeyFrame
	NetworkStream []byte
	DebugStrings  []DebugString
	Ticks         []Tick
	Packages      []string
	Objects       []string
	Names         []string
	ClassIndex    map[string]uint
	NetClass      map[string]NetClass
}
