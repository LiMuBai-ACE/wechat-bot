package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// 获取当前运行的文件路径和文件名
func GetDirPath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("无法获取运行程序的信息")
	}

	// 获取当前文件所在目录
	dirPath := filepath.Dir(file)

	// 循环直到找到项目根路径（即存在 go.mod 的目录）
	for ; dirPath != filepath.Dir(dirPath); dirPath = filepath.Dir(dirPath) {
		if _, err := os.Stat(filepath.Join(dirPath, "go.mod")); err == nil {
			return dirPath
		}
	}

	panic("未能找到当前项目的根路径")
}
