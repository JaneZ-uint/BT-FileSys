package main

import (
	"fmt"
	"os"
	"sync"
)

const PieceSize = 256 << 10

type UploadInfo struct {
	hashes [20]byte
	index  int
}

type DownloadInfo struct {
}

func upload(inputPath, outputPath string, node *dhtNode) error {
	file, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}
	fileLength := len(file)
	var blockNum int
	if fileLength%PieceSize == 0 {
		blockNum = fileLength / PieceSize
	} else {
		blockNum = fileLength/PieceSize + 1
	}

	var hash []byte = make([]byte, 20*blockNum)
	var pieces chan string
	var info chan bencodeInfo
	var channel chan string
	var ch chan UploadInfo
	var wg sync.WaitGroup
	for i := 0; i < blockNum; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := i * PieceSize
			end := start + PieceSize
			if end > fileLength {
				end = fileLength
			}
			putPieces(node, blockNum, i, file[start:end], ch)
		}(i)
	}
	wg.Wait()
	go ToDotTorrentFile(inputPath, outputPath, pieces, info, channel)
	var flag bool
	flag = true
	for flag {
		select {
		case result := <-ch:
			index := result.index
			copy(hash[index*20:index*20+20], result.hashes[:])
		default:
			flag = false
		}
	}
	//传回去哈希后各个块的索引
	pieces <- string(hash)
	content := <-channel
	ViewInfo := <-info
	key, err := ViewInfo.hash()
	if err != nil {
		return err
	}
	(*node).Put(fmt.Sprintf("%x", key), content)
	return nil
}

func putPieces(node *dhtNode, blocknum int, index int, data []byte, ch chan UploadInfo) {
	info := pieceInfo{data, index}
	HashResult, err := info.pieceHash()
	if err != nil {
		return
	}
	ok := (*node).Put(fmt.Sprintf("%x", HashResult), string(data))
	if !ok {
		return
	}
	ch <- UploadInfo{HashResult, index}
}

func download(inputPath, outputPath string, node *dhtNode) {

}
