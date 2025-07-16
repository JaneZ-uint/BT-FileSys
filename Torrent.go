// Reference: https://github.com/veggiedefender/torrent-client/blob/master/torrentfile/torrentfile.go

package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
)

const PieceSizes = 256 << 10

// TorrentFile encodes the metadata from a .torrent file
type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

type pieceInfo struct {
	Data  []byte
	index int
}

func Open(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	defer file.Close()

	bto := bencodeTorrent{}
	err = bencode.Unmarshal(file, &bto)
	if err != nil {
		return TorrentFile{}, err
	}
	return bto.toTorrentFile()
}

func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20 // Length of SHA-1 hash
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("invalid pieces length: %d, must be a multiple of %d", len(buf), hashLen)
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (bto *bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}
	pieceHashes, err := bto.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, err
	}
	t := TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}
	return t, nil
}

func (i *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (tmp *pieceInfo) pieceHash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *tmp)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func ToDotTorrentFile(inputPath, outputPath string, pieces chan string, info chan bencodeInfo, ch chan string) error {
	file, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	tmp := bencodeTorrent{}
	tmp.Announce = "Tracker"
	tmp.Info = bencodeInfo{"", PieceSizes, int(file.Size()), file.Name()}
	tmp.Info.Pieces = <-pieces
	if outputPath == "" {
		outputPath = file.Name() + ".torrent"
	} else {
		outputPath = outputPath + "/" + file.Name() + ".torrent"
	}
	newfile, err1 := os.Create(outputPath)
	fmt.Println("Torrent file created:", outputPath)
	if err1 != nil {
		fmt.Println("Error creating torrent file:", err1)
		return err1
	}
	//fmt.Println("Torrent file created:", outputPath)
	err2 := bencode.Marshal(newfile, tmp)
	//fmt.Println("Torrent file created:", outputPath)
	if err2 != nil {
		return err2
	}
	content, err3 := os.ReadFile(outputPath)
	if err3 != nil {
		return err3
	}
	info <- tmp.Info
	ch <- string(content)
	return nil
}
