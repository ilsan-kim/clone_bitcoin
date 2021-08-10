package cli

import (
	"flag"
	"fmt"
	"github.com/ddhyun93/seancoin/explorer"
	"github.com/ddhyun93/seancoin/rest"
	"os"
)

func usage() {
	fmt.Printf("Welcome to Seancoin CLI environment\n\n")
	fmt.Printf("Please use the following flags:\n")
	fmt.Printf("-port=4000: Set port of the server\n")
	fmt.Printf("-mode=html: Choose between 'html' and 'rest * if you select 'multi'mode, you dont need to add '-port' flag '\n\n")
	os.Exit(0)
}

func Start() {
	mode := flag.String("mode", "multi", "Choose between 'html' and 'rest'")
	port := flag.Int("port", 4000, "Set port of the server")
	flag.Parse()

	switch *mode {
	case "rest":
		// start rest api
		rest.Start(*port)
		fmt.Println(*port, *mode)	// port 와 mode 모두 포인터 이기 때문에, look through 하기 위해서 *를 써줘야함
	case "html":
		// start html explorer
		explorer.Start(*port)
		fmt.Println(*port, *mode)	// port 와 mode 모두 포인터 이기 때문에, look through 하기 위해서 *를 써줘야함
	case "multi":
		// start html and rest api at one time
		go explorer.Start(5000)
		rest.Start(4000)
	default:
		usage()
	}
}