package main


//import "com.abneptis.oss/aws"
//import "com.abneptis.oss/aws/sqs"
//import "http"
import "flag"
import "log"
import "fmt"

func main(){
  flag.Parse()
  q, err := GetQueue()
  if err != nil {
    log.Fatalf("Couldn't create queue: %v\n", err)
  }
  fmt.Printf("%v\n", q.Name)
}
