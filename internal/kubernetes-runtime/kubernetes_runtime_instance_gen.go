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

// KubernetesRuntimeInstanceReconciler reconciles system state when a KubernetesRuntimeInstance
// is created, updated or deleted.
func KubernetesRuntimeInstanceReconciler(r *controller.Reconciler) {
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
				log.V(1).Info("kubernetes runtime instance reconciliation requeued with identical payload and fixed delay")
				continue
			}

			// decode the object that was sent in the notification
			var kubernetesRuntimeInstance v0.KubernetesRuntimeInstance
			if err := kubernetesRuntimeInstance.DecodeNotifObject(notif.Object); err != nil {
				log.Error(err, "failed to marshal object map from consumed notification message")
				r.RequeueRaw(msg)
				log.V(1).Info("kubernetes runtime instance reconciliation requeued with identical payload and fixed delay")
				continue
			}
			log = log.WithValues("kubernetesRuntimeInstanceID", kubernetesRuntimeInstance.ID)

			// back off the requeue delay as needed
			requeueDelay := controller.SetRequeueDelay(
				notif.CreationTime,
			)

			// check for lock on object
			locked, ok := r.CheckLock(&kubernetesRuntimeInstance)
			if locked || ok == false {
				r.Requeue(&kubernetesRuntimeInstance, requeueDelay, msg)
				log.V(1).Info("kubernetes runtime instance reconciliation requeued")
				continue
			}

			// set up handler to unlock and requeue on termination signal
			go func() {
				select {
				case <-osSignals:
					log.V(1).Info("received termination signal, performing unlock and requeue of kubernetes runtime instance")
					r.UnlockAndRequeue(&kubernetesRuntimeInstance, requeueDelay, lockReleased, msg)
				case <-lockReleased:
					log.V(1).Info("reached end of reconcile loop for kubernetes runtime instance, closing out signal handler")
				}
			}()

			// put a lock on the reconciliation of the created object
			if ok := r.Lock(&kubernetesRuntimeInstance); !ok {
				r.Requeue(&kubernetesRuntimeInstance, requeueDelay, msg)
				log.V(1).Info("kubernetes runtime instance reconciliation requeued")
				continue
			}

			// retrieve latest version of object
			latestKubernetesRuntimeInstance, err := client.GetKubernetesRuntimeInstanceByID(
				r.APIClient,
				r.APIServer,
				*kubernetesRuntimeInstance.ID,
			)
			// check if error is 404 - if object no longer exists, no need to requeue
			if errors.Is(err, client.ErrObjectNotFound) {
				log.Info(fmt.Sprintf(
					"object with ID %d no longer exists - halting reconciliation",
					*kubernetesRuntimeInstance.ID,
				))
				r.ReleaseLock(&kubernetesRuntimeInstance, lockReleased, msg, true)
				continue
			}
			if err != nil {
				log.Error(err, "failed to get kubernetes runtime instance by ID from API")
				r.UnlockAndRequeue(&kubernetesRuntimeInstance, requeueDelay, lockReleased, msg)
				continue
			}
			kubernetesRuntimeInstance = *latestKubernetesRuntimeInstance

			// determine which operation and act accordingly
			switch notif.Operation {
			case notifications.NotificationOperationCreated:
				if kubernetesRuntimeInstance.DeletionScheduled != nil {
					log.Info("kubernetes runtime instance scheduled for deletion - skipping create")
					break
				}
				customRequeueDelay, err := kubernetesRuntimeInstanceCreated(r, &kubernetesRuntimeInstance, &log)
				if err != nil {
					log.Error(err, "failed to reconcile created kubernetes runtime instance object")
					r.UnlockAndRequeue(
						&kubernetesRuntimeInstance,
						requeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				if customRequeueDelay != 0 {
					log.Info("create requeued for future reconciliation")
					r.UnlockAndRequeue(
						&kubernetesRuntimeInstance,
						customRequeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
			case notifications.NotificationOperationUpdated:
				customRequeueDelay, err := kubernetesRuntimeInstanceUpdated(r, &kubernetesRuntimeInstance, &log)
				if err != nil {
					log.Error(err, "failed to reconcile updated kubernetes runtime instance object")
					r.UnlockAndRequeue(
						&kubernetesRuntimeInstance,
						requeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				if customRequeueDelay != 0 {
					log.Info("update requeued for future reconciliation")
					r.UnlockAndRequeue(
						&kubernetesRuntimeInstance,
						customRequeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
			case notifications.NotificationOperationDeleted:
				customRequeueDelay, err := kubernetesRuntimeInstanceDeleted(r, &kubernetesRuntimeInstance, &log)
				if err != nil {
					log.Error(err, "failed to reconcile deleted kubernetes runtime instance object")
					r.UnlockAndRequeue(
						&kubernetesRuntimeInstance,
						requeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				if customRequeueDelay != 0 {
					log.Info("deletion requeued for future reconciliation")
					r.UnlockAndRequeue(
						&kubernetesRuntimeInstance,
						customRequeueDelay,
						lockReleased,
						msg,
					)
					continue
				}
				deletionTimestamp := util.Ptr(time.Now().UTC())
				deletedKubernetesRuntimeInstance := v0.KubernetesRuntimeInstance{
					Common: v0.Common{ID: kubernetesRuntimeInstance.ID},
					Reconciliation: v0.Reconciliation{
						DeletionAcknowledged: deletionTimestamp,
						DeletionConfirmed:    deletionTimestamp,
						Reconciled:           util.Ptr(true),
					},
				}
				if err != nil {
					log.Error(err, "failed to update kubernetes runtime instance to mark as reconciled")
					r.UnlockAndRequeue(&kubernetesRuntimeInstance, requeueDelay, lockReleased, msg)
					continue
				}
				_, err = client.UpdateKubernetesRuntimeInstance(
					r.APIClient,
					r.APIServer,
					&deletedKubernetesRuntimeInstance,
				)
				if err != nil {
					log.Error(err, "failed to update kubernetes runtime instance to mark as deleted")
					r.UnlockAndRequeue(&kubernetesRuntimeInstance, requeueDelay, lockReleased, msg)
					continue
				}
				_, err = client.DeleteKubernetesRuntimeInstance(
					r.APIClient,
					r.APIServer,
					*kubernetesRuntimeInstance.ID,
				)
				if err != nil {
					log.Error(err, "failed to delete kubernetes runtime instance")
					r.UnlockAndRequeue(&kubernetesRuntimeInstance, requeueDelay, lockReleased, msg)
					continue
				}
			default:
				log.Error(
					errors.New("unrecognized notifcation operation"),
					"notification included an invalid operation",
				)
				r.UnlockAndRequeue(
					&kubernetesRuntimeInstance,
					requeueDelay,
					lockReleased,
					msg,
				)
				continue

			}

			// set the object's Reconciled field to true if not deleted
			if notif.Operation != notifications.NotificationOperationDeleted {
				reconciledKubernetesRuntimeInstance := v0.KubernetesRuntimeInstance{
					Common:         v0.Common{ID: kubernetesRuntimeInstance.ID},
					Reconciliation: v0.Reconciliation{Reconciled: util.Ptr(true)},
				}
				updatedKubernetesRuntimeInstance, err := client.UpdateKubernetesRuntimeInstance(
					r.APIClient,
					r.APIServer,
					&reconciledKubernetesRuntimeInstance,
				)
				if err != nil {
					log.Error(err, "failed to update kubernetes runtime instance to mark as reconciled")
					r.UnlockAndRequeue(&kubernetesRuntimeInstance, requeueDelay, lockReleased, msg)
					continue
				}
				log.V(1).Info(
					"kubernetes runtime instance marked as reconciled in API",
					"kubernetes runtime instanceName", updatedKubernetesRuntimeInstance.Name,
				)
			}

			// release the lock on the reconciliation of the created object
			if ok := r.ReleaseLock(&kubernetesRuntimeInstance, lockReleased, msg, true); !ok {
				log.Error(errors.New("kubernetes runtime instance remains locked - will unlock when TTL expires"), "")
			} else {
				log.V(1).Info("kubernetes runtime instance unlocked")
			}

			log.Info(fmt.Sprintf(
				"kubernetes runtime instance successfully reconciled for %s operation",
				notif.Operation,
			))
		}
	}

	r.Sub.Unsubscribe()
	reconcilerLog.Info("reconciler shutting down")
	r.ShutdownWait.Done()
}
