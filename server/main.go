package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"syscall"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"gopkg.in/gographics/imagick.v3/imagick"
	"gopkg.in/ini.v1"
)

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	fmt.Println("Get token from file")
	tok, err := tokenFromFile("token.json")
	if err != nil {
		fmt.Printf("Token file not found \"%v\"\n", err)
		tok = getTokenFromWeb(config)
		if err != nil {
			fmt.Println("error")
		}
		// tok, err := config.Exchange(oauth2.NoContext, tok)
		fmt.Printf("token get %v", tok)
		if tok != nil {
			saveToken("token.json", tok)
		}
	}
	return config.Client(ctx, tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	err := exec.Command("open", authURL).Start()
	if err != nil {
		panic(err)
	}
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		fmt.Printf("Err %v\n", err)
	}
	code, _ = url.QueryUnescape(code)
	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("Err %v\n", err)
	}
	fmt.Println("getTokenFromWeb Complete")
	return tok
}

func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Save token as file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Save error %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getPDF(client *http.Client, folderID string) ([]string, int) {
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	r, err := srv.Files.List().Q("mimeType='application/pdf' and trashed = false and '" + folderID + "' in parents").PageSize(30).Fields("nextPageToken, files(id, name, mimeType)").
		Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	var fileName []string
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			r, err := srv.Files.Get(i.Id).Download()
			if err != nil {
				log.Fatalf("%v", err)
			}
			output, err := os.Create("tmp/" + i.Name)
			io.Copy(output, r.Body)
			// fmt.Println(n)
			fileName = append(fileName, i.Name)
		}
	}
	imagick.Initialize()
	defer imagick.Terminate()
	var imageList []string
	for _, i := range fileName {
		mw := imagick.NewMagickWand()
		defer mw.Destroy()
		mw.ReadImage("tmp/" + i)
		mw.SetIteratorIndex(0)
		mw.SetImageFormat("jpg")
		mw.WriteImage("page/" + i + ".jpg")
		imageList = append(imageList, i+".jpg")
	}
	return imageList, len(imageList)
}

func image2base64(path string) []byte {
	file, _ := os.Open("page/" + path)
	defer file.Close()
	f, err := file.Stat()
	if err != nil {
		fmt.Printf("The file is %d bytes long %v", err)
	}
	size := f.Size()
	data := make([]byte, size)
	file.Read(data)
	return data
}

func makeTmpFolder() {
	if err := os.RemoveAll("tmp"); err != nil {
		fmt.Println(err)
	}
	if err := os.Mkdir("tmp", 0777); err != nil {
		fmt.Println(err)
	}
	if err := os.RemoveAll("page"); err != nil {
		fmt.Println(err)
	}
	if err := os.Mkdir("page", 0777); err != nil {
		fmt.Println(err)
	}
}
func setLog() {
	f, err := os.OpenFile("signage_server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
}

func openConf() string {
	c, _ := ini.Load("config.ini")
	return c.Section("INFO").Key("FOLDERID").String()
}

func main() {
	folderID := openConf()
	makeTmpFolder()
	setLog()
	ctx := context.Background()
	b, err := ioutil.ReadFile("client_secret.json")
	config, err := google.ConfigFromJSON(b, drive.DriveReadonlyScope, drive.DriveFileScope)
	if err != nil {
		log.Println(err)
	}
	client := getClient(ctx, config)
	var imageList []string
	var filelen int
	l, err := net.Listen("tcp4", "0.0.0.0:30000")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	for {
		imageList, filelen = getPDF(client, folderID)
		cnt := 0
		fmt.Println("wait")
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go func(conn net.Conn) {
			fmt.Println("Connected")
			for {
				fmt.Println(cnt)
				data := image2base64(imageList[cnt%filelen])
				time.Sleep(time.Second * 20)
				_, werr := conn.Write(data)
				if werr != nil {
					if opErr, ok := werr.(*net.OpError); ok {
						if sysErr, okok := opErr.Err.(*os.SyscallError); okok && sysErr.Err == syscall.EPIPE {
							fmt.Println("Disconnected")
							break
						}
					}
				}
				if cnt == (filelen*2)-1 {
					imageList, filelen = getPDF(client, folderID)
					cnt = 0
				} else {
					cnt++
				}
			}
			conn.Close()
		}(conn)
	}
}
