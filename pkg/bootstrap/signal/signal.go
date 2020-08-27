package signal

import "context"

type App interface  {
	Stop () error
	GracefulStop(ctx context.Context) error
}