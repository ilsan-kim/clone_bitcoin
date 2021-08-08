package explorer

import (
	"fmt"
	"github.com/ddhyun93/seancoin/blockchain"
	"log"
	"net/http"
	"text/template"
)

const templateDir string = "templates/"

// 이미 있는 템플릿 파일을 파싱하지않고, 정의해둔 템플릿을 "load" 하기 위한 객체
var templates *template.Template

type homeData struct {
	// home.gohtml 이라는 템플릿 (외부파일)에 데이터를 보낼거니까 public 으로 만들어야함
	PageTitle 	string
	Blocks		[]*blockchain.Block
}

func Start(port int) {
	handler := http.NewServeMux()
	// templates 라는 변수에 템플릿 담기
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))		// text/template 패키지의 템플릿
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))	// templates 변수  -> pages와 partials 두 폴더로 나누지 않았으면 이러지 않아도 댐

	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	// 에러가 있다면 에러를 출력하며 끝나고, 에러가 없다면 함수를 실행하는 것이 log.Fatal의 기능
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

// 데이터를 쓰는 writer는 포인터가 필요없고, 리턴해주는 request는 파일이 될수도있으니 포인터를 사용해야한다.
func home(writer http.ResponseWriter,  request *http.Request) {
	data := homeData{"동더리움", blockchain.GetBlockChain().AllBlocks()}
	templates.ExecuteTemplate(writer, "home", data)
}

func add(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		templates.ExecuteTemplate(writer, "add", nil)
	case "POST":	// POST일떄는 유저가 입력한 데이터 (폼데이터 중 name="blockData"인것)를 가져와서 새로운 블록을 생성해야함
		request.ParseForm()
		data := request.Form.Get("blockData")
		blockchain.GetBlockChain().AddBlock(data)
		http.Redirect(writer, request, "/home", http.StatusPermanentRedirect)
	}
}