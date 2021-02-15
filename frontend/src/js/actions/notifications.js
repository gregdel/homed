export const notificationAdd = (
  message,
  type = "info",
  duration = 10,
  id = Math.random().toString(36).substring(7)
) => ({
  type: "NOTIFICATION_ADD",
  payload: {
    message,
    type,
    duration,
    id,
  },
});

export const notificationRemove = (id) => ({
  type: "NOTIFICATION_REMOVE",
  payload: {
    id,
  },
});
