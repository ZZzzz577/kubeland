package main

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
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
	//pushImage(imageTag)
	//time.Sleep(time.Hour)
}
func cloneCode(codePath string) {
	fmt.Println("###### Step1: Clone code")
	gitUrl := os.Getenv("GIT_URL")
	if gitUrl == "" {
		panic("GIT_URL is empty")
	}

	// todo: consider no git token
	gitTokenFile := "/app/config/git/GIT_TOKEN"
	file, err := os.Open(gitTokenFile)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()
	gitToken := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		gitToken += scanner.Text() + "\n"
	}
	if err = scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println("GIT_TOKEN:", gitToken)

	fmt.Printf("clone code from %s\n", gitUrl)
	repository, err := git.PlainClone(codePath, &git.CloneOptions{
		URL: gitUrl,
		//Auth: &http.BasicAuth{
		//	Username: "token",
		//	Password: gitToken,
		//},
		Progress: os.Stdout,
	})
	if err != nil {
		panic(err)
	}
	gitCommit := os.Getenv("GIT_COMMIT")
	if gitCommit == "" {
		panic("GIT_COMMIT is empty")
	}

	worktree, err := repository.Worktree()
	if err != nil {
		panic(err)
	}
	hash, ok := plumbing.FromHex(gitCommit)
	if !ok {
		panic("invalid git commit")
	}
	if err = worktree.Checkout(&git.CheckoutOptions{Hash: hash}); err != nil {
		panic(err)
	}

}

func buildImage(imageTag string) {
	fmt.Println("###### Step2: Build image")
	cmd := exec.Command("buildah", "build",
		"--tag", imageTag,
		"/app/config/dockerfile",
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
