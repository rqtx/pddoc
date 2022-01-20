package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Subsection struct {
	Name  string
	Table Table
}

type Section struct {
	Name        string
	Subsections []Subsection
}

type Document struct {
	Name     string
	Sections []Section
}

const BEGING_MARK string = "<!-- BEGIN_PDOCS -->"
const END_MARK string = "<!-- END_PDOCS -->"

func NewSubsection(name string, table Table) Subsection {
	return Subsection{
		Name:  name,
		Table: table,
	}
}

func NewSection(name string, subs []Subsection) Section {
	return Section{
		Name:        name,
		Subsections: subs,
	}
}

func NewDocumet(secs []Section) Document {
	return Document{
		Sections: secs,
	}
}

func (sec Section) WriteSection(fd *os.File) {
	fd.WriteString(fmt.Sprintf("## %s\n\n", sec.Name))
	for _, subs := range sec.Subsections {
		fd.WriteString(fmt.Sprintf("### %s\n\n", subs.Name))
		subs.Table.WriteToFile(fd)
	}
}

func (doc Document) WriteDocument(file string) (err error) {
	var fd *os.File

	if err = sanitize(file); err != nil {
		return
	}

	if fd, err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0); err != nil {
		return
	}
	defer func() {
		if cErr := fd.Close(); err == nil {
			err = cErr
		}
	}()

	fd.WriteString(fmt.Sprintf(BEGING_MARK + "\n"))
	for _, secs := range doc.Sections {
		secs.WriteSection(fd)
	}
	fd.WriteString(fmt.Sprintf(END_MARK))

	return
}

func sanitize(file string) (err error) {
	var fd *os.File
	var beginLine, endLine int
	if fd, err = os.OpenFile(file, os.O_RDWR, 0); err != nil {
		return
	}
	defer func() {
		// Check if last caracter is \n
		fd.Seek(-1, io.SeekEnd)
		b1 := make([]byte, 1)
		_, err := fd.Read(b1)
		if string(b1) != "\n" {
			fd.Seek(0, io.SeekEnd)
			fd.WriteString("\n")
		}
		// Close file
		if cErr := fd.Close(); err == nil {
			err = cErr
		}
	}()
	if beginLine, err = scanner(fd, BEGING_MARK); beginLine == 0 || err != nil {
		return
	}
	if endLine, err = lineCounter(fd); err != nil {
		return
	}
	if err = removeLines(fd, beginLine, endLine); err != nil {
		return
	}
	return
}

func scanner(input io.ReadSeeker, mark string) (int, error) {
	scanner := bufio.NewScanner(input) // Splits on newlines by default.
	line := 1

	defer input.Seek(0, io.SeekStart)

	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), mark) {
			return line, nil
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		return -1, err
	}
	return 0, nil
}

func removeLines(f *os.File, start, end int) (err error) {
	//https://www.rosettacode.org/wiki/Remove_lines_from_a_file#Go
	n := end - start + 1
	defer f.Seek(0, io.SeekStart)

	if start < 1 {
		return errors.New("invalid request.  line numbers start at 1.")
	}
	if n < 0 {
		return errors.New("invalid request.  negative number to remove.")
	}
	if n == 0 {
		return nil
	}

	var b []byte
	if b, err = ioutil.ReadAll(f); err != nil {
		return
	}

	cut, ok := skip(b, start-1)
	if !ok {
		return fmt.Errorf("less than %d lines", start)
	}

	tail, ok := skip(cut, n)
	if !ok {
		return fmt.Errorf("less than %d lines after line %d", n, start)
	}

	t := int64(len(b) - len(cut))
	if err = f.Truncate(t); err != nil {
		return
	}
	if len(tail) > 0 {
		_, err = f.WriteAt(tail, t)
	}

	return
}

func skip(b []byte, n int) ([]byte, bool) {
	//https://www.rosettacode.org/wiki/Remove_lines_from_a_file#Go
	for ; n > 0; n-- {
		if len(b) == 0 {
			return nil, false
		}
		x := bytes.IndexByte(b, '\n')
		if x < 0 {
			x = len(b)
		} else {
			x++
		}
		b = b[x:]
	}
	return b, true
}

func lineCounter(r io.ReadSeeker) (int, error) {
	buf := make([]byte, 32*1024)
	count := 1
	lineSep := []byte{'\n'}
	defer r.Seek(0, io.SeekStart)

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
