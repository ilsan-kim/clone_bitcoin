package rest

import (
	"encoding/json"
	"fmt"
	"github.com/ddhyun93/seancoin/blockchain"
	"github.com/ddhyun93/seancoin/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var port string
type url string

type urlDescription struct {
	URL         url    `json:"url"`               // json은 소문자여야하는데 그럼 패키지 외부에서 쓸수없다. 이 때 필요한게 struct field tag
	Method      string `json:"method"`            // 뜻은 "내 struct가 json이라면, Field가 특정한 방식으로 보여질 것"이라고 설정하는 것
	Description string `json:"description"`       // 백틱으로 따옴표 감싸는 형식
	Payload     string `json:"payload,omitempty"` // omitempty를 붙이면 해당 밸류의 값이 없을 때 숨김
}
type addBlockBody struct {
	Message		string
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func (u url) marshalText() ([]byte, error) { // 내장 인터페이스인 MarshalText를 사용하여 문자열 포매팅
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

func Start(aPort int) {
	router := mux.NewRouter()	// multiplexer 라는 뜻으로 url 헨들러를 분리하여 html파일로 접근하는 explorer패키지/ restAPI서버 두개를 동시에 띄울 수 있게 해줌
	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleware)	// router 가 호출될 때 이 함수가 실행 된 후 nextHandler가 실행됨
	router.HandleFunc("/", documentation).Methods("GET")	// gorilla mux의 특징으로 어떤 라우터가 어떤 메서드를 처리할지 특정할 수 있음 >> Method not allowed 처리를 안해도됨
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func documentation(writer http.ResponseWriter, request *http.Request) {
	data := []urlDescription{
		{
			URL:			url("/"),
			Method: 		"GET",
			Description: 	"See Documentation",
			// 이건 아직 JSON 이 아님 -> Marshal 을 통해 Go의 Struct를 JSON으로 encoding한 interface를 리턴함 (https://en.wikipedia.org/wiki/Marshalling_(computer_science))
		},
		{
			URL:			url("/blocks"),
			Method: 		"POST",
			Description:	"Add A Block",
			Payload:		"data:string",
		},
		{
			URL:			url("/blocks"),
			Method: 		"GET",
			Description:	"See All Block",
		},
		{
			URL:			url("/blocks/{hash}"),
			Method: 		"GET",
			Description:	"See A Block",
		},
	}
	/*	[어려운방법]
		b, e := json.Marshal(data)
		utils.rrHandleErr(err)
		fmt.Fprintf(writer, "%s", b)	// Marshalling을 하게되면 바이트 형식의 슬라이스가 리턴된다, 이를 스트링으로 변환해주어야 적절히 쓸 수 있다.
		// 여기까지 하면 "writer" 객체에 JSON 포매팅 된 data를 전달하여 API에 띄울 수 있게 된다. 다만 이때까지만 해도 콘텐츠 타입은 json이 아닌, text이다.
		// 따라서 맨위에 writer 객체의 header를 추가하여 content-type이 application/json 임을 알려야 한다.
	*/
	json.NewEncoder(writer).Encode(data)	// 간단히 처리할 수도 있다. >> Encode 메서드가 마샬링해주고 그 값을 NewEncoder 메서드를 통해 writer에 저장
}

func block(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)	// 리퀘스트에 담긴 변수 (패스 파라미터)를 vars 라는 인자로 던져주고
	hash := vars["hash"] // 리퀘스트에 담긴 변수 (패스 파라미터)를 담은 vars에서 우리가 필요한 "height"라는 값을 꺼내옴 // str > int 형변환
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(writer)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})			// 에러를 string으로 바꿈
	} else {
		encoder.Encode(block)
	}
}

func blocks(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		utils.HandleErr(json.NewEncoder(writer).Encode(blockchain.BlockChain().Blocks()))
	case "POST":
		var addBlockBody addBlockBody
		utils.HandleErr(json.NewDecoder(request.Body).Decode(&addBlockBody)) // 바로 위 라인에서 선언한 addBlockBody의 주소값에다가 디코딩을 해야함
		fmt.Printf("Request Body : %s", addBlockBody)
		blockchain.BlockChain().AddBlock(addBlockBody.Message)
		writer.WriteHeader(http.StatusCreated)
	}

}

// 다른 함수 시작 전 먼저 시작되도록 함 -> 여기서 정의한 담에 func Start 안에서 사용해준다.
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}