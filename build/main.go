package main

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v6"
	"io"
	"os"
	"os/exec"
	"sync"
)

func main() {

	codePath := "/tmp/foo"
	imageTag := "crpi-mgl4ujhwwhrsi5e3.cn-hangzhou.personal.cr.aliyuncs.com/kubeland/test:v1"
	cloneCode(codePath)
	buildImage(imageTag)
	pushImage(imageTag)
}
func cloneCode(codePath string) {
	fmt.Println("###### Step1: Clone code")
	gitUrl := os.Getenv("GIT_URL")
	if gitUrl == "" {
		panic("GIT_URL is empty")
	}
	fmt.Printf("clone code from %s\n", gitUrl)
	_, err := git.PlainClone(codePath, &git.CloneOptions{
		URL:      gitUrl,
		Progress: os.Stdout,
	})
	if err != nil {
		panic(err)
	}
}

func buildImage(imageTag string) {
	fmt.Println("###### Step2: Build image")
	cmd := exec.Command("buildah", "build",
		"--tag", imageTag,
		"/app/config",
	)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	if err = cmd.Start(); err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		printOutput(stdoutPipe)
	}()
	go func() {
		defer wg.Done()
		printOutput(stderrPipe)
	}()
	wg.Wait()
	if err = cmd.Wait(); err != nil {
		panic(err)
	}
}

func pushImage(imageTag string) {
	fmt.Println("###### Step3: Push image")

	//cmd := exec.Command("buildah", "push",
	//	imageTag,
	//)
	//stdoutPipe, err := cmd.StdoutPipe()
	//if err != nil {
	//	panic(err)
	//}
	//stderrPipe, err := cmd.StderrPipe()
	//if err != nil {
	//	panic(err)
	//}
	//if err = cmd.Start(); err != nil {
	//	panic(err)
	//}
	//var wg sync.WaitGroup
	//wg.Add(2)
	//go func() {
	//	defer wg.Done()
	//	printOutput(stdoutPipe)
	//}()
	//go func() {
	//	defer wg.Done()
	//	printOutput(stderrPipe)
	//}()
	//wg.Wait()
	//if err = cmd.Wait(); err != nil {
	//	panic(err)
	//}

}

func printOutput(pipe io.Reader) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%s\n", line)
	}

}
