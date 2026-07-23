package controller

// conditionReady is the standard condition type used across all controller
// status conditions. Keeping it in a shared file avoids duplicate
// declarations and makes it easy for new controllers to reference.
const conditionReady = "Ready"
