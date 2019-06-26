package connection

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/scrape"
	"github.com/stackrox/rox/central/sensor/networkpolicies"
	"github.com/stackrox/rox/central/sensor/service/pipeline"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/reflectutils"
	"github.com/stackrox/rox/pkg/sac"
	"golang.org/x/time/rate"
	"google.golang.org/grpc/metadata"
)

var (
	log = logging.LoggerForModule()
)

type sensorConnection struct {
	clusterID           string
	stopSig, stoppedSig concurrency.ErrorSignal

	sendC chan *central.MsgToSensor

	scrapeCtrl          scrape.Controller
	networkPoliciesCtrl networkpolicies.Controller

	sensorEventHandler *sensorEventHandler

	queues map[string]*dedupingQueue

	eventPipeline pipeline.ClusterPipeline

	clusterMgr               ClusterManager
	checkInRecordRateLimiter *rate.Limiter
}

func newConnection(ctx context.Context, clusterID string, pf pipeline.Factory, clusterMgr ClusterManager) (*sensorConnection, error) {
	eventPipeline, err := pf.PipelineForCluster(ctx, clusterID)
	if err != nil {
		return nil, errors.Wrap(err, "creating event pipeline")
	}

	conn := &sensorConnection{
		stopSig:       concurrency.NewErrorSignal(),
		stoppedSig:    concurrency.NewErrorSignal(),
		sendC:         make(chan *central.MsgToSensor),
		eventPipeline: eventPipeline,
		queues:        make(map[string]*dedupingQueue),

		clusterID:  clusterID,
		clusterMgr: clusterMgr,

		checkInRecordRateLimiter: rate.NewLimiter(rate.Every(10*time.Second), 1),
	}

	// Need a reference to conn for injector
	conn.sensorEventHandler = newSensorEventHandler(eventPipeline, conn, &conn.stopSig)

	conn.scrapeCtrl = scrape.NewController(conn, &conn.stopSig)
	conn.networkPoliciesCtrl = networkpolicies.NewController(conn, &conn.stopSig)
	return conn, nil
}

func (c *sensorConnection) Terminate(err error) bool {
	return c.stopSig.SignalWithError(err)
}

func (c *sensorConnection) Stopped() concurrency.ReadOnlyErrorSignal {
	return &c.stoppedSig
}

// Record the check-in if the rate limiter allows it.
func (c *sensorConnection) recordCheckInRateLimited(ctx context.Context) {
	if c.checkInRecordRateLimiter.Allow() {
		err := c.clusterMgr.UpdateClusterContactTime(ctx, c.clusterID, time.Now())
		if err != nil {
			log.Warnf("Could not record cluster contact: %v", err)
		}
	}
}

func (c *sensorConnection) multiplexedPush(ctx context.Context, msg *central.MsgFromSensor) {
	typ := reflectutils.Type(msg.Msg)
	queue := c.queues[typ]
	if queue == nil {
		queue = newDedupingQueue()
		c.queues[typ] = queue
		go c.handleMessages(ctx, queue)
	}
	queue.push(msg)
}

func (c *sensorConnection) runRecv(ctx context.Context, server central.SensorService_CommunicateServer) {
	for !c.stopSig.IsDone() {
		msg, err := server.Recv()
		if err != nil {
			c.stopSig.SignalWithError(errors.Wrap(err, "recv error"))
			return
		}
		c.recordCheckInRateLimited(ctx)

		c.multiplexedPush(ctx, msg)
	}
}

func (c *sensorConnection) handleMessages(ctx context.Context, queue *dedupingQueue) {
	for msg := queue.pullBlocking(&c.stopSig); msg != nil; msg = queue.pullBlocking(&c.stopSig) {
		if err := c.handleMessage(ctx, msg); err != nil {
			log.Errorf("Error handling sensor message: %v", err)
		}
	}
	c.eventPipeline.OnFinish(c.clusterID)
	c.stoppedSig.SignalWithError(c.stopSig.Err())
}

func (c *sensorConnection) runSend(server central.SensorService_CommunicateServer) {
	for !c.stopSig.IsDone() {
		select {
		case <-c.stopSig.Done():
			return
		case <-server.Context().Done():
			c.stopSig.SignalWithError(errors.Wrap(server.Context().Err(), "context error"))
			return
		case msg := <-c.sendC:
			if err := server.Send(msg); err != nil {
				c.stopSig.SignalWithError(errors.Wrap(err, "send error"))
				return
			}
		}
	}
}

func (c *sensorConnection) Scrapes() scrape.Controller {
	return c.scrapeCtrl
}

func (c *sensorConnection) InjectMessageIntoQueue(msg *central.MsgFromSensor) {
	c.multiplexedPush(sac.WithAllAccess(context.Background()), msg)
}

func (c *sensorConnection) NetworkPolicies() networkpolicies.Controller {
	return c.networkPoliciesCtrl
}

func (c *sensorConnection) InjectMessage(ctx concurrency.Waitable, msg *central.MsgToSensor) error {
	select {
	case c.sendC <- msg:
		return nil
	case <-ctx.Done():
		return errors.New("context aborted")
	case <-c.stopSig.Done():
		return errors.Wrap(c.stopSig.Err(), "could not send message as sensor connection was stopped")
	}
}

func (c *sensorConnection) handleMessage(ctx context.Context, msg *central.MsgFromSensor) error {
	switch m := msg.Msg.(type) {
	case *central.MsgFromSensor_ScrapeUpdate:
		return c.scrapeCtrl.ProcessScrapeUpdate(m.ScrapeUpdate)
	case *central.MsgFromSensor_NetworkPoliciesResponse:
		return c.networkPoliciesCtrl.ProcessNetworkPoliciesResponse(m.NetworkPoliciesResponse)
	case *central.MsgFromSensor_Event:
		// Special case the reprocess deployment because its fields are already set
		if msg.GetEvent().GetReprocessDeployment() != nil {
			c.sensorEventHandler.addMultiplexed(ctx, msg)
			return nil
		}
		// Only dedupe on non-creates
		if msg.GetEvent().GetAction() != central.ResourceAction_CREATE_RESOURCE {
			msg.DedupeKey = msg.GetEvent().GetId()
		}
		// Set the hash key for all values
		msg.HashKey = msg.GetEvent().GetId()

		c.sensorEventHandler.addMultiplexed(ctx, msg)
		return nil
	}
	return c.eventPipeline.Run(ctx, msg, c)
}

func (c *sensorConnection) getClusterConfigMsg(ctx context.Context) (*central.MsgToSensor, error) {
	cluster, exists, err := c.clusterMgr.GetCluster(ctx, c.clusterID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.Errorf("could not pull config for cluster %q because it does not exist", c.clusterID)
	}
	return &central.MsgToSensor{
		Msg: &central.MsgToSensor_ClusterConfig{
			ClusterConfig: &central.ClusterConfig{
				Config: cluster.GetDynamicConfig(),
			},
		},
	}, nil
}

func (c *sensorConnection) Run(ctx context.Context, server central.SensorService_CommunicateServer) error {
	if err := server.SendHeader(metadata.MD{}); err != nil {
		return errors.Wrap(err, "sending initial metadata")
	}

	// Synchronously send the config to ensure syncing before Sensor marks the connection as Central reachable
	msg, err := c.getClusterConfigMsg(ctx)
	if err != nil {
		return errors.Wrapf(err, "unable to get cluster config for %q", c.clusterID)
	}

	if err := server.Send(msg); err != nil {
		return errors.Wrapf(err, "unable to sync config to cluster %q", c.clusterID)
	}

	go c.runSend(server)

	c.runRecv(ctx, server)
	return c.stopSig.Err()
}

func (c *sensorConnection) ClusterID() string {
	return c.clusterID
}
