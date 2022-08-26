package utils

import (
	"github.com/ShauryaAg/opfs/types"
	"github.com/google/uuid"
)

func DivideIntoChunks(data []byte, chunkSize int) (map[string]types.Chunk, []string) {
	chunks := make(map[string]types.Chunk)
	var sequence []string

	for i := 0; i < len(data); i += chunkSize {
		chunk := new(types.Chunk)
		chunk.Id = uuid.New().String()
		chunk.Data = data[i : i+chunkSize]

		chunks[chunk.Id] = *chunk
		sequence = append(sequence, chunk.Id)
	}
	return chunks, sequence
}
