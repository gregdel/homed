import React from "react";

interface ToastContainerProps {
  children?: React.ReactNode;
}

export const ToastContainer: React.FC<ToastContainerProps> = ({ children }) => (
  <div className="toast-container">{children}</div>
);

interface ToastProps {
  type?: "info" | "success" | "error";
  content: string;
  onClick?: () => void;
}

export const Toast: React.FC<ToastProps> = ({
  type = "info",
  content,
  onClick,
}) => {
  const className = `toast toast-${type}`;

  return (
    <div className={className} onClick={onClick}>
      {content}
    </div>
  );
};

export default Toast;
