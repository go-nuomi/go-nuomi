package nuomi

import (
  "bufio"
  "bytes"
  "fmt"
  "io"
  "os"
)

func (n *NuoMi) getWordlist()(*bufio.Scanner, error){
  wordlist, err := os.Open(n.Opts.Wordlist)
  if err != nil{
    return  nil, fmt.Errorf("failed to open wordlist: %v",err)
  }

  lines, err := lineCounter(wordlist)
  if err != nil{
    return nil, fmt.Errorf("failed to get number of lines: %v", err)
  }

  n.requestsExpected = lines
  n.requestsIssued = 0

  _, err = wordlist.Seek(0,0)
  if err != nil{
    return nil, fmt.Errorf("failed to rewind wordlist: %v", err)
  }
  return bufio.NewScanner(wordlist),nil
}

func lineCounter(r io.Reader)(int, error){
  buf := make([]byte, 32 * 1024)
  count := 1
  lineSep := []byte{'\n'}

  for{
    c, err := r.Read(buf)
    count += bytes.Count(buf[:c],lineSep)

    switch {
    case err == io.EOF:
      return count,nil
    case err != nil:
      return count,err
    }
  }

  }
