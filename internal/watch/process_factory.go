package watch

import (
	"context"
	"io"
	"time"

	"github.com/liupzmin/weewoe/serialize"

	"github.com/liupzmin/weewoe/internal/render"

	pb "github.com/liupzmin/weewoe/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type ProcessFactory struct {
	serialize.Serializable
	conn   *grpc.ClientConn
	client pb.StateClient
	stop   chan struct{}
}

func (p ProcessFactory) SendCommand(i int64) error {
	_, err := p.client.SendCommand(context.Background(), &pb.Command{
		Kind: "process",
		ID:   i,
	})
	return err
}

func NewProcessFactory(conn *grpc.ClientConn) *ProcessFactory {
	f := &ProcessFactory{
		Serializable: serialize.ProcessGob{},
		conn:         conn,
		stop:         make(chan struct{}),
	}
	f.client = pb.NewStateClient(conn)
	return f
}

func (p ProcessFactory) Terminate() {
	close(p.stop)
}

func (p ProcessFactory) Stream(cat string) <-chan render.Rows {
	ch := make(chan render.Rows)

	go func() {
	Repeat:
		ctx, cancel := context.WithCancel(context.Background())
		stream, err := p.client.Drain(ctx, &pb.Kind{
			Name: cat,
		})
		if err != nil {
			cancel()
			time.Sleep(5 * time.Second)
			log.Info().Msgf("GRPC Stream Call failed: %s, go to Repeat!", err)
			goto Repeat
		}

		log.Info().Msgf("GRPC Stream Call Successful!")

		go func() {
			<-p.stop
			_ = stream.CloseSend()
			cancel()
		}()

		for {
			data, err := stream.Recv()
			if err == io.EOF {
				log.Info().Msg("end of DrainProcessState, go to Repeat!")
				goto Repeat
			}
			if err != nil {
				if stream.Context().Err() == context.Canceled {
					log.Info().Msg("DrainProcessState canceled!")
					return
				}
				log.Error().Msgf("DrainProcessState recv failed: %s, go to Repeat!", err)

				goto Repeat

			}
			b, err := p.Decode(data.Content)
			if err != nil {
				log.Error().Msgf("Decode failed: %s", err)
				continue
			}
			log.Debug().Msgf("receive msg :%+v", b)
			ch <- b
		}
	}()

	return ch
}
