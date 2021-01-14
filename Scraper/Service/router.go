package Crawler

import (

	/*"os"
	"io"*/
	"net/http"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"	
	"github.com/chromedp/cdproto/input"	

	"Decorations/Scraper/Func/Worker"
)

var (
	R 				*gin.Engine
	w 				*worker.Worker = worker.InitWorker()
	allocCtx 		context.Context 
	Cancel 			context.CancelFunc
)


func init() {
	/*gin.SetMode(gin.ReleaseMode)

	var f *os.File
	if _, err := os.Stat("./Scraper.log"); err == nil {
		f,_ = os.OpenFile("./Scraper.log", os.O_RDWR|os.O_CREATE, 0755)
	} else if os.IsNotExist(err) {
		f,_ = os.Create("./Scraper.log")
	} else {
		f,_ = os.OpenFile("./Scraper.log", os.O_RDWR|os.O_CREATE, 0755)
	}

	gin.DefaultWriter = io.MultiWriter(f)*/


	R = gin.Default()

	R.GET("/", alive)
	facebook := R.Group("/facebook")
	{
		facebook.POST("/init"		, FacebookInit)
		facebook.POST("/crawler"	, FacebookScraper)
		facebook.POST("/replying"	, FacebookReplyer)
		facebook.POST("/announce"	, FacebookAnnounce)
		//facebook.POST("/whisper" 	, FacebookWhisper)
	}

	R.NoRoute(pageNotFound)
	R.NoMethod(pageNotFound)
}


func alive(c *gin.Context) {
	c.JSON( http.StatusOK, gin.H{
		"Msg":"ok",
		"StatusCode":"200",
	})	
}

func pageNotFound(c *gin.Context){
	c.JSON( http.StatusNotFound, gin.H{
		"Msg":"Error",
		"StatusCode":"404",
	})
}

func init() {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("disable-notifications", true),
		chromedp.Flag("mute-audio", false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.183 Safari/537.36`),
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:], options...)
	allocCtx, Cancel = chromedp.NewExecAllocator(context.Background(), opts...)
}






func MouseMoveNode(n *cdp.Node, opts ...chromedp.MouseOption) chromedp.MouseAction {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		boxes, err := dom.GetContentQuads().WithNodeID(n.NodeID).Do(ctx)
		if err != nil {
			return err
		}
		content := boxes[0]

		c := len(content)
		if c%2 != 0 || c < 1 {
			return chromedp.ErrInvalidDimensions
		}

		var x, y float64
		for i := 0; i < c; i += 2 {
			x += content[i]
			y += content[i+1]
		}
		x /= float64(c / 2)
		y /= float64(c / 2)

		return MouseMoveXY(x, y, opts...).Do(ctx)
	})
}

func MouseMoveXY(x, y float64, opts ...chromedp.MouseOption) chromedp.MouseAction {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		p := &input.DispatchMouseEventParams{
			Type:       input.MouseMoved,
			X:          x,
			Y:          y,
		}

		// apply opts
		for _, o := range opts {
			p = o(p)
		}

		return p.Do(ctx)
	})
}