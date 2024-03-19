// generated by 'threeport-sdk codegen controller' - do not edit

package kubernetesruntime

import (
	"errors"
	"fmt"
	v0 "github.com/threeport/threeport/pkg/api/v0"
	client "github.com/threeport/threeport/pkg/client/v0"
	controller "github.com/threeport/threeport/pkg/controller/v0"
	notifications "github.com/threeport/threeport/pkg/notifications/v0"
	util "github.com/threeport/threeport/pkg/util/v0"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// KubernetesRuntimeDefinitionReconciler reconciles system state when a KubernetesRuntimeDefinition
// is created, updated or deleted.
func KubernetesRuntimeDefinitionReconciler(r *controller.Reconciler) {
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
				log.V(1).Info("kubernetes runtime definition reconciliation requeued with identical payload and fixed delay")
				continue
			}

			// decode the object that was sent in the notification
			var kubernetesRuntimeDefinition v0.KubernetesRuntimeDefinition
			if err := kubernetesRuntimeDefinition.DecodeNotifObject(notif.Object); err != nil {
				log.Error(err, "failed to marshal object map from consumed notification message")
				r.RequeueRaw(msg)
				log.V(1).Info("kubernetes runtime definition reconciliation requeued with identical payload and fixed delay")
				continue
			}
			log = log.WithValues("kubernetesRuntimeDefinitionID", kubernetesRuntimeDefinition.ID)

			// back off the requeue delay as needed
			requeueDelay := controller.SetRequeueDelay(
				notif.CreationTime,
			)

			// check for lock on object
			locked, ok := r.CheckLock(&kubernetesRuntimeDefinition)
			if locked || ok == false {
				r.Requeue(&kubernetesRuntimeDefinition, requeueDelay, msg)
				log.V(1).Info("kubernetes runtime definition reconciliation requeued")
				continue
			}

			// set up handler to unlock and requeue on termination signal
			go func() {
				select {
				case <-osSignals:
					log.V(1).Info("received termination signal, performing unlock and requeue of kubernetes runtime definition")
					r.UnlockAndRequeue(&kubernetesRuntimeDefinition, requeueDelay, lockReleased, msg)
				case <-lockReleased:
					log.V(1).Info("reached end of reconcile loop for kubernetes runtime definition, closing out signal handler")
				}
			}()

			// put a lock on the reconciliation of the created object
			if ok := r.Lock(&kubernetesRuntimeDefinition); !ok {
				r.Requeue(&kubernetesRuntimeDefinition, requeueDelay, msg)
				log.V(1).Info("kubernetes runtime definition reconciliation requeued")
				continue
			}

			// retrieve latest version of object
			latestKubernetesRuntimeDefinition, err := client.GetKubernetesRuntimeDefinitionByID(
				r.APIClient,
				r.APIServer,
				*kubernetesRuntimeDefinition.ID,
			)
			// check if error is 404 - if object no longer exists, no need to requeue
			if errors.Is(err, client.ErrObjectNotFound) {
				log.Info(fmt.Sprintf(
					"object with ID %d no longer exists - halting reconciliation",
					*kubernetesRuntimeDefinition.ID,
				))
				r.ReleaseLock(&kubernetesRuntimeDefinition, lockReleased, msg, true)
				continue
			}
			if err != nil {
				log.Error(err, "failed to get kubernetes runtime definition by ID from API")
				r.UnlockAndRequeue(&kubernetesRuntimeDefinition, requeueDelay, lockReleased, msg)
				continue
			}
			kubernetesRuntimeDefinition = *latestKubernetesRuntimeDefinition

			// determine which operation and act accordingly
			switch notif.Operation {
			case notifications.NotificationOperationCreated:
				if kubernetesRuntimeDefinition.DeletionScheduled != nil {
					log.Info("kubernetes runtime definition scheduled for deletion - skipping create")
					break
				}
				customRequeueDelay, err := kubernetesRuntimeDefinitionCreated(r, &kubernetesRuntimeDefinition, &log)
				if err != nil {
					log.Error(err, "failed to reconcile created kubernetes runtime definition object")
					r.UnlockAndRequeue(
						&kubernetesRuntimeDefinition,
						requeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				if customRequeueDelay != 0 {
					log.Info("create requeued for future reconciliation")
					r.UnlockAndRequeue(
						&kubernetesRuntimeDefinition,
						customRequeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
			case notifications.NotificationOperationUpdated:
				customRequeueDelay, err := kubernetesRuntimeDefinitionUpdated(r, &kubernetesRuntimeDefinition, &log)
				if err != nil {
					log.Error(err, "failed to reconcile updated kubernetes runtime definition object")
					r.UnlockAndRequeue(
						&kubernetesRuntimeDefinition,
						requeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				if customRequeueDelay != 0 {
					log.Info("update requeued for future reconciliation")
					r.UnlockAndRequeue(
						&kubernetesRuntimeDefinition,
						customRequeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
			case notifications.NotificationOperationDeleted:
				customRequeueDelay, err := kubernetesRuntimeDefinitionDeleted(r, &kubernetesRuntimeDefinition, &log)
				if err != nil {
					log.Error(err, "failed to reconcile deleted kubernetes runtime definition object")
					r.UnlockAndRequeue(
						&kubernetesRuntimeDefinition,
						requeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				if customRequeueDelay != 0 {
					log.Info("deletion requeued for future reconciliation")
					r.UnlockAndRequeue(
						&kubernetesRuntimeDefinition,
						customRequeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				deletionTimestamp := util.Ptr(time.Now().UTC())
				deletedKubernetesRuntimeDefinition := v0.KubernetesRuntimeDefinition{
					Common: v0.Common{ID: kubernetesRuntimeDefinition.ID},
					Reconciliation: v0.Reconciliation{
						DeletionAcknowledged: deletionTimestamp,
						DeletionConfirmed:    deletionTimestamp,
						Reconciled:           util.Ptr(true),
					},
				}
				if err != nil {
					log.Error(err, "failed to update kubernetes runtime definition to mark as reconciled")
					r.UnlockAndRequeue(&kubernetesRuntimeDefinition, requeueDelay, lockReleased, msg)
					continue
				}
				_, err = client.UpdateKubernetesRuntimeDefinition(
					r.APIClient,
					r.APIServer,
					&deletedKubernetesRuntimeDefinition,
				)
				if err != nil {
					log.Error(err, "failed to update kubernetes runtime definition to mark as deleted")
					r.UnlockAndRequeue(&kubernetesRuntimeDefinition, requeueDelay, lockReleased, msg)
					continue
				}
				_, err = client.DeleteKubernetesRuntimeDefinition(
					r.APIClient,
					r.APIServer,
					*kubernetesRuntimeDefinition.ID,
				)
				if err != nil {
					log.Error(err, "failed to delete kubernetes runtime definition")
					r.UnlockAndRequeue(&kubernetesRuntimeDefinition, requeueDelay, lockReleased, msg)
					continue
				}
			default:
				log.Error(
					errors.New("unrecognized notifcation operation"),
					"notification included an invalid operation",
				)
				r.UnlockAndRequeue(
					&kubernetesRuntimeDefinition,
					requeueDelay,
					lockReleased,
					msg,
				)
				continue

			}

			// set the object's Reconciled field to true if not deleted
			if notif.Operation != notifications.NotificationOperationDeleted {
				reconciledKubernetesRuntimeDefinition := v0.KubernetesRuntimeDefinition{
					Common:         v0.Common{ID: kubernetesRuntimeDefinition.ID},
					Reconciliation: v0.Reconciliation{Reconciled: util.Ptr(true)},
				}
				updatedKubernetesRuntimeDefinition, err := client.UpdateKubernetesRuntimeDefinition(
					r.APIClient,
					r.APIServer,
					&reconciledKubernetesRuntimeDefinition,
				)
				if err != nil {
					log.Error(err, "failed to update kubernetes runtime definition to mark as reconciled")
					r.UnlockAndRequeue(&kubernetesRuntimeDefinition, requeueDelay, lockReleased, msg)
					continue
				}
				log.V(1).Info(
					"kubernetes runtime definition marked as reconciled in API",
					"kubernetes runtime definitionName", updatedKubernetesRuntimeDefinition.Name,
				)
			}

			// release the lock on the reconciliation of the created object
			if ok := r.ReleaseLock(&kubernetesRuntimeDefinition, lockReleased, msg, true); !ok {
				log.Error(errors.New("kubernetes runtime definition remains locked - will unlock when TTL expires"), "")
			} else {
				log.V(1).Info("kubernetes runtime definition unlocked")
			}

			log.Info(fmt.Sprintf(
				"kubernetes runtime definition successfully reconciled for %s operation",
				notif.Operation,
			))
		}
	}

	r.Sub.Unsubscribe()
	reconcilerLog.Info("reconciler shutting down")
	r.ShutdownWait.Done()
}
