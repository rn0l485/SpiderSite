package Crawler

import (
	"github.com/gin-gonic/gin"

	"Decorations/Func/Worker"
	"Decorations/Scraper/Service/Models"
)

var (
	R 				*gin.Engine
	w 				*worker.Worker = worker.InitWorker()
)



func init() {
	R = gin.Default()

	R.GET("/", alive)
	facebook := R.Group("/facebook")
	{
		facebook.POST("/crawler", PersonalCrawler)
		facebook.GET("/replyer", FBReplyer)
	}

	R.NoRoute(pageNotFound)
	R.NoMethod(pageNotFound)
}


func alive(c *gin.Context) {
	c.JSON( http.StatusNotFound, gin.H{
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