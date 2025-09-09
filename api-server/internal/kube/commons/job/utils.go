package job

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

func hasJobCondition(conditions []batchv1.JobCondition, conditionType batchv1.JobConditionType) bool {
	for _, condition := range conditions {
		if condition.Type == conditionType {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

type Status int32

const (
	StatusUnknown Status = iota
	StatusRunning
	StatusComplete
	StatusFailed
	StatusTerminating
	StatusSuspended
	StatusFailureTarget
	StatusSuccessCriteriaMet
)

func GetJobStatus(job batchv1.Job) Status {
	status := StatusUnknown
	if hasJobCondition(job.Status.Conditions, batchv1.JobComplete) {
		status = StatusComplete
	} else if hasJobCondition(job.Status.Conditions, batchv1.JobFailed) {
		status = StatusFailed
	} else if job.ObjectMeta.DeletionTimestamp != nil {
		status = StatusTerminating
	} else if hasJobCondition(job.Status.Conditions, batchv1.JobSuspended) {
		status = StatusSuspended
	} else if hasJobCondition(job.Status.Conditions, batchv1.JobFailureTarget) {
		status = StatusFailureTarget
	} else if hasJobCondition(job.Status.Conditions, batchv1.JobSuccessCriteriaMet) {
		status = StatusSuccessCriteriaMet
	} else {
		status = StatusRunning
	}
	return status
}
