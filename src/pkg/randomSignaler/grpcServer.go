package randServer

import (
	"log"
	"math/rand"
	api "randsig/pkg/api"
	"time"

	"github.com/google/uuid"
)

type GRPCServer struct {
	api.UnimplementedRandomSignalerServer
}

func (s *GRPCServer) RandSignal(req *api.RandSignalRequest, stream api.RandomSignaler_RandSignalServer) error {
	sessionId := uuid.New()
	rand.Seed(time.Now().UnixNano())
	mean := rand.Intn(21) - 10
	std := 0.3 + rand.Float64()*(1.5-0.3)
	log.Println("Mean", mean, "STD", std)
	for {
		frequency := rand.NormFloat64()*std + float64(mean)
		err := stream.Send(
			&api.RandSignalResponse{
				SessionId:        sessionId.String(),
				Frequency:        frequency,
				CurrentTimestamp: time.Now().String(),
			})
		if err != nil {
			return err
		}
		time.Sleep(100000)
	}
}
