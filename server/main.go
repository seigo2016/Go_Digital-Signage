package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	// "os/exec"
	"image/jpeg"
	"path/filepath"
	"syscall"
	"time"

	"github.com/nfnt/resize"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"

	"github.com/gen2brain/go-fitz"

	// "gopkg.in/gographics/imagick.v3/imagick"
	"gopkg.in/ini.v1"
)

// type Template struct {
// 	templates *template.Template
// }

// func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
// 	return t.templates.ExecuteTemplate(w, name, data)
// }

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
	// err := exec.Command("open", authURL).Start()
	fmt.Println(authURL)
	// if err != nil {
	//	panic(err)
	// }
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
	p, _ := os.Executable()
	p = filepath.Dir(p)
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
		fileName = append(fileName, "404.pdf")
		src, err := os.Open(filepath.Join(p, "views", "404.pdf"))
		if err != nil {
			panic(err)
		}
		defer src.Close()
		dst, err := os.Create(filepath.Join(p, "tmp", "404.pdf"))
		if err != nil {
			panic(err)
		}
		defer dst.Close()
		_, err = io.Copy(dst, src)
		if err != nil {
			panic(err)
		}
	} else {
		for _, i := range r.Files {
			r, err := srv.Files.Get(i.Id).Download()
			if err != nil {
				log.Fatalf("%v", err)
			}
			output, err := os.Create(filepath.Join(p, "tmp", i.Name))
			io.Copy(output, r.Body)
			fileName = append(fileName, i.Name)
		}
	}
	// imagick.Initialize()
	// defer imagick.Terminate()
	var imageList []string
	// const layout = "Monday-Jan-02-15:04:05-JST-2006"
	for _, i := range fileName {
		// t := time.Now()
		name := strings.Replace(i, ".pdf", "", -1)
		doc, err := fitz.New(filepath.Join(p, "tmp", i))
		if err != nil {
			panic(err)
		}
		defer doc.Close()
		for n := 0; n < 1; n++ {
			img, err := doc.Image(n)
			if err != nil {
				panic(err)
			}

			f, err := os.Create(filepath.Join(p, "page", name+".jpeg"))
			if err != nil {
				panic(err)
			}
			img = resize.Thumbnail(1080, 1920, img, resize.Lanczos3)
			err = jpeg.Encode(f, img, &jpeg.Options{jpeg.DefaultQuality})
			if err != nil {
				panic(err)
			}

			f.Close()

		}
		imageList = append(imageList, name+".jpeg")
	}
	return imageList, len(imageList)
}

func image2base64(path string) ([]byte, int64) {
	p, _ := os.Executable()
	p = filepath.Dir(p)
	file, _ := os.Open(filepath.Join(p, "page", path))
	defer file.Close()
	f, err := file.Stat()
	if err != nil {
		fmt.Printf("The file is %d bytes long %v", err)
	}
	size := f.Size()
	// size = fmt.Printf("%06d", size)
	data := make([]byte, size)
	// sizeb := make([]byte, binary.MaxVarintLen64)
	// binary.PutVarint(sizeb, size)
	file.Read(data)
	data = append(data, []byte("\n\n")...)
	// data = append(sizeb, data...)
	return data, size
}

func makeTmpFolder() {
	p, _ := os.Executable()
	p = filepath.Dir(p)
	if err := os.RemoveAll(filepath.Join(p, "tmp")); err != nil {
		fmt.Println(err)
	}
	if err := os.Mkdir(filepath.Join(p, "tmp"), 0777); err != nil {
		fmt.Println(err)
	}
	if err := os.RemoveAll(filepath.Join(p, "page")); err != nil {
		fmt.Println(err)
	}
	if err := os.Mkdir(filepath.Join(p, "page"), 0777); err != nil {
		fmt.Println(err)
	}
}
func setLog() {
	p, _ := os.Executable()
	p = filepath.Dir(p)
	f, err := os.OpenFile(filepath.Join(p, "signage_server.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
}

func openConf() string {
	p, _ := os.Executable()
	p = filepath.Dir(p)
	c, _ := ini.Load(filepath.Join(p, "config.ini"))
	return c.Section("INFO").Key("FOLDERID").String()
}

func main() {
	p, _ := os.Executable()
	p = filepath.Dir(p)
	folderID := openConf()
	makeTmpFolder()
	setLog()
	ctx := context.Background()
	b, err := ioutil.ReadFile(filepath.Join(p, "client_secret.json"))
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
				body, _ := image2base64(imageList[cnt%filelen])
				_, werr := conn.Write(body)
				if werr != nil {
					fmt.Printf("Error: %v\n", werr)
					if errors.Is(werr, syscall.EPIPE) {
						fmt.Println("Disconnected")
						break
					}
				}
				time.Sleep(time.Second * 20)
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
