package main

import (
	"fmt"
	"github.com/abdfnx/gosh"
	. "github.com/klauspost/cpuid/v2"
	"github.com/rdegges/go-ipify"
	"github.com/ricochet2200/go-disk-usage/du"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"gopkg.in/mail.v2"
	"log"
	"net/http"
	"os"
	"os/exec"
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

// This checks the error given from running a program
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func showDraft(email string) {
	f, err := os.Create("./dependencies/index.html")
	check(err)
	defer f.Close()

	email = strings.Replace(email, "cid:", "", -1)
	_, err = f.WriteString(email)
	check(err)
	f.Sync()

	fs := http.FileServer(http.Dir("./dependencies"))
	http.Handle("/", fs)

	//multithreading to keep the server up but finish the program
	go func() {
		err = http.ListenAndServe(":3000", nil)
	}()

	err = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:3000").Start()
}

func sendEmail(m mail.Message, email string, usrname string, passwd string) {
	m.SetBody("text/html", email)

	d := mail.NewDialer("smtp.office365.com", 587, usrname, passwd)

	if err := d.DialAndSend(&m); err != nil {

		panic(err)

	}
}

func cleanup() {
	os.Remove("./dependencies/index.html")
}

func main() {

	//Setting the credentials
	var to string = os.Getenv("SMTP_TO")
	var from string = os.Getenv("SMTP_FROM")
	var pass string = os.Getenv("SMTP_PASS")

	//Getting memory variables
	vmem, err := mem.VirtualMemory()
	check(err)
	freeMem := fmt.Sprintf("%.2f GB", float64(vmem.Free)/float64(GB))
	totalMem := fmt.Sprintf("%.2f GB", float64(vmem.Total)/float64(GB))

	//Getting CPU & Host Info
	cpuSpeed := fmt.Sprintf("%.2f GHz", float64(CPU.Hz)/float64(GHz))
	hostInfo, err := host.Info()
	check(err)

	//Getting services using powershell
	err, out, errout := gosh.PowershellOutput(`Get-Service | Where-Object {$_.Status -eq "Running"} | Select Name`)
	services := strings.ReplaceAll(out, "\r\n", "<br>")
	if err != nil {
		log.Printf("error: %v\n", err)
		fmt.Print(errout)
	}

	//getting disk usage for the C drive
	usage := du.NewDiskUsage("C:\\")
	usageString := fmt.Sprintf("%.2f GB", float64(usage.Available())/float64(GB))

	//Get the IP using the ipify API
	ip, err := ipify.GetIp()
	check(err)

	//Build the message
	m := mail.NewMessage()

	m.SetHeader("From", from)

	m.SetHeader("To", to)

	m.SetHeader("Subject", "Your System Information")

	//Reading the email.html and using it as a template
	fileBytes, err := os.ReadFile("email.html")
	fileString := string(fileBytes)
	m.Embed("dependencies/logo-no-background.png")
	fileString = strings.Replace(fileString, "[public-ip]", ip, -1)
	fileString = strings.Replace(fileString, "[mem-info]", freeMem, -1)
	fileString = strings.Replace(fileString, "[disk-space]", usageString, -1)
	fileString = strings.Replace(fileString, "[os-name]", runtime.GOOS, -1)
	fileString = strings.Replace(fileString, "[os-type]", hostInfo.Platform, -1)
	fileString = strings.Replace(fileString, "[mem-total]", totalMem, -1)
	fileString = strings.Replace(fileString, "[CPU]", CPU.BrandName, -1)
	fileString = strings.Replace(fileString, "[CPU-Speed]", cpuSpeed, -1)
	fileString = strings.Replace(fileString, "[services]", services, -1)

	//Checking if the user would like to see a draft of their
	fmt.Println("Would you like to see a draft (Y,N):")
	var response string
	fmt.Scanln(&response)
	if response == "Y" {
		showDraft(fileString)
	}

	//Checking if I would like to send the email
	fmt.Println("Would you like to send the email (Y,N):")
	var response2 string
	fmt.Scanln(&response2)
	if response2 == "Y" {
		fmt.Println("Sending email")
		sendEmail(*m, fileString, from, pass)
	} else {
		fmt.Println("Quitting the Program...")
		cleanup()
		os.Exit(0)
	}
	cleanup()

}
