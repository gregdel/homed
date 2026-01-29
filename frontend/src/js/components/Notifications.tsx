import type React from "react";
import { useEffect, useState } from "react";
import { useNotifications } from "./NotificationsContext";

import { Toast, ToastContainer } from "./ui/Toast";

export const Notifications: React.FC = () => {
  const { notifications } = useNotifications();
  if (!notifications || Object.keys(notifications).length === 0) {
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

interface NotificationProps {
  id: string;
  type: "success" | "error";
  content: string;
  duration: number;
}

const Notification: React.FC<NotificationProps> = ({
  id,
  type,
  content,
  duration,
}) => {
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
