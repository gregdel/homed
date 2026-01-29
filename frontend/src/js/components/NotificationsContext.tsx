import type React from "react";
import { createContext, useCallback, useContext, useState } from "react";
import type { ReactNode } from "react";

interface Notification {
  id: string;
  message: string;
  type: "success" | "error";
  duration: number;
}

interface NotificationsContextType {
  notifications: Record<string, Notification>;
  addNotificationOk: (message: string) => void;
  addNotificationError: (message: string) => void;
  removeNotification: (id: string) => void;
}

const NotificationsContext = createContext<
  NotificationsContextType | undefined
>(undefined);

export const useNotifications = (): NotificationsContextType => {
  const context = useContext(NotificationsContext);
  if (!context) {
    throw new Error(
      "useNotifications must be used within NotificationsProvider",
    );
  }
  return context;
};

interface NotificationsProviderProps {
  children: ReactNode;
}

export const NotificationsProvider: React.FC<NotificationsProviderProps> = ({
  children,
}) => {
  const [notifications, setNotifications] = useState<
    Record<string, Notification>
  >({});

  const addNotification = useCallback(
    (message: string, type: "success" | "error", duration: number) => {
      const id = Math.random().toString(36).substring(7);
      setNotifications((prevNotifications) => ({
        ...prevNotifications,
        [id]: { id, message, type, duration },
      }));
    },
    [],
  );

  const addNotificationOk = useCallback(
    (message: string) => {
      addNotification(message, "success", 1.5);
    },
    [addNotification],
  );

  const addNotificationError = useCallback(
    (message: string) => {
      addNotification(message, "error", 8);
    },
    [addNotification],
  );

  const removeNotification = useCallback((id: string) => {
    setNotifications((prevNotifications) => {
      const { [id]: _, ...remainingNotifications } = prevNotifications;
      return remainingNotifications;
    });
  }, []);

  const value: NotificationsContextType = {
    notifications,
    addNotificationOk,
    addNotificationError,
    removeNotification,
  };

  return (
    <NotificationsContext.Provider value={value}>
      {children}
    </NotificationsContext.Provider>
  );
};
