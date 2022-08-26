package types

// A file gets divided into multiple chunk
type Chunk struct {
	Id   string
	Data []byte
}

// A file gets distributed among the peers
type File struct {
	Name     string
	Chunks   map[string]Chunk
	Sequence []string // sequence of chunk ids
}

func NewFile(name string, chunks map[string]Chunk, sequence []string) *File {
	return &File{Name: name, Chunks: chunks, Sequence: sequence}
}

func (f *File) JoinChunks() []byte {
	var data []byte
	for _, id := range f.Sequence {
		data = append(data, f.Chunks[id].Data...)
	}

	return data
}
