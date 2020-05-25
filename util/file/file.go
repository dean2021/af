package file

import (
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"os"
)

// 获取路径是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 获取文件md5 hash
// 仅适合小文件
//
func Md5File(path string) (string, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return "", err
	}

	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// 获取文件md5 hash
// 分块计算
func Md5FileWithChunk(path string) (string, error) {

	const filechunk = 8192 // we settle for 8KB

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	// calculate the file size
	info, err := file.Stat()
	if err != nil {
		return "", err
	}

	fileSize := info.Size()
	blocks := uint64(math.Ceil(float64(fileSize) / float64(filechunk)))
	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blockSize := int(math.Min(filechunk, float64(fileSize-int64(i*filechunk))))
		buf := make([]byte, blockSize)
		_, err = file.Read(buf)
		if err != nil {
			return "", err
		}
		_, err = io.WriteString(hash, string(buf)) // append into the hash
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// 拷贝文件
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

// 替换文件
func ReplaceFile(dstName, srcName string) (written int64, err error) {
	err = os.Remove(dstName)
	if err != nil {
		return
	}
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}
