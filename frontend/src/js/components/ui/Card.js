import React from "react";
import PropTypes from "prop-types";

export const Card = ({ title, extra, actions, style, children }) => {
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

Card.propTypes = {
  title: PropTypes.node,
  extra: PropTypes.node,
  actions: PropTypes.array,
  style: PropTypes.object,
  children: PropTypes.node,
};

export default Card;
