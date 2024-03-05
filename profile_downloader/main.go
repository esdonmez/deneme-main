package main

import (
	"context"
	"fmt"
	pgo "github.com/esdonmez/deneme-main"
	"os/exec"
	"strings"
)

func main() {
	ctx := context.Background()

	endpoint := "localhost:9000"
	accessKeyID := "admin"
	secretAccessKey := "admin123"

	minioUploader := pgo.NewUploader(endpoint, accessKeyID, secretAccessKey, "")

	// Make a new bucket called pgos.
	bucketName := "pgos"
	location := "us-east-1"

	minioUploader.Create(ctx, bucketName, location)

	// Download the files
	objectNames := minioUploader.DownloadAll(ctx, bucketName)

	paths := fmt.Sprintf("../ %s.pprof", strings.Join(objectNames, ".pprof "))
	command := fmt.Sprintf("go tool pprof -proto %s > merged.pprof", paths)

	// create a new *Cmd instance
	cmd := exec.Command("bash", "-c", command)

	// The `Output` method executes the command and
	// collects the output, returning its value
	out, err := cmd.Output()
	if err != nil {
		// if there was any error, print it here
		fmt.Println("could not run command: ", err)
	}
	// otherwise, print the output from running the command
	fmt.Println("Output: ", string(out))
}
