package graph

type Vertex struct {
	key string
}

func (v *Vertex) Init(key string) {
	if len(v.key) != 0 {
		return
	}

	v.key = key
}

func (v *Vertex) GetKey() string {
	return v.key
}
