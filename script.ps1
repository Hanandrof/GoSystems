$To = Read-Host -Prompt 'Who are you sending this email to?'
$From = Read-Host -Prompt 'What email are you sending this from?'
$pass = Read-Host -Prompt 'What is your password?' -AsSecureString

$pass = [Runtime.InteropServices.Marshal]::PtrToStringAuto(
        [Runtime.InteropServices.Marshal]::SecureStringToBSTR($pass))

[Environment]::SetEnvironmentVariable('SMTP_TO',$To)
[Environment]::SetEnvironmentVariable('SMTP_FROM',$From)
[Environment]::SetEnvironmentVariable('SMTP_PASS',$pass)

go run smtp.go

[Environment]::SetEnvironmentVariable('SMTP_TO','')
[Environment]::SetEnvironmentVariable('SMTP_FROM','')
[Environment]::SetEnvironmentVariable('SMTP_PASS','')