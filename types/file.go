package types

// A file gets divided into multiple chunk
type Chunk struct {
	Id   string
	Data []byte
}

// A file gets distributed among the peers
type File struct {
	Name   string
	Chunks []Chunk
}

func NewFile(name string, chunks []Chunk) *File {
	return &File{Name: name, Chunks: chunks}
}

func (f *File) JoinChunks() []byte {
	var data []byte
	for _, chunk := range f.Chunks {
		data = append(data, chunk.Data...)
	}
	return data
}
