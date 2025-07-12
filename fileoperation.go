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
	data  []byte
	index int
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
			putPieces(node, i, file[start:end], ch)
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
	fmt.Println("Upload completed, torrent file saved to:", content)
	return nil
}

func putPieces(node *dhtNode, index int, data []byte, ch chan UploadInfo) {
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

func download(inputPath, outputPath string, node *dhtNode) error {
	fileInfo, err := Open(inputPath)
	if err != nil {
		return err
	}
	var downloadFile []byte
	downloadFile = make([]byte, fileInfo.Length)
	var blocknum int
	blocknum = len(fileInfo.PieceHashes)
	var wg sync.WaitGroup
	var ch chan DownloadInfo
	ch = make(chan DownloadInfo, blocknum+5)
	for i := 0; i < blocknum; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			getPieces(node, i, fileInfo.PieceHashes[i], ch)
		}(i)
	}
	wg.Wait()
	var flag bool
	flag = true
	for flag {
		select {
		case result := <-ch:
			index := result.index
			copy(downloadFile[index*PieceSize:index*PieceSize+len(result.data)], result.data)
		default:
			flag = false
			var downloadfileName string
			if outputPath == "" {
				downloadfileName = fileInfo.Name
			} else {
				downloadfileName = fileInfo.Name + "/" + outputPath
			}
			err = os.WriteFile(downloadfileName, downloadFile, 0644)
			if err != nil {
				return err
			}
			fmt.Println("Download completed, file saved to:", downloadfileName)
		}
	}
	return nil
}

func getPieces(node *dhtNode, index int, hash [20]byte, ch chan DownloadInfo) {
	ok, data := (*node).Get(fmt.Sprintf("%x", hash))
	if !ok {
		return
	}
	ch <- DownloadInfo{data: []byte(data), index: index}
}
