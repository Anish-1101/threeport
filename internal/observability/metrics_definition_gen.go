// generated by 'threeport-sdk gen' for controller scaffolding - do not edit

package observability

import (
	"errors"
	"fmt"
	v0 "github.com/threeport/threeport/pkg/api/v0"
	v1 "github.com/threeport/threeport/pkg/api/v1"
	client "github.com/threeport/threeport/pkg/client/v0"
	controller "github.com/threeport/threeport/pkg/controller/v0"
	notifications "github.com/threeport/threeport/pkg/notifications/v0"
	util "github.com/threeport/threeport/pkg/util/v0"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// MetricsDefinitionReconciler reconciles system state when a MetricsDefinition
// is created, updated or deleted.
func MetricsDefinitionReconciler(r *controller.Reconciler) {
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
				log.V(1).Info("metrics definition reconciliation requeued with identical payload and fixed delay")
				continue
			}

			// decode the object that was sent in the notification
			var metricsDefinition v0.MetricsDefinition
			if err := metricsDefinition.DecodeNotifObject(notif.Object); err != nil {
				log.Error(err, "failed to marshal object map from consumed notification message")
				r.RequeueRaw(msg)
				log.V(1).Info("metrics definition reconciliation requeued with identical payload and fixed delay")
				continue
			}
			log = log.WithValues("metricsDefinitionID", metricsDefinition.ID)

			// back off the requeue delay as needed
			requeueDelay := controller.SetRequeueDelay(
				notif.CreationTime,
			)

			// check for lock on object
			locked, ok := r.CheckLock(&metricsDefinition)
			if locked || ok == false {
				r.Requeue(&metricsDefinition, requeueDelay, msg)
				log.V(1).Info("metrics definition reconciliation requeued")
				continue
			}

			// set up handler to unlock and requeue on termination signal
			go func() {
				select {
				case <-osSignals:
					log.V(1).Info("received termination signal, performing unlock and requeue of metrics definition")
					r.UnlockAndRequeue(&metricsDefinition, requeueDelay, lockReleased, msg)
				case <-lockReleased:
					log.V(1).Info("reached end of reconcile loop for metrics definition, closing out signal handler")
				}
			}()

			// put a lock on the reconciliation of the created object
			if ok := r.Lock(&metricsDefinition); !ok {
				r.Requeue(&metricsDefinition, requeueDelay, msg)
				log.V(1).Info("metrics definition reconciliation requeued")
				continue
			}

			// retrieve latest version of object
			latestMetricsDefinition, err := client.GetMetricsDefinitionByID(
				r.APIClient,
				r.APIServer,
				*metricsDefinition.ID,
			)
			// check if error is 404 - if object no longer exists, no need to requeue
			if errors.Is(err, client.ErrObjectNotFound) {
				log.Info(fmt.Sprintf(
					"object with ID %d no longer exists - halting reconciliation",
					*metricsDefinition.ID,
				))
				r.ReleaseLock(&metricsDefinition, lockReleased, msg, true)
				continue
			}
			if err != nil {
				log.Error(err, "failed to get metrics definition by ID from API")
				r.UnlockAndRequeue(&metricsDefinition, requeueDelay, lockReleased, msg)
				continue
			}
			metricsDefinition = *latestMetricsDefinition

			// determine which operation and act accordingly
			switch notif.Operation {
			case notifications.NotificationOperationCreated:
				if metricsDefinition.DeletionScheduled != nil {
					log.Info("metrics definition scheduled for deletion - skipping create")
					break
				}
				customRequeueDelay, err := metricsDefinitionCreated(r, &metricsDefinition, &log)
				if err != nil {
					errorMsg := "failed to reconcile created metrics definition object"
					log.Error(err, errorMsg)
					r.EventsRecorder.HandleEventOverride(
						&v1.Event{
							Note:   util.Ptr(errorMsg),
							Reason: util.Ptr("MetricsDefinitionNotCreated"),
							Type:   util.Ptr("Normal"),
						},
						metricsDefinition.ID,
						err,
						&log,
					)
					r.UnlockAndRequeue(
						&metricsDefinition,
						requeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				if customRequeueDelay != 0 {
					log.Info("create requeued for future reconciliation")
					r.UnlockAndRequeue(
						&metricsDefinition,
						customRequeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
			case notifications.NotificationOperationUpdated:
				customRequeueDelay, err := metricsDefinitionUpdated(r, &metricsDefinition, &log)
				if err != nil {
					errorMsg := "failed to reconcile updated metrics definition object"
					log.Error(err, errorMsg)
					r.EventsRecorder.HandleEventOverride(
						&v1.Event{
							Note:   util.Ptr(errorMsg),
							Reason: util.Ptr("MetricsDefinitionNotUpdated"),
							Type:   util.Ptr("Normal"),
						},
						metricsDefinition.ID,
						err,
						&log,
					)
					r.UnlockAndRequeue(
						&metricsDefinition,
						requeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				if customRequeueDelay != 0 {
					log.Info("update requeued for future reconciliation")
					r.UnlockAndRequeue(
						&metricsDefinition,
						customRequeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
			case notifications.NotificationOperationDeleted:
				customRequeueDelay, err := metricsDefinitionDeleted(r, &metricsDefinition, &log)
				if err != nil {
					errorMsg := "failed to reconcile deleted metrics definition object"
					log.Error(err, errorMsg)
					r.EventsRecorder.HandleEventOverride(
						&v1.Event{
							Note:   util.Ptr(errorMsg),
							Reason: util.Ptr("MetricsDefinitionNotUpdated"),
							Type:   util.Ptr("Normal"),
						},
						metricsDefinition.ID,
						err,
						&log,
					)
					r.UnlockAndRequeue(
						&metricsDefinition,
						requeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				if customRequeueDelay != 0 {
					log.Info("deletion requeued for future reconciliation")
					r.UnlockAndRequeue(
						&metricsDefinition,
						customRequeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				deletionTimestamp := util.Ptr(time.Now().UTC())
				deletedMetricsDefinition := v0.MetricsDefinition{
					Common: v0.Common{ID: metricsDefinition.ID},
					Reconciliation: v0.Reconciliation{
						DeletionAcknowledged: deletionTimestamp,
						DeletionConfirmed:    deletionTimestamp,
						Reconciled:           util.Ptr(true),
					},
				}
				if err != nil {
					log.Error(err, "failed to update metrics definition to mark as reconciled")
					r.UnlockAndRequeue(&metricsDefinition, requeueDelay, lockReleased, msg)
					continue
				}
				_, err = client.UpdateMetricsDefinition(
					r.APIClient,
					r.APIServer,
					&deletedMetricsDefinition,
				)
				if err != nil {
					log.Error(err, "failed to update metrics definition to mark as deleted")
					r.UnlockAndRequeue(&metricsDefinition, requeueDelay, lockReleased, msg)
					continue
				}
				_, err = client.DeleteMetricsDefinition(
					r.APIClient,
					r.APIServer,
					*metricsDefinition.ID,
				)
				if err != nil {
					log.Error(err, "failed to delete metrics definition")
					r.UnlockAndRequeue(&metricsDefinition, requeueDelay, lockReleased, msg)
					continue
				}
			default:
				log.Error(
					errors.New("unrecognized notifcation operation"),
					"notification included an invalid operation",
				)
				r.UnlockAndRequeue(
					&metricsDefinition,
					requeueDelay,
					lockReleased,
					msg,
				)
				continue

			}

			// set the object's Reconciled field to true if not deleted
			if notif.Operation != notifications.NotificationOperationDeleted {
				reconciledMetricsDefinition := v0.MetricsDefinition{
					Common:         v0.Common{ID: metricsDefinition.ID},
					Reconciliation: v0.Reconciliation{Reconciled: util.Ptr(true)},
				}
				updatedMetricsDefinition, err := client.UpdateMetricsDefinition(
					r.APIClient,
					r.APIServer,
					&reconciledMetricsDefinition,
				)
				if err != nil {
					log.Error(err, "failed to update metrics definition to mark as reconciled")
					r.UnlockAndRequeue(&metricsDefinition, requeueDelay, lockReleased, msg)
					continue
				}
				log.V(1).Info(
					"metrics definition marked as reconciled in API",
					"metrics definitionName", updatedMetricsDefinition.Name,
				)
			}

			// release the lock on the reconciliation of the created object
			if ok := r.ReleaseLock(&metricsDefinition, lockReleased, msg, true); !ok {
				log.Error(errors.New("metrics definition remains locked - will unlock when TTL expires"), "")
			} else {
				log.V(1).Info("metrics definition unlocked")
			}

			successMsg := fmt.Sprintf(
				"metrics definition successfully reconciled for %s operation",
				strings.ToLower(string(notif.Operation)),
			)
			if err := r.EventsRecorder.RecordEvent(
				&v1.Event{
					Note:   util.Ptr(successMsg),
					Reason: util.Ptr("MetricsDefinitionSuccessfullyReconciled"),
					Type:   util.Ptr("Normal"),
				},
				metricsDefinition.ID,
			); err != nil {
				log.Error(err, "failed to record event for successful metrics definition reconciliation")
			}
		}
	}

	r.Sub.Unsubscribe()
	reconcilerLog.Info("reconciler shutting down")
	r.ShutdownWait.Done()
}
