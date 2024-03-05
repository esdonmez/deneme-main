package example

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Output: ", "")

	router := gin.New()

	pprof.Register(router)

	_ = router.Run("8082")
}
