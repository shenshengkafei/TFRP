package rpforrp

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func iterDirectory(dirPath string, tw *tar.Writer, originPath string) {
	dir, _ := os.Open(dirPath)
	defer dir.Close()
	fis, _ := dir.Readdir(0)
	for _, fi := range fis {
		curPath := dirPath + "/" + fi.Name()
		if fi.IsDir() {
			iterDirectory(curPath, tw, originPath)
		} else {
			fmt.Printf("adding... %s\n", curPath)

			dockerFile := curPath[len(originPath):len(curPath)]
			dockerFileReader, err := os.Open(dir.Name() + string(filepath.Separator) + fi.Name())
			if err != nil {
				log.Fatal(err, " :unable to open Dockerfile")
			}
			readDockerFile, err := ioutil.ReadAll(dockerFileReader)
			if err != nil {
				log.Fatal(err, " :unable to read dockerfile")
			}

			tarHeader := &tar.Header{
				Name: dockerFile,
				Size: int64(len(readDockerFile)),
			}
			err = tw.WriteHeader(tarHeader)
			if err != nil {
				log.Fatal(err, " :unable to write tar header")
			}
			_, err = tw.Write(readDockerFile)
			if err != nil {
				log.Fatal(err, " :unable to write tar body")
			}
		}
	}
}

// BuildAndPushImageFromDockerfile will build a image and push it to dockerhub.
func BuildAndPushImageFromDockerfile(swaggerPath string, imageName string, dockerhubAccount string, dockerhubPassword string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err, " :unable to init client")
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	originPath := swaggerPath + "src/IO.Swagger"

	dockerFile := "myDockerfile"
	dockerFileReader, err := os.Open(originPath + "/Dockerfile")
	if err != nil {
		log.Fatal(err, " :unable to open Dockerfile")
	}
	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	if err != nil {
		log.Fatal(err, " :unable to read dockerfile")
	}

	tarHeader := &tar.Header{
		Name: dockerFile,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		log.Fatal(err, " :unable to write tar header")
	}
	_, err = tw.Write(readDockerFile)
	if err != nil {
		log.Fatal(err, " :unable to write tar body")
	}

	iterDirectory(originPath, tw, originPath)

	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Tags:       []string{imageName},
			Context:    dockerFileTarReader,
			Dockerfile: dockerFile,
			NoCache:    true,
			Remove:     true})
	if err != nil {
		log.Fatal(err, " :unable to build docker image")
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}

	auth := types.AuthConfig{
		Username: dockerhubAccount,
		Password: dockerhubPassword,
	}
	authBytes, _ := json.Marshal(auth)
	authBase64 := base64.URLEncoding.EncodeToString(authBytes)

	imagePushResponse, err := cli.ImagePush(
		context.Background(),
		imageName,
		types.ImagePushOptions{
			All:          true,
			RegistryAuth: authBase64,
		})
	if err != nil {
		log.Fatal(err, " :unable to push docker image")
	}
	defer imagePushResponse.Close()
	_, err = io.Copy(os.Stdout, imagePushResponse)
	if err != nil {
		log.Fatal(err, " :unable to read image push response")
	}
}
