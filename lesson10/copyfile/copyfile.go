package copyfile

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

const bufferSize = 1024 * 1024 * 10

// Copy from <from> file from <offset> position to <to> file <limit> bytes
func Copy(from, to string, limit, offset int) error {
	if offset < 0 {
		return errors.New("offset of source file must be greater then or equal to zero")
	}

	if limit < 0 {
		return errors.New("limit must be greater then or equal to zero (to read all file)")
	}

	sourceFile, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("could not open source file: %s", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("could not create destination file: %s", err)
	}
	defer destFile.Close()

	sourceFileStat, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("could not get file stat for %s: %s", from, err)
	}

	if sourceFileStat.Size() <= int64(offset) {
		return fmt.Errorf("size of source file %s is %d and is less then or equal to start copy offset %d",
			from, sourceFileStat.Size(), offset)
	}

	if limit > 0 && sourceFileStat.Size() < int64(offset+limit) {
		return fmt.Errorf("size of source file %s is %d and is less then start copy offset + bytes to copy",
			from, sourceFileStat.Size())
	}

	_, err = sourceFile.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not set start offset for source file: %s", err)
	}

	var maxBytesToRead int64
	if limit == 0 {
		maxBytesToRead = sourceFileStat.Size() - int64(offset)
	} else {
		maxBytesToRead = int64(limit)
	}

	bar := pb.Full.Start64(maxBytesToRead)
	defer bar.Finish()

	teeFileReader := io.TeeReader(
		bar.NewProxyReader(io.LimitReader(sourceFile, maxBytesToRead)),
		destFile)

	buffer := make([]byte, bufferSize)
	for err != io.EOF {
		_, err = teeFileReader.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
	}

	return nil
}
