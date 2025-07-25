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
	fmt.Println("File length:", fileLength, "Block number:", blockNum)
	var hash []byte = make([]byte, 20*blockNum+20)
	var pieces chan string
	pieces = make(chan string, 1)
	var info chan bencodeInfo
	info = make(chan bencodeInfo, 1)
	var channel chan string
	channel = make(chan string)
	var ch chan UploadInfo
	ch = make(chan UploadInfo, blockNum+5)
	wg := new(sync.WaitGroup)
	for i := 0; i < blockNum; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			start := j * PieceSize
			end := start + PieceSize
			if end > fileLength {
				end = fileLength
			}
			putPieces(node, j, file[start:end], ch)
		}(i)
	}
	wg.Wait()
	//fmt.Println("All pieces uploaded, preparing torrent file...")
	go ToDotTorrentFile(inputPath, outputPath, pieces, info, channel)
	//fmt.Println("Torrent file creation in progress...")
	var flag bool
	flag = true
	count := 0
	for flag {
		select {
		case result := <-ch:
			count++
			index := result.index
			fmt.Print(index)
			copy(hash[index*20:index*20+20], result.hashes[:])
		default:
			flag = false
		}
	}
	fmt.Println("All pieces processed, preparing to save torrent file...")
	//传回去哈希后各个块的索引
	fmt.Print(count)
	fmt.Println(len(hash), "pieces to be saved in torrent file")
	pieces <- string(hash)
	content := <-channel
	//fmt.Println("Torrent file content:", content)
	ViewInfo := <-info
	//fmt.Println("Torrent file content:", content)
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
		fmt.Print("Failed to upload piece", index, "to DHT node.")
		return
	}
	fmt.Println("Piece", index, "uploaded successfully")
	ch <- UploadInfo{HashResult, index}
}

func download(inputPath, outputPath string, node *dhtNode) error {
	fmt.Print("Downloading file from DHT node...\n")
	fileInfo, err := Open(inputPath)
	if err != nil {
		return err
	}
	var downloadFile []byte
	downloadFile = make([]byte, fileInfo.Length)
	var blocknum int
	blocknum = len(fileInfo.PieceHashes)
	wg := new(sync.WaitGroup)
	var ch chan DownloadInfo
	ch = make(chan DownloadInfo, blocknum+5)
	for i := 0; i < blocknum; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			getPieces(node, j, fileInfo.PieceHashes[j], ch)
		}(i)
	}
	wg.Wait()
	//fmt.Println("All pieces downloaded, preparing to save file...")
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
				downloadfileName = outputPath + "/" + fileInfo.Name
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
		//fmt.Println("Failed to download piece", index, "from DHT node.")
		return
	}
	//fmt.Println(data)
	ch <- DownloadInfo{data: []byte(data), index: index}
}
