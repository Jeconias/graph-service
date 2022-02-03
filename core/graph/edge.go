package graph

type EdgeInfo struct {
	Url  string `json:"url,omitempty"`
	Date int64  `json:"date,omitempty"`
}

type Edge struct {
	From  *Vertex
	To    *Vertex
	Infos []EdgeInfo
}
