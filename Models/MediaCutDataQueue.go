package Models

type MediaCutDataQueue struct {
	ChunkId int32 `json:"chunkId"`
	From float64 `json:"from"`
	To float64 `json:"to"`
	Position int32 `json:"position"`
} 
