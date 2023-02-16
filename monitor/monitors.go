package monitor

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

type Monitor interface {
	FileInitialState() (bool, error)
	FileNowState() (bool, error)
	VerifyFileChange() (string, int, error)
	OutFileInfo() (FileInfo, FileInfo, bool)
}

type monitors struct {
	MonitorId string
	FilePath  string
	Observe   func(msg string, level int, err error)
	InitState FileInfo
	NowState  FileInfo
}

type FileInfo struct {
	os.FileInfo
	FileHash string
	FileAuth string
}

var fileNilInfo = FileInfo{}

func NewMonitor(filePath string, id string, on func(msg string, level int, err error)) Monitor {
	return &monitors{
		MonitorId: id,
		FilePath:  filePath,
		Observe:   on,
		InitState: FileInfo{},
		NowState:  FileInfo{},
	}
}

// FileInitialState Load the original state of the file
func (m *monitors) FileInitialState() (bool, error) {
	return func(info FileInfo, ok bool, err error) (bool, error) {
		m.InitState = info
		return ok, err
	}(getFileInfo(m.FilePath))
}

// FileNowState Load the current state of the file
func (m *monitors) FileNowState() (bool, error) {
	return func(info FileInfo, ok bool, err error) (bool, error) {
		m.NowState = info
		return ok, err
	}(getFileInfo(m.FilePath))
}

// getFileInfo Get file status information
func getFileInfo(filepath string) (FileInfo, bool, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return fileNilInfo, false, err
	}
	fileInfo, err := file.Stat()
	if fileInfo.IsDir() {
		return fileNilInfo, false, fmt.Errorf("this is a directory not a file")
	}
	if err != nil {
		return fileNilInfo, false, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	all, err := io.ReadAll(file)
	if err != nil {
		return fileNilInfo, false, err
	}
	return FileInfo{
		FileInfo: fileInfo,
		FileHash: fmt.Sprintf("%x", md5.Sum(all)),
		FileAuth: "",
	}, true, nil
}

// VerifyFileChange Verify that the file has changed
func (m *monitors) VerifyFileChange() (string, int, error) {
	ok, err := m.FileNowState()
	if err != nil {
		return "", -1, err
	}

	return func(msg string, l int, err error) (string, int, error) {
		defer m.Observe(msg, l, err)
		return msg, l, err
	}(verifyFileChange(m.InitState, m.NowState, ok))
}

func verifyFileChange(init FileInfo, now FileInfo, ok bool) (string, int, error) {
	if !ok || init.FileInfo.IsDir() || now.FileInfo.IsDir() {
		return "ERROR ", -1, nil
	}
	switch {
	case init.FileHash != now.FileHash:
		return "file content changed. <HASH changed>", 1, nil
	case init.FileInfo.ModTime() != now.FileInfo.ModTime():
		return "file Time changed. <TIME changed>", 1, nil
	case init.FileInfo.Mode() != now.FileInfo.Mode():
		return "file MODE changed. <MODE changed>", 1, nil
	default:
		return "NOT CHANGE.", 0, nil
	}
}

func (m *monitors) OutFileInfo() (FileInfo, FileInfo, bool) {
	return m.InitState, m.NowState, true
}
