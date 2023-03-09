package main

import (
	"bufio"
	"fmt"
	"github.com/abdfnx/gosh"
	_ "github.com/abdfnx/gosh"
	. "github.com/klauspost/cpuid/v2"
	"github.com/rdegges/go-ipify"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"gopkg.in/mail.v2"
	"log"
	"os"
	"runtime"
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

	m := mail.NewMessage()

	m.SetHeader("From", fileLines[1])

	m.SetHeader("To", fileLines[0])

	m.SetHeader("Subject", "Hello!")

	fileBytes, err := os.ReadFile("email.html")
	fileString := string(fileBytes)
	m.Embed("images/logo-no-background.png")
	m.Embed("images/logo-color.png")
	m.SetBody("text/html", fileString)

	fmt.Print(runtime.GOOS + "\n")
	vmem, err := mem.VirtualMemory()
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", vmem.Total, vmem.Free, vmem.UsedPercent)
	fmt.Println("\nName:", CPU.BrandName, "\nSpeed:", CPU.Hz)

	hostInfo, err := host.Info()
	fmt.Printf("Platform: %v\n", hostInfo.Platform)
	fmt.Printf("Total Processes: %d\n", hostInfo.Procs)
	err, out, errout := gosh.PowershellOutput(`Get-Service | Where-Object {$_.Status -eq "Running"}`)

	if err != nil {
		log.Printf("error: %v\n", err)
		fmt.Print(errout)
	}

	fmt.Print(out)

	ip, err := ipify.GetIp()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(ip)

	d := mail.NewDialer("smtp.office365.com", 587, fileLines[1], fileLines[2])

	// Send the email to Kate, Noah and Oliver.

	if err := d.DialAndSend(m); err != nil {

		panic(err)

	}

}
