package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/micmonay/keybd_event"
)

type PowerShell struct {
	powerShell string
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "jonatan.me:9009")
	if err != nil {
		log.Println("Couldn't resolve the address")
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println("Couldn't connect to the server")
	}
	posh := New()
	for {
		// receive message
		var buff [30]byte
		_, err = conn.Read(buff[0:])
		if err != nil {
			log.Println("Couldn't receive message")
		}
		var messageWithoutExcessBytes []byte = bytes.Trim(buff[:], "\x00")

		var message string = BytesToString(messageWithoutExcessBytes[:]) // Convert to string
		log.Println(message)
		if strings.Contains(message, "skip") {
			log.Println("Skipped song")
			go skip()
		} else if strings.Contains(message, "timesync") {
			t := cleanUnixString(message)
			command := fmt.Sprintf("Set-Date %d:%d:%d", t.Hour(), t.Minute(), t.Second())
			cmd := fmt.Sprintf("%s", command)
			stdOut, stdErr, err := posh.execute(cmd)
			fmt.Printf("\nCommand : %s\nStdOut : '%s'\nStdErr: '%s'\nErr: %s", command, strings.TrimSpace(stdOut), stdErr, err)
		}
	}

}

func cleanUnixString(message string) time.Time {
	var (
		timeString    = strings.Split(message, ":")
		timeUncleaned = timeString[1]
		timeCleaned   = timeUncleaned[:len(timeUncleaned)-1]
		timeUnix, _   = strconv.ParseInt(timeCleaned, 10, 64)
	)
	return time.Unix(timeUnix/1000, 0)
}

func BytesToString(b []byte) string {
	var strToConvert string = bytes.NewBuffer(b).String()
	return strToConvert
}

func skip() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	kb.SetKeys(keybd_event.VK_MEDIA_NEXT_TRACK)
	kb.Press()
}

func New() *PowerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		powerShell: ps,
	}
}

func (p *PowerShell) execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.powerShell, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}

var (
	admin = `
	$myWindowsID = [System.Security.Principal.WindowsIdentity]::GetCurrent();
	$myWindowsPrincipal = New-Object System.Security.Principal.WindowsPrincipal($myWindowsID);
	$adminRole = [System.Security.Principal.WindowsBuiltInRole]::Administrator;
	if (-Not ($myWindowsPrincipal.IsInRole($adminRole))) {
		$newProcess = New-Object System.Diagnostics.ProcessStartInfo "PowerShell";
		$newProcess.Arguments = "& '" + $script:MyInvocation.MyCommand.Path + "'"
		$newProcess.Verb = "runAs";
		[System.Diagnostics.Process]::Start($newProcess);
	}
	`
)
