export let constants = {};
constants.db = {};
constants.db.UserSequenceKinds = {
	REVIEW: 1,
	APPROVAL: 2,
};
constants.db.ProcessStatuses = {
	DRAFT: 1,
	ACTIVE: 2,
	FINISHED: 3,
	SUSPENDED: 4,
	AWAITING_REVIEW_CONFIRMATION: 5,
	AWAITING_APPROVAL_CONFIRMATION: 6,
};
constants.ui = {};
constants.ui.RoutingModal = {};
constants.ui.RoutingModal.ContextSwitch = {
	MAINTAIN_STATE: 0,
	CHANGE_STATE: 1,
};
