package runner

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
)

const (
	defaultPreExitSleetDuration = 5 * time.Second
)

type Starter interface {
	Start()
	Exit()
}

type Initier interface {
	Init() error
}

type Runner struct {
	services []Starter
	notifyCh chan os.Signal

	preExitSleepDuration time.Duration
}

func New() *Runner {
	return &Runner{
		notifyCh:             make(chan os.Signal),
		preExitSleepDuration: defaultPreExitSleetDuration,
	}
}

func (r *Runner) Add(srv Starter) {
	r.services = append(r.services, srv)
}

func (r *Runner) Run() {
	log.Print(color.GreenString("[runner]: starting %d services", len(r.services)))

	for _, srv := range r.services {
		if srv == nil {
			log.Panicf("%+v service is nil", srv)
		}
		if initier, ok := srv.(Initier); ok {
			if err := initier.Init(); err != nil {
				log.Panicf("failed to init service: %s", err.Error())
			}
		}

		go srv.Start()
	}

	signal.Notify(r.notifyCh, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-r.notifyCh:
			r.Exit()

			log.Print(color.YellowString("[runner] waiting %v before exiting", r.preExitSleepDuration))
			t := time.NewTimer(r.preExitSleepDuration)
			<-t.C
			os.Exit(0)
		}
	}
}

func (r *Runner) Exit() {
	log.Print(color.YellowString("[runner] exiting"))
	for _, srv := range r.services {
		srv.Exit()
	}
}

func (r *Runner) SetPreExitSleepDuration(d time.Duration) {
	r.preExitSleepDuration = d
}
