import React, { useState, useEffect } from "react";
import PropTypes from "prop-types";
import { useNotifications } from "./NotificationsContext";

import { message } from "antd";

export const Notifications = () => {
  const { notifications } = useNotifications();
  if (!notifications || notifications.size === 0) {
    return null;
  }

  return (
    <>
      {Object.values(notifications).map((value) => (
        <Notification
          key={value.id}
          id={value.id}
          type={value.type}
          content={value.message}
          duration={value.duration}
        />
      ))}
    </>
  );
};

const Notification = ({ id, type, content, duration }) => {
  const { removeNotification } = useNotifications();
  const [delay, setDelay] = useState(0);
  const [closed, setClosed] = useState(false);

  const close = () => {
    if (closed) {
      return;
    }

    setClosed(true);
    setDelay(0.1);
    setTimeout(() => {
      removeNotification(id);
    }, 1000);
  };

  useEffect(() => {
    const t = setTimeout(close, duration * 1000);
    return () => clearTimeout(t);
  }, []);

  const config = {
    key: id,
    content,
    duration: delay,
    onClose: close,
    onClick: close,
    style: {
      cursor: "pointer",
    },
  };

  useEffect(() => {
    switch (type) {
      case "error":
        message.error(config);
        break;
      case "success":
        message.success(config);
        break;
      default:
        message.info(config);
        break;
    }
  }, [config]);

  return <></>;
};
Notification.propTypes = {
  id: PropTypes.string.isRequired,
  type: PropTypes.string.isRequired,
  content: PropTypes.string.isRequired,
  duration: PropTypes.number.isRequired,
};
