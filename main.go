package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	checkInterval = 2 * time.Second
	bytesInMB     = 1024 * 1024
)

func main() {
	totalMemoryBytes, err := totalMemory()
	if err != nil {
		fmt.Println("get total memory: %w", err)
		return
	}

	maxMemoryBytes := maxMemory()

	fmt.Printf("total memory:\t%10.dMB\n", totalMemoryBytes/bytesInMB)
	fmt.Printf("max memory:\t%10.dMB\n", maxMemoryBytes/bytesInMB)

	termSigC := waitForTermination()
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-termSigC:
			fmt.Print("\nmemkill terminated")
			return

		case <-ticker.C:
			pids, err := findProcessesOverLimit(maxMemoryBytes, totalMemoryBytes)
			if err != nil {
				fmt.Println("Error finding processes:", err)
				continue
			}

			if len(pids) == 0 {
				fmt.Print(".")
				continue
			}

			fmt.Printf("\nterminating pids: %v\n", pids)
			for _, pid := range pids {
				err := terminateProcess(pid)
				if err != nil {
					fmt.Printf("Error terminating process %d: %v\n", pid, err)
				}
			}
		}
	}
}

func findProcessesOverLimit(limit, totalMemoryBytes int64) ([]int, error) {
	cmd := exec.Command("ps", "-eo", "pid,%mem")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var pids []int

	lines := string(output)
	for _, line := range strings.Split(lines, "\n")[1:] {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}

		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}

		memoryPercentage, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return nil, err
		}
		if memoryPercentage < 1 {
			continue
		}
		memoryPercentage /= 100

		memoryBytes := int64(memoryPercentage * float64(totalMemoryBytes))
		if memoryBytes >= limit {
			pids = append(pids, pid)
		}
	}

	return pids, nil
}

func terminateProcess(pid int) error {
	return errors.Join(
		syscall.Kill(pid, syscall.SIGTERM),
		syscall.Kill(pid, syscall.SIGINT),
		syscall.Kill(pid, syscall.SIGKILL),
	)
}

func waitForTermination() <-chan struct{} {
	sig := make(chan os.Signal, 1)
	done := make(chan struct{})

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		close(done)
	}()

	return done
}

func totalMemory() (int64, error) {
	sysInfo := new(syscall.Sysinfo_t)
	if err := syscall.Sysinfo(sysInfo); err != nil {
		return 0, err
	}

	return int64(sysInfo.Totalram) * int64(sysInfo.Unit), nil
}

func maxMemory() int64 {
	if len(os.Args) != 2 {
		fmt.Println("Usage: program_name <max_memory_usage_in_megabytes>")
		os.Exit(1)
	}

	maxMemoryUsage, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Println("Invalid max_memory_usage argument:", err)
		os.Exit(1)
	}

	return maxMemoryUsage * bytesInMB
}
