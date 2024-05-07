// generated by 'threeport-sdk gen' for controller scaffolding - do not edit

package observability

import (
	"errors"
	"fmt"
	v0 "github.com/threeport/threeport/pkg/api/v0"
	v1 "github.com/threeport/threeport/pkg/api/v1"
	client "github.com/threeport/threeport/pkg/client/v0"
	controller "github.com/threeport/threeport/pkg/controller/v0"
	event "github.com/threeport/threeport/pkg/event/v0"
	notifications "github.com/threeport/threeport/pkg/notifications/v0"
	util "github.com/threeport/threeport/pkg/util/v0"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// MetricsInstanceReconciler reconciles system state when a MetricsInstance
// is created, updated or deleted.
func MetricsInstanceReconciler(r *controller.Reconciler) {
	r.ShutdownWait.Add(1)
	reconcilerLog := r.Log.WithValues("reconcilerName", r.Name)
	reconcilerLog.Info("reconciler started")
	shutdown := false

	// create a channel to receive OS signals
	osSignals := make(chan os.Signal, 1)
	lockReleased := make(chan bool, 1)

	// register the os signals channel to receive SIGINT and SIGTERM signals
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	for {
		// create a fresh log object per reconciliation loop so we don't
		// accumulate values across multiple loops
		log := r.Log.WithValues("reconcilerName", r.Name)

		if shutdown {
			break
		}

		// check for shutdown instruction
		select {
		case <-r.Shutdown:
			shutdown = true
		default:
			// pull message off queue
			msg := r.PullMessage()
			if msg == nil {
				continue
			}

			// consume message data to capture notification from API
			notif, err := notifications.ConsumeMessage(msg.Data)
			if err != nil {
				log.Error(
					err, "failed to consume message data from NATS",
					"msgData", string(msg.Data),
				)
				r.RequeueRaw(msg)
				log.V(1).Info("metrics instance reconciliation requeued with identical payload and fixed delay")
				continue
			}

			// decode the object that was sent in the notification
			var metricsInstance v0.MetricsInstance
			if err := metricsInstance.DecodeNotifObject(notif.Object); err != nil {
				log.Error(err, "failed to marshal object map from consumed notification message")
				r.RequeueRaw(msg)
				log.V(1).Info("metrics instance reconciliation requeued with identical payload and fixed delay")
				continue
			}
			log = log.WithValues("metricsInstanceID", metricsInstance.ID)

			// back off the requeue delay as needed
			requeueDelay := controller.SetRequeueDelay(
				notif.CreationTime,
			)

			// check for lock on object
			locked, ok := r.CheckLock(&metricsInstance)
			if locked || ok == false {
				r.Requeue(&metricsInstance, requeueDelay, msg)
				log.V(1).Info("metrics instance reconciliation requeued")
				continue
			}

			// set up handler to unlock and requeue on termination signal
			go func() {
				select {
				case <-osSignals:
					log.V(1).Info("received termination signal, performing unlock and requeue of metrics instance")
					r.UnlockAndRequeue(&metricsInstance, requeueDelay, lockReleased, msg)
				case <-lockReleased:
					log.V(1).Info("reached end of reconcile loop for metrics instance, closing out signal handler")
				}
			}()

			// put a lock on the reconciliation of the created object
			if ok := r.Lock(&metricsInstance); !ok {
				r.Requeue(&metricsInstance, requeueDelay, msg)
				log.V(1).Info("metrics instance reconciliation requeued")
				continue
			}

			// retrieve latest version of object
			latestMetricsInstance, err := client.GetMetricsInstanceByID(
				r.APIClient,
				r.APIServer,
				*metricsInstance.ID,
			)
			// check if error is 404 - if object no longer exists, no need to requeue
			if errors.Is(err, client.ErrObjectNotFound) {
				log.Info(fmt.Sprintf(
					"object with ID %d no longer exists - halting reconciliation",
					*metricsInstance.ID,
				))
				r.ReleaseLock(&metricsInstance, lockReleased, msg, true)
				continue
			}
			if err != nil {
				log.Error(err, "failed to get metrics instance by ID from API")
				r.UnlockAndRequeue(&metricsInstance, requeueDelay, lockReleased, msg)
				continue
			}
			metricsInstance = *latestMetricsInstance

			// determine which operation and act accordingly
			switch notif.Operation {
			case notifications.NotificationOperationCreated:
				if metricsInstance.DeletionScheduled != nil {
					log.Info("metrics instance scheduled for deletion - skipping create")
					break
				}
				customRequeueDelay, err := metricsInstanceCreated(r, &metricsInstance, &log)
				if err != nil {
					errorMsg := "failed to reconcile created metrics instance object"
					log.Error(err, errorMsg)
					r.EventsRecorder.HandleEventOverride(
						&v1.Event{
							Note:   util.Ptr(errorMsg),
							Reason: util.Ptr(event.ReasonFailedCreate),
							Type:   util.Ptr(event.TypeNormal),
						},
						metricsInstance.ID,
						err,
						&log,
					)
					r.UnlockAndRequeue(
						&metricsInstance,
						requeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				if customRequeueDelay != 0 {
					log.Info("create requeued for future reconciliation")
					r.UnlockAndRequeue(
						&metricsInstance,
						customRequeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
			case notifications.NotificationOperationUpdated:
				customRequeueDelay, err := metricsInstanceUpdated(r, &metricsInstance, &log)
				if err != nil {
					errorMsg := "failed to reconcile updated metrics instance object"
					log.Error(err, errorMsg)
					r.EventsRecorder.HandleEventOverride(
						&v1.Event{
							Note:   util.Ptr(errorMsg),
							Reason: util.Ptr(event.ReasonFailedUpdate),
							Type:   util.Ptr(event.TypeNormal),
						},
						metricsInstance.ID,
						err,
						&log,
					)
					r.UnlockAndRequeue(
						&metricsInstance,
						requeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				if customRequeueDelay != 0 {
					log.Info("update requeued for future reconciliation")
					r.UnlockAndRequeue(
						&metricsInstance,
						customRequeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
			case notifications.NotificationOperationDeleted:
				customRequeueDelay, err := metricsInstanceDeleted(r, &metricsInstance, &log)
				if err != nil {
					errorMsg := "failed to reconcile deleted metrics instance object"
					log.Error(err, errorMsg)
					r.EventsRecorder.HandleEventOverride(
						&v1.Event{
							Note:   util.Ptr(errorMsg),
							Reason: util.Ptr(event.ReasonFailedUpdate),
							Type:   util.Ptr(event.TypeNormal),
						},
						metricsInstance.ID,
						err,
						&log,
					)
					r.UnlockAndRequeue(
						&metricsInstance,
						requeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				if customRequeueDelay != 0 {
					log.Info("deletion requeued for future reconciliation")
					r.UnlockAndRequeue(
						&metricsInstance,
						customRequeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				deletionTimestamp := util.Ptr(time.Now().UTC())
				deletedMetricsInstance := v0.MetricsInstance{
					Common: v0.Common{ID: metricsInstance.ID},
					Reconciliation: v0.Reconciliation{
						DeletionAcknowledged: deletionTimestamp,
						DeletionConfirmed:    deletionTimestamp,
						Reconciled:           util.Ptr(true),
					},
				}
				_, err = client.UpdateMetricsInstance(
					r.APIClient,
					r.APIServer,
					&deletedMetricsInstance,
				)
				if err != nil {
					log.Error(err, "failed to update metrics instance to mark as deleted")
					r.UnlockAndRequeue(&metricsInstance, requeueDelay, lockReleased, msg)
					continue
				}
				_, err = client.DeleteMetricsInstance(
					r.APIClient,
					r.APIServer,
					*metricsInstance.ID,
				)
				if err != nil {
					log.Error(err, "failed to delete metrics instance")
					r.UnlockAndRequeue(&metricsInstance, requeueDelay, lockReleased, msg)
					continue
				}
			default:
				log.Error(
					errors.New("unrecognized notifcation operation"),
					"notification included an invalid operation",
				)
				r.UnlockAndRequeue(
					&metricsInstance,
					requeueDelay,
					lockReleased,
					msg,
				)
				continue

			}

			// set the object's Reconciled field to true if not deleted
			if notif.Operation != notifications.NotificationOperationDeleted {
				reconciledMetricsInstance := v0.MetricsInstance{
					Common:         v0.Common{ID: metricsInstance.ID},
					Reconciliation: v0.Reconciliation{Reconciled: util.Ptr(true)},
				}
				updatedMetricsInstance, err := client.UpdateMetricsInstance(
					r.APIClient,
					r.APIServer,
					&reconciledMetricsInstance,
				)
				if err != nil {
					log.Error(err, "failed to update metrics instance to mark as reconciled")
					r.UnlockAndRequeue(&metricsInstance, requeueDelay, lockReleased, msg)
					continue
				}
				log.V(1).Info(
					"metrics instance marked as reconciled in API",
					"metrics instanceName", updatedMetricsInstance.Name,
				)
			}

			// release the lock on the reconciliation of the created object
			if ok := r.ReleaseLock(&metricsInstance, lockReleased, msg, true); !ok {
				log.Error(errors.New("metrics instance remains locked - will unlock when TTL expires"), "")
			} else {
				log.V(1).Info("metrics instance unlocked")
			}

			successMsg := fmt.Sprintf(
				"metrics instance successfully reconciled for %s operation",
				strings.ToLower(string(notif.Operation)),
			)
			if err := r.EventsRecorder.RecordEvent(
				&v1.Event{
					Note:   util.Ptr(successMsg),
					Reason: util.Ptr(event.GetSuccessReasonForOperation(notif.Operation)),
					Type:   util.Ptr(event.TypeNormal),
				},
				metricsInstance.ID,
			); err != nil {
				log.Error(err, "failed to record event for successful metrics instance reconciliation")
			}
		}
	}

	r.Sub.Unsubscribe()
	reconcilerLog.Info("reconciler shutting down")
	r.ShutdownWait.Done()
}
