package torrent

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/jackpal/bencode-go"
)

// 默认参数
type Config struct {
	PieceLength int //256KB
	AnnounceURL string
}

func DefaultConfig() Config {
	return Config{
		PieceLength: 256 << 10, // 256KB
		AnnounceURL: "",
	}
}

// .torrent 文件信息
type MetaInfo struct {
	Announce string   `bencode:"announce"`
	Info     InfoDict `bencode:"info"`
}

type InfoDict struct {
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
	Length      int    `bencode:"length,omitempty"`
	Name        string `bencode:"name"`
	Files       []File `bencode:"files,omitempty"`
}

type File struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}

// 一种人类能读懂的torrent文件
type Torrent struct {
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	TotalSize   int64
	Name        string
	IsMultiFile bool
	Files       []File
}

func ParseFile(path string) (*Torrent, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ReadFile(file)
}

// 解析Bencode编码文件
func ReadFile(r io.Reader) (*Torrent, error) {
	var meta MetaInfo
	err := bencode.Unmarshal(r, &meta)
	if err != nil {
		return nil, err
	}
	return processMetaInfo(&meta)
}

func Create(inputPath, outputPath string, config Config) error {
	info, err := buildInfoDict(inputPath, config.PieceLength)
	if err != nil {
		return err
	}
	meta := MetaInfo{
		Announce: config.AnnounceURL,
		Info:     *info,
	}

	if outputPath == "" {
		outputPath = inputPath + ".torrent"
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return bencode.Marshal(file, meta)
}

func processMetaInfo(meta *MetaInfo) (*Torrent, error) {
	infoHash, err := hashInfoDict(&meta.Info)
	if err != nil {
		return nil, err
	}
	pieceHashes, err := splitPieceHashes(meta.Info.Pieces)
	if err != nil {
		return nil, err
	}
	var totalSize int64
	if meta.Info.Length > 0 {
		totalSize = int64(meta.Info.Length)
	} else {
		for _, f := range meta.Info.Files {
			totalSize += int64(f.Length)
		}
	}

	return &Torrent{
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: meta.Info.PieceLength,
		TotalSize:   totalSize,
		Name:        meta.Info.Name,
		IsMultiFile: meta.Info.Length == 0,
		Files:       meta.Info.Files,
	}, nil
}

func hashInfoDict(info *InfoDict) ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, info)
	if err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(buf.Bytes()), nil
}

func splitPieceHashes(pieces string) ([][20]byte, error) {
	const hashSize = 20
	buf := []byte(pieces)
	if len(buf)%hashSize != 0 {
		return nil, errors.New("")
	}
	count := len(buf) / hashSize
	hashes := make([][20]byte, count)
	for i := 0; i < count; i++ {
		copy(hashes[i][:], buf[i*hashSize:(i+1)*hashSize])
	}
	return hashes, nil
}

func buildInfoDict(path string, pieceLength int) (*InfoDict, error) {
	info := &InfoDict{
		PieceLength: pieceLength,
		Name:        filepath.Base(path),
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, errors.New("")
	}
	info.Length = int(fileInfo.Size())
	pieces, err := hashPieces(path, pieceLength)
	if err != nil {
		return nil, err
	}
	info.Pieces = pieces
	return info, nil
}

func hashPieces(path string, pieceLength int) (string, error) {
	var pieces = make([]byte, pieceLength)
	var buf bytes.Buffer
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	var Lock sync.RWMutex
	var wg sync.WaitGroup
	var errChan = make(chan error, 1)
	for {
		n, err := io.ReadFull(file, pieces)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		wg.Add(1)
		go func(data []byte) {
			defer wg.Done()
			h := sha1.Sum(data)
			Lock.Lock()
			defer Lock.Unlock()
			_, err := buf.Write(h[:])
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}(pieces[:n])
	}
	wg.Wait()
	select {
	case err := <-errChan:
		return "", err
	default:
	}
	return buf.String(), nil
}

func (t *Torrent) Print() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Name: %s\n", t.Name)
	fmt.Fprintf(&buf, "InfoHash: %s\n", hex.EncodeToString(t.InfoHash[:]))
	fmt.Fprintf(&buf, "Piece Length: %d bytes\n", t.PieceLength)
	fmt.Fprintf(&buf, "Total Size: %d bytes\n", t.TotalSize)
	fmt.Fprintf(&buf, "Pieces: %d\n", len(t.PieceHashes))
	if t.IsMultiFile {
		fmt.Fprintln(&buf, "\nFiles:")
		for _, f := range t.Files {
			fmt.Fprintf(&buf, "  %s - %d bytes\n", filepath.Join(f.Path...), f.Length)
		}
	}
	return buf.String()
}
