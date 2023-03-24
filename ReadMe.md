![logo](dependencies/logo-color.png)

This program is built to use SMTP to send an email with your system information using GoLang. It uses environment
variables from a PowerShell script and accesses them when running a Go Program

## Dependencies
+ Go
  + Make sure the installation is referenced in your **PATH** environment variable
+ PowerShell
+ An **Outlook** email
  + Google took out the ability for insecure connections so outlook is currently the only way to send emails

## Installation
```
git clone https://github.com/Hanandrof/GoSystems.git
cd GoSystems
powershell script.ps1
```

## How It Works

This section is going to quickly talk about how everything works, it will be split into two sections, powershell and GO

### Powershell

Powershell uses environmental variables and collects your information via simple input and output. It then executes the GO program and subsequently deletes the environment variables.

### GO

Go uses different API's to collect your system information such as your IP address. It then also grabs other system information using different GO libraries all detailed in the go.mod requirements. It uses GoMail to send an email using any outlook email to any email.

When viewing a draft of your email it will serve the html form on a localhost and open your default web browser. 

It builds your html email using a template and then it uploads the contents of that file to GoMail, it then deletes the newly created file and keeps the template
