package nuomi

import (
  "context"
  "fmt"
  "os"
  "strings"
  "sync"
  "time"
)

//糯米开始以各种姿势翻滚

func (n *NuoMi) Errors() <- chan error{
  return n.errorChan
}

func (n *NuoMi) Results() <- chan Result{
  return n.resultChan
}

func (n *NuoMi) incrementRequests() {
  n.mu.Lock()
  n.requestsIssued++
  n.mu.Unlock()
}

func (n *NuoMi) roll()error{
  defer close(n.resultChan)
  defer close(n.errorChan)

  if err := n.plugin.PreRun();err != nil{
    return err
  }

  var workerGroup sync.WaitGroup
  workerGroup.Add(n.Opts.Coroutines)

  wordChan := make(chan string, n.Opts.Coroutines)
  for i := 0; i < n.Opts.Coroutines; i++{
    go n.explodeWorker(wordChan, &workerGroup)
  }

  scanner, err := n.getWordlist()
  if err != nil{
    return err
  }
Scan:
  for scanner.Scan(){
    select{
    case <- n.context.Done():
      break Scan
      case wordChan <- scanner.Text():
    }
  }
  close(wordChan)
  workerGroup.Wait()
  return nil
}



// 爆破工作者
func (n *NuoMi) explodeWorker(wordChan <- chan string, wg *sync.WaitGroup) {
  defer wg.Done()
  for{
    select {
    case <- n.context.Done():
      return
    case word, ok := <- wordChan:
      if !ok{
        return
      }
      n.incrementRequests()
      wordCleaned := strings.TrimSpace(word)
      if strings.HasPrefix(wordCleaned,"#") || len(wordCleaned) == 0{
        break
      }

      // 将读取到的内容发给对应字段
      res, err := n.plugin.Run(wordCleaned)
      if err != nil{
        n.errorChan <- err
        continue
      }else {
        //获取到结果
        for _, r := range res{
          n.resultChan <- r
        }
      }

      select {
      //读取单行数据时需要结束或者延时
      case <-n.context.Done():

        case <- time.After(n.Opts.Delay):

      }
    }
  }
}

func errorWorker(n *NuoMi, wg *sync.WaitGroup){
  defer wg.Done()
  //todo:这里目前就简单地打印出来即可
  for e := range n.Errors(){
    fmt.Println(e)
  }
}

//todo:目前还是filename，将结果保存在文件中，要入库
func resultWorker(n *NuoMi, filename string, wg *sync.WaitGroup){
  defer wg.Done()
  var f *os.File
  var err error
  if filename != ""{
    f, err =os.Create(filename)
    if err != nil{
      n.LogError.Fatalf("error on creating output file: %v", err)
    }
    defer f.Close()
  }

  for r := range n.Results(){
    s, err := r.ToString(n)
    if err != nil{
      n.LogError.Fatal(err)
    }
    if s != ""{
      //todo:入库操作
    }
  }

}

func progressWorker(c context.Context, n *NuoMi, wg *sync.WaitGroup){
  defer wg.Done()
  tick := time.NewTicker(1 * time.Second)

  for {
    select {
    case <-tick.C:
      fmt.Println("Hello World?")
    case <-c.Done():
      return
    }
  }
}
