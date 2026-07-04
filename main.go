package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func printMountNS(where string) {
	out, err := exec.Command("readlink", "/proc/self/ns/mnt").Output()
	must(err)
	fmt.Printf("%s (PID %d): %s", where, os.Getpid(), out)
}

func main() {
	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		panic("what should I do")
	}
}

func parent() {
	//for windows
	/*exe, winErr := os.Executable()
	if winErr != nil {
		panic(winErr)
	}
	cmd := exec.Command(exe, append([]string{"child"}, os.Args[2:]...)...)*/

	//for wsl or linux
	printMountNS("Parent")
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS, //todo: add more namespaces and find a way to acess them
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
}

func child() {
	must(syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))
	must(os.Mkdir("rootfs/oldrootfs", 0700))
	must(syscall.PivotRoot("rootfs", "rootfs/oldrootfs"))
	must(os.Chdir("/"))
	printMountNS("Child")
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
