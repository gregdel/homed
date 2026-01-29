import React from "react";
import PropTypes from "prop-types";

export const ToastContainer = ({ children }) => (
  <div className="toast-container">{children}</div>
);

ToastContainer.propTypes = {
  children: PropTypes.node,
};

export const Toast = ({ type = "info", content, onClick }) => {
  const className = `toast toast-${type}`;

  return (
    <div className={className} onClick={onClick}>
      {content}
    </div>
  );
};

Toast.propTypes = {
  type: PropTypes.oneOf(["info", "success", "error"]),
  content: PropTypes.string.isRequired,
  onClick: PropTypes.func,
};

export default Toast;
