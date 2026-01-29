import React, { useState, useEffect } from "react";
import PropTypes from "prop-types";
import { useNotifications } from "./NotificationsContext";

import { ToastContainer, Toast } from "./ui/Toast";

export const Notifications = () => {
  const { notifications } = useNotifications();
  if (!notifications || notifications.size === 0) {
    return null;
  }

  return (
    <ToastContainer>
      {Object.values(notifications).map((value) => (
        <Notification
          key={value.id}
          id={value.id}
          type={value.type}
          content={value.message}
          duration={value.duration}
        />
      ))}
    </ToastContainer>
  );
};

const Notification = ({ id, type, content, duration }) => {
  const { removeNotification } = useNotifications();
  const [closed, setClosed] = useState(false);

  const close = () => {
    if (closed) {
      return;
    }

    setClosed(true);
    setTimeout(() => {
      removeNotification(id);
    }, 100);
  };

  useEffect(() => {
    const t = setTimeout(close, duration * 1000);
    return () => clearTimeout(t);
  }, []);

  if (closed) {
    return null;
  }

  return <Toast type={type} content={content} onClick={close} />;
};
Notification.propTypes = {
  id: PropTypes.string.isRequired,
  type: PropTypes.string.isRequired,
  content: PropTypes.string.isRequired,
  duration: PropTypes.number.isRequired,
};
