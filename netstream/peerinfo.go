package netstream


/*
TODO:
1. 配置改成读文件
2. 将main里读配置删掉，移到这里
 */

var PeerTag string

const (
	PeerLocalProxy = "LocalProxy"
	PeerRemoteProxy = "RemoteProxy"
	PeerTCPClient = "TCPClient"
	PeerTCPServer = "TCPServer"
	PeerUDPClient = "UDPClient"
	PeerUDPServer = "UDPServer"
)

//func main() {
//    inputFile, inputError := os.Open("address.txt")
//    if inputError != nil {
//        fmt.Printf("Cannot open file", inputError.Error())
//        return
//    }
//    defer inputFile.Close()
//    inputReader := bufio.NewReader(inputFile)
//    i := 0
//    for {
//        inputString, readerError := inputReader.ReadString('\n')
//        if readerError == io.EOF {
//            return
//        }
//        i++
//        fmt.Printf("IP Address:%s", i, inputString)
//    }
//}

//func get_