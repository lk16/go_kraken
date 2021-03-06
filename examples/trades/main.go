package main

import (
	"log"
	"time"

	ws "github.com/aopoltorzhicky/go_kraken/websocket"
)

func main() {
	c := ws.New()
	err := c.Connect()
	if err != nil {
		log.Fatal("Error connecting to web socket : ", err)
	}

	pairs := []string{ws.BTCUSD}

	// subscribe to BTCUSD trades
	err = c.SubscribeTrades(pairs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		time.Sleep(time.Second * 30)
		log.Print("Unsubsribing...")
		err = c.Unsubscribe(ws.ChanTrades, pairs)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Success!")
		c.Close()
	}()

	for obj := range c.Listen() {
		switch obj.(type) {
		case error:
			log.Printf("channel closed: %s", obj)
		case ws.DataUpdate:
			data := obj.(ws.DataUpdate)
			for _, trade := range data.Data.([]ws.TradeUpdate) {
				log.Print("----------------")
				log.Printf("Price: %f", trade.Price)
				log.Printf("Volume: %f", trade.Volume)
				log.Printf("Time: %f", trade.Time)
				log.Printf("Pair: %s", trade.Pair)
				log.Printf("Order type: %s", trade.OrderType)
				log.Printf("Side: %s", trade.Side)
				log.Printf("Misc: %s", trade.Misc)
			}
		}

	}
}
