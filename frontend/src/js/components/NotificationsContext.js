import React, { createContext, useContext, useState } from "react";
import PropTypes from "prop-types";

const NotificationsContext = createContext();

export const useNotifications = () => {
  const context = useContext(NotificationsContext);
  if (!context) {
    throw new Error(
      "useNotifications must be used within NotificationsProvider"
    );
  }
  return context;
};

export const NotificationsProvider = ({ children }) => {
  const [notifications, setNotifications] = useState({});

  const addNotification = (message, type, duration) => {
    const id = Math.random().toString(36).substring(7);
    setNotifications((prevNotifications) => ({
      ...prevNotifications,
      [id]: { id, message, type, duration },
    }));
  };

  const addNotificationOk = (message) => {
    addNotification(message, "success", 1.5);
  };

  const addNotificationError = (message) => {
    addNotification(message, "error", 1);
  };

  const removeNotification = (id) => {
    setNotifications((prevNotifications) => {
      // eslint-disable-next-line no-unused-vars
      const { [id]: _, ...remainingNotifications } = prevNotifications;
      return remainingNotifications;
    });
  };

  const value = {
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

NotificationsProvider.propTypes = {
  children: PropTypes.node,
};
