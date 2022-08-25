package utils

import (
	"github.com/ShauryaAg/opfs/types"
	"github.com/google/uuid"
)

func DivideIntoChunks(data []byte, chunkSize int) []types.Chunk {
	var chunks []types.Chunk
	for i := 0; i < len(data); i += chunkSize {
		chunk := new(types.Chunk)
		chunk.Id = uuid.New().String()
		chunk.Data = data[i : i+chunkSize]

		chunks = append(chunks, *chunk)
	}
	return chunks
}
