// Program: crlf2lf
// Description: Change the real CRLF into LF
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

const sniffLen = 8192

var root string

func main() {
	flag.StringVar(&root, "root", ".", "root directory")
	flag.Parse()

	if root == "" {
		if len(os.Args) > 1 {
			root = os.Args[1]
		}
	}
	if root == "" {
		log.Fatalln("root directory required")
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		return normalizeFile(path)
	})

	if err != nil {
		log.Fatalln(err)
	}
}

func normalizeFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if !isTextFile(data) {
		// 跳过二进制文件
		log.Printf("skip binary file: %s\n", path)
		return nil
	}

	if !bytes.Contains(data, []byte("\r\n")) {
		// 跳过不包含 CRLF 的文件
		return nil
	}

	out := bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))

	return atomicWriteFile(path, out)
}

// isTextFile 通过检测前 8KB 是否含有 null byte 来判断是不是纯文本文件
// 这与 git 的二进制检测逻辑一致
func isTextFile(data []byte) bool {
	sniff := bytes.Clone(data)
	if len(sniff) > sniffLen {
		sniff = sniff[:sniffLen]
	}
	return !bytes.ContainsRune(sniff, 0)
}

// atomicWriteFile 以原子方式将数据写回目标文件：
//  1. 在同目录创建临时文件
//  2. Sync 确保落盘
//  3. 继承原文件权限
//  4. Rename 原子替换（POSIX 保证原子性）
//  5. defer 确保中途失败时清理临时文件
func atomicWriteFile(path string, data []byte) error {
	// 继承原文件权限
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	perm := info.Mode().Perm()

	// 创建临时文件
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".crlf_normalize_*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}

	// 确保失败时清理临时文件
	tmpName := tmp.Name()
	committed := false
	defer func() {
		if !committed {
			_ = os.Remove(tmpName)
		}
	}()

	// 写入数据
	_, err = tmp.Write(data)
	if err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write temp file: %w", err)
	}

	// 修正权限
	err = tmp.Chmod(perm)
	if err != nil {
		_ = tmp.Close()
		return fmt.Errorf("chmod temp file: %w", err)
	}

	// 确保落盘
	err = tmp.Sync()
	if err != nil {
		_ = tmp.Close()
		return fmt.Errorf("sync temp file: %w", err)
	}

	// 关闭临时文件
	err = tmp.Close()
	if err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	// Rename 原子替换
	err = os.Rename(tmpName, path)
	if err != nil {
		return fmt.Errorf("rename: %w", err)
	}

	committed = true
	return nil
}
