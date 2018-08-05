package handler

import (
	"fmt"
	"net/http"

	"github.com/kataras/iris"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// @Summary Shows OK as the ping-pong result
// @Description Shows OK as the ping-pong result
// @Tags sd
// @Accept  json
// @Produce  json
// @Success 200 {string} plain "OK"
// @Router /sd/health [get]
func HealthCheck(c iris.Context) {
	message := "OK"
	c.ContentType("text/plain")
	c.StatusCode(http.StatusOK)
	c.WriteString("\n"+message)
}

// @Summary Checks the disk usage
// @Description Checks the disk usage
// @Tags sd
// @Accept  json
// @Produce  json
// @Success 200 {string} plain "OK - Free space: 17233MB (16GB) / 51200MB (50GB) | Used: 33%"
// @Router /sd/disk [get]
func DiskCheck(c iris.Context) {
	u, _ := disk.Usage("/")

	usedMB := int(u.Used) / MB
	usedGB := int(u.Used) / GB
	totalMB := int(u.Total) / MB
	totalGB := int(u.Total) / GB
	usedPercent := int(u.UsedPercent)

	status := http.StatusOK
	text := "OK"

	if usedPercent >= 95 {
		status = http.StatusOK
		text = "CRITICAL"
	} else if usedPercent >= 90 {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
	c.ContentType("text/plain")
	c.StatusCode(status)
	c.WriteString("\n"+message)
}

// @Summary Checks the cpu usage
// @Description Checks the cpu usage
// @Tags sd
// @Accept  json
// @Produce  json
// @Success 200 {string} plain "CRITICAL - Load average: 1.78, 1.99, 2.02 | Cores: 2"
// @Router /sd/cpu [get]
func CPUCheck(c iris.Context) {
	cores, _ := cpu.Counts(false)

	a, _ := load.Avg()
	l1 := a.Load1
	l5 := a.Load5
	l15 := a.Load15

	status := http.StatusOK
	text := "OK"

	if l5 >= float64(cores-1) {
		status = http.StatusInternalServerError
		text = "CRITICAL"
	} else if l5 >= float64(cores-2) {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	message := fmt.Sprintf("%s - Load average: %.2f, %.2f, %.2f | Cores: %d", text, l1, l5, l15, cores)
	c.ContentType("text/plain")
	c.StatusCode(status)
	c.WriteString("\n"+message)
}

// @Summary Checks the ram usage
// @Description Checks the ram usage
// @Tags sd
// @Accept  json
// @Produce  json
// @Success 200 {string} plain "OK - Free space: 402MB (0GB) / 8192MB (8GB) | Used: 4%"
// @Router /sd/ram [get]
func RAMCheck(c iris.Context) {
	u, _ := mem.VirtualMemory()

	usedMB := int(u.Used) / MB
	usedGB := int(u.Used) / GB
	totalMB := int(u.Total) / MB
	totalGB := int(u.Total) / GB
	usedPercent := int(u.UsedPercent)

	status := http.StatusOK
	text := "OK"

	if usedPercent >= 95 {
		status = http.StatusInternalServerError
		text = "CRITICAL"
	} else if usedPercent >= 90 {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
	c.ContentType("text/plain")
	c.StatusCode(status)
	c.WriteString("\n"+message)
}
