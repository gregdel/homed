import type React from "react";
import { useEffect, useRef } from "react";

interface ModalProps {
  title: string;
  open: boolean;
  onOk: () => void;
  onCancel: () => void;
  children?: React.ReactNode;
}

export const Modal: React.FC<ModalProps> = ({
  title,
  open,
  onOk,
  onCancel,
  children,
}) => {
  const dialogRef = useRef<HTMLDialogElement>(null);

  useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;

    if (open) {
      dialog.showModal();
    } else {
      dialog.close();
    }
  }, [open]);

  const handleBackdropClick = (e: React.MouseEvent<HTMLDialogElement>) => {
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

export default Modal;
