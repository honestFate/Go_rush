package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	api "randsig/pkg/api"
	p "randsig/pkg/postgres"
	"strconv"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getMean(pool []float64) float64 {
	var mean float64
	for i, val := range pool {
		mean += (val - mean) / float64((i + 1))
	}
	return mean
}

func getSD(pool []float64) float64 {
	mean := getMean(pool)
	var dispersion float64
	for _, val := range pool {
		dispersion += math.Pow(float64(val-mean), 2)
	}
	dispersion /= float64(len(pool))
	return math.Sqrt(dispersion)
}

var messagePool = sync.Pool{
	New: func() interface{} { return []float64{} },
}

func receiveData(stream api.RandomSignaler_RandSignalClient, k float64) {
	pool := messagePool.Get().([]float64)
	for i := 0; i < 1000; i++ {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		pool = append(pool, message.Frequency)
	}
	messagePool.Put(&pool)
	mean := getMean(pool)
	sd := getSD(pool)
	fmt.Println("Mean:", mean, "/ sd:", sd)
	conn := p.NewDBConn()
	count := 0
	for i := 0; i < 1000; i++ {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		div := message.Frequency - mean
		if math.Abs(div) > (sd * k) {
			err = p.InsertDB(conn, &p.Message{
				Session_id:        message.SessionId,
				Frequency:         strconv.FormatFloat(message.Frequency, 'f', 6, 64),
				Current_timestamp: message.CurrentTimestamp,
			})
			if err != nil {
				fmt.Println("cant insert")
				log.Fatal(err)
			}
			count++
		}
	}
	if err := conn.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Total ", count)
}

func flagParse() float64 {
	k := flag.Float64("k", 1.00, "[0; +fmax]")
	flag.Parse()
	if *k < 0.00 {
		flag.PrintDefaults()
		log.Fatal("Wrong usage")
	}
	return *k
}

func main() {
	k := flagParse()
	connection, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	client := api.NewRandomSignalerClient(connection)
	stream, err := client.RandSignal(context.Background(), &api.RandSignalRequest{})
	if err != nil {
		log.Fatal(err)
	}
	receiveData(stream, k)
}
