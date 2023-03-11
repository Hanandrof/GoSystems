package main

import (
	"bufio"
	"fmt"
	"github.com/abdfnx/gosh"
	_ "github.com/abdfnx/gosh"
	. "github.com/klauspost/cpuid/v2"
	"github.com/rdegges/go-ipify"
	"github.com/ricochet2200/go-disk-usage/du"
	_ "github.com/ricochet2200/go-disk-usage/du"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"gopkg.in/mail.v2"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	B   = 1
	KB  = 1024 * B
	MB  = 1024 * KB
	GB  = 1024 * MB
	GHz = GB
)

func main() {

	readFile, err := os.Open("config")
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	fmt.Print(runtime.GOOS + "\n")
	vmem, err := mem.VirtualMemory()
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", vmem.Total, vmem.Free, vmem.UsedPercent)
	fmt.Printf("Available Memory: %.2f GB\n", float64(vmem.Free)/float64(GB))
	freeMem := fmt.Sprintf("%.2f GB", float64(vmem.Free)/float64(GB))
	totalMem := fmt.Sprintf("%.2f GB", float64(vmem.Total)/float64(GB))

	fmt.Println("\nName:", CPU.BrandName)
	fmt.Printf("Speed: %.2f GHZ\n", float64(CPU.Hz)/float64(GHz))
	cpuSpeed := fmt.Sprintf("%.2f GHz\n", float64(CPU.Hz)/float64(GHz))
	hostInfo, err := host.Info()
	fmt.Printf("Platform: %v\n", hostInfo.Platform)
	err, out, errout := gosh.PowershellOutput(`Get-Service | Where-Object {$_.Status -eq "Running"} | Select Name`)

	if err != nil {
		log.Printf("error: %v\n", err)
		fmt.Print(errout)
	}

	//fmt.Print(out)
	//fmt.Println(strings.Replace(out, string('\n'), "<br>", -1))
	fmt.Println(strings.ReplaceAll(out, "\r\n", "<br>"))

	usage := du.NewDiskUsage("C:\\")
	fmt.Printf("Disk Usage %.2f\n", float64(usage.Available())/float64(GB))
	usageString := fmt.Sprintf("%.2f\n", float64(usage.Available())/float64(GB))
	fmt.Print(usageString)
	ip, err := ipify.GetIp()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(ip)

	m := mail.NewMessage()

	m.SetHeader("From", fileLines[1])

	m.SetHeader("To", fileLines[0])

	m.SetHeader("Subject", "Your System Information")

	fileBytes, err := os.ReadFile("email.html")
	fileString := string(fileBytes)
	m.Embed("images/logo-no-background.png")
	fileString = strings.Replace(fileString, "[public-ip]", ip, -1)
	fileString = strings.Replace(fileString, "[mem-info]", freeMem, -1)
	fileString = strings.Replace(fileString, "[disk-space]", usageString, -1)
	fileString = strings.Replace(fileString, "[os-name]", runtime.GOOS, -1)
	fileString = strings.Replace(fileString, "[os-type]", hostInfo.Platform, -1)
	fileString = strings.Replace(fileString, "[mem-total]", totalMem, -1)
	fileString = strings.Replace(fileString, "[CPU]", CPU.BrandName, -1)
	fileString = strings.Replace(fileString, "[CPU-Speed]", cpuSpeed, -1)
	fileString = strings.Replace(fileString, "[services]", strings.ReplaceAll(out, "\r\n", "<br>"), -1)
	m.SetBody("text/html", fileString)

	d := mail.NewDialer("smtp.office365.com", 587, fileLines[1], fileLines[2])

	if err := d.DialAndSend(m); err != nil {

		panic(err)

	}

}
