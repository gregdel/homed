import React from "react";

interface CardProps {
  title?: React.ReactNode;
  extra?: React.ReactNode;
  actions?: React.ReactNode[];
  style?: React.CSSProperties;
  children?: React.ReactNode;
}

export const Card: React.FC<CardProps> = ({
  title,
  extra,
  actions,
  style,
  children,
}) => {
  return (
    <div className="card" style={style}>
      {(title || extra) && (
        <div className="card-header">
          {title && <div className="card-header-title">{title}</div>}
          {extra && <div className="card-header-extra">{extra}</div>}
        </div>
      )}
      <div className="card-body">{children}</div>
      {actions && actions.length > 0 && (
        <div className="card-actions">
          {actions.map((action, i) => (
            <div key={i}>{action}</div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Card;
