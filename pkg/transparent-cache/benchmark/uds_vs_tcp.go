//go:build ignore

package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"text/tabwriter"
	"time"
)

type Metrics struct {
	rps        float64
	avgLatency time.Duration
	mbps       float64
}

type BenchmarkCase struct {
	name       string
	payload    int
	conns      int
	total      int
	downstream time.Duration
}

type ExperimentResult struct {
	name       string
	payload    int
	conns      int
	total      int
	downstream time.Duration
	tcp        Metrics
	uds        Metrics
}

func main() {
	fmt.Println("======================================")
	fmt.Println(" Upload Pipeline Benchmark (TCP vs UDS)")
	fmt.Println("======================================\n")

	tests := []BenchmarkCase{
		{"Small upload", 1024, 20, 100, 200 * time.Microsecond},
		{"Medium upload", 4096, 20, 100, 200 * time.Microsecond},
		{"Large upload", 16384, 20, 100, 200 * time.Microsecond},
		{"Very large upload", 65536, 20, 100, 200 * time.Microsecond},
	}

	results := make([]ExperimentResult, 0, len(tests))
	for _, t := range tests {
		results = append(results, runExperiment(t))
	}

	printResultsTable(results)
}

func runExperiment(test BenchmarkCase) ExperimentResult {
	tcpLn, _ := startServer("tcp", "127.0.0.1:0", test.downstream)
	defer tcpLn.Close()
	tcpAddr := tcpLn.Addr().String()

	udsAddr := fmt.Sprintf("/tmp/upload_%d.sock", time.Now().UnixNano())
	udsLn, _ := startServer("unix", udsAddr, test.downstream)
	defer func() {
		udsLn.Close()
		os.Remove(udsAddr)
	}()

	tcp := runClient("tcp", tcpAddr, test.conns, test.total, test.payload)
	uds := runClient("unix", udsAddr, test.conns, test.total, test.payload)

	return ExperimentResult{
		name:       test.name,
		payload:    test.payload,
		conns:      test.conns,
		total:      test.total,
		downstream: test.downstream,
		tcp:        tcp,
		uds:        uds,
	}
}

func startServer(network, addr string, downstream time.Duration) (net.Listener, error) {
	if network == "unix" {
		os.Remove(addr)
	}

	ln, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go handleConn(conn, downstream)
		}
	}()

	return ln, nil
}

func handleConn(conn net.Conn, downstream time.Duration) {
	defer conn.Close()

	lenBuf := make([]byte, 4)

	for {
		// Read payload size
		_, err := io.ReadFull(conn, lenBuf)
		if err != nil {
			return
		}
		size := binary.BigEndian.Uint32(lenBuf)

		// Read payload
		payload := make([]byte, size)
		_, err = io.ReadFull(conn, payload)
		if err != nil {
			return
		}

		// Simulate downstream (S3 upload)
		if downstream > 0 {
			time.Sleep(downstream)
		}

		// Return small metadata response (8 bytes)
		resp := make([]byte, 8)
		binary.BigEndian.PutUint64(resp, uint64(size)) // pretend "uploaded size"
		conn.Write(resp)
	}
}

func runClient(network, addr string, conns, total, payload int) Metrics {
	var wg sync.WaitGroup
	var totalLatency int64

	perConn := total / conns
	start := time.Now()

	wg.Add(conns)

	for i := 0; i < conns; i++ {
		go func() {
			defer wg.Done()

			conn, err := net.Dial(network, addr)
			if err != nil {
				return
			}
			defer conn.Close()

			data := make([]byte, payload)
			lenBuf := make([]byte, 4)
			resp := make([]byte, 8)

			binary.BigEndian.PutUint32(lenBuf, uint32(payload))

			for j := 0; j < perConn; j++ {
				t1 := time.Now()

				// send size + payload
				conn.Write(lenBuf)
				conn.Write(data)

				// read metadata response
				io.ReadFull(conn, resp)

				lat := time.Since(t1)
				atomic.AddInt64(&totalLatency, lat.Nanoseconds())
			}
		}()
	}

	wg.Wait()

	elapsed := time.Since(start)
	rps := float64(total) / elapsed.Seconds()

	totalBytes := float64(total * payload)
	mbps := totalBytes / elapsed.Seconds() / (1024 * 1024)

	avgLat := time.Duration(totalLatency / int64(total))

	return Metrics{
		rps:        rps,
		avgLatency: avgLat,
		mbps:       mbps,
	}
}

func printResultsTable(results []ExperimentResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "Test\tPayload\tTCP req/s\tUDS req/s\tThroughput benefit\tTCP latency\tUDS latency\tLatency benefit\tTCP MB/s\tUDS MB/s")
	for _, result := range results {
		fmt.Fprintf(
			w,
			"%s\t%dB\t%.0f\t%.0f\t%s\t%s\t%s\t%s\t%.2f\t%.2f\n",
			result.name,
			result.payload,
			result.tcp.rps,
			result.uds.rps,
			formatPercent(throughputBenefit(result.tcp, result.uds)),
			result.tcp.avgLatency,
			result.uds.avgLatency,
			formatPercent(latencyBenefit(result.tcp, result.uds)),
			result.tcp.mbps,
			result.uds.mbps,
		)
	}

	w.Flush()

	fmt.Println("\nPositive throughput benefit = UDS higher req/s. Positive latency benefit = UDS lower latency.")
}

func throughputBenefit(tcp, uds Metrics) float64 {
	if tcp.rps == 0 {
		return 0
	}

	return ((uds.rps - tcp.rps) / tcp.rps) * 100
}

func latencyBenefit(tcp, uds Metrics) float64 {
	tcpLatency := tcp.avgLatency.Nanoseconds()
	if tcpLatency == 0 {
		return 0
	}

	return (float64(tcpLatency-uds.avgLatency.Nanoseconds()) / float64(tcpLatency)) * 100
}

func formatPercent(value float64) string {
	return fmt.Sprintf("%+.2f%%", value)
}

// low payload, better UDS benefit
// if we increase the total requests, the benefits shrink. 


// Run 1
// ======================================
//  Upload Pipeline Benchmark (TCP vs UDS)
// ======================================

// Test               Payload  TCP req/s  UDS req/s  Throughput benefit  TCP latency  UDS latency  Latency benefit  TCP MB/s  UDS MB/s
// Small upload       1024B    28101      44888      +59.74%             523.306µs    382.372µs    +26.93%          27.44     43.84
// Medium upload      4096B    35191      68695      +95.21%             425.182µs    246.319µs    +42.07%          137.46    268.34
// Large upload       16384B   26566      40894      +53.93%             623.251µs    391.412µs    +37.20%          415.09    638.96
// Very large upload  65536B   20199      18365      -9.08%              818.659µs    957.803µs    -17.00%          1262.45   1147.81

// Run 2
// ======================================
//  Upload Pipeline Benchmark (TCP vs UDS)
// ======================================

// Test               Payload  TCP req/s  UDS req/s  Throughput benefit  TCP latency  UDS latency  Latency benefit  TCP MB/s  UDS MB/s
// Small upload       1024B    29559      56540      +91.28%             522.401µs    291.325µs    +44.23%          28.87     55.21
// Medium upload      4096B    28551      53164      +86.21%             560.567µs    301.608µs    +46.20%          111.53    207.67
// Large upload       16384B   28075      35612      +26.84%             570.949µs    406.074µs    +28.88%          438.68    556.44
// Very large upload  65536B   21731      20282      -6.67%              776.397µs    795.588µs    -2.47%           1358.19   1267.64