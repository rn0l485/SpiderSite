package Crawler

import (
	"github.com/gin-gonic/gin"

	"Decorations/Func/Worker"
	"Decorations/Scraper/Service/Models"
)

var (
	R 				*gin.Engine
	CTX 			chromedp.Context
	Cancel 			func()
	w 				*worker.Worker = worker.InitWorker()
	setting 		ScraperModels.Setting
)



func init() {
	R = gin.Default()

	R.GET("/", alive)
	facebook := R.Group("/facebook")
	{
		facebook.GET("/crawler", FBCrawler)
	}
	


	R.NoRoute(pageNotFound)
	R.NoMethod(pageNotFound)
}


func alive(c *gin.Context) {
	c.JSON( http.StatusNotFound, gin.H{
		"Msg":"Well, nothing can help you.",
		"StatusCode":"200",
	})	
}

func pageNotFound(c *gin.Context){
	c.JSON( http.StatusNotFound, gin.H{
		"Msg":"Path error",
		"StatusCode":"404",
	})
}

func init() {
	options := []chromedp.ExecAllocatorOption{
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.183 Safari/537.36`),
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:], options...)
	allocCtx, Cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	// create chrome instance
	CTX, Cancel = chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(logger), // todo, can set the request log
	)
}

func logger(format string, v ...interface{}) {
	resp, err := w.Post(config.MongoDBApi+"/v1/add", false, gin.H{
		"DataBaseName"	: "Spider",
		"CollectionName" : "Log",
		"Record": {
			"CreateTime": time.Now().Format("2006-01-02 15:04:05"),
			"Context": fmt.Sprintf(format, v...),
			"Service": "Crawler",
		},
	})
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