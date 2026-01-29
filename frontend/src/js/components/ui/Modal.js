import React, { useRef, useEffect } from "react";
import PropTypes from "prop-types";

export const Modal = ({ title, open, onOk, onCancel, children }) => {
  const dialogRef = useRef(null);

  useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;

    if (open) {
      dialog.showModal();
    } else {
      dialog.close();
    }
  }, [open]);

  const handleBackdropClick = (e) => {
    if (e.target === dialogRef.current) {
      onCancel();
    }
  };

  return (
    <dialog ref={dialogRef} onClick={handleBackdropClick}>
      <div className="modal-header">
        <span>{title}</span>
        <button className="modal-close" onClick={onCancel} type="button">
          ✕
        </button>
      </div>
      <div className="modal-body">{children}</div>
      <div className="modal-footer">
        <button className="btn" onClick={onCancel} type="button">
          Cancel
        </button>
        <button className="btn btn-primary" onClick={onOk} type="button">
          OK
        </button>
      </div>
    </dialog>
  );
};

Modal.propTypes = {
  title: PropTypes.string.isRequired,
  open: PropTypes.bool.isRequired,
  onOk: PropTypes.func.isRequired,
  onCancel: PropTypes.func.isRequired,
  children: PropTypes.node,
};

export default Modal;
