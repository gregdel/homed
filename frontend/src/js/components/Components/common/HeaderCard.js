import React from "react";
import PropTypes from "prop-types";
import { relativeTime } from "../../../utils/relativeTime";

import { Icon } from "../../ui/Icon";
import { Card } from "../../ui/Card";

import { Link } from "../../Navigation";
import { useComponents } from "../../ComponentsContext";

export const HeaderCard = ({ id, children }) => {
  const { getComponentById } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const {
    updated_at: updatedAt,
    graph_url: graphURL,
    friendly_name: friendlyName,
    device: device,
    type,
  } = component.values;

  var extras = [];
  extras.push(
    <div key="status">
      {device.online === true ? relativeTime(updatedAt) : "offline"}
    </div>
  );

  if (graphURL && graphURL != "") {
    extras.push(
      <div key="graphIcon" style={{ marginRight: "-1em", marginLeft: "0.3em" }}>
        <Link to={`/components/${id}/graph`}>
          <Icon name="chartLine" style={{ color: "#000000d9" }} size={0.8} />
        </Link>
      </div>
    );
  }

  const title =
    friendlyName !== ""
      ? friendlyName
      : `${device.roomName} - ${type} - ${device.name}`;
  return (
    <Card title={title} extra={<div style={{ display: "flex" }}>{extras}</div>}>
      {children}
    </Card>
  );
};
HeaderCard.propTypes = {
  id: PropTypes.string.isRequired,
  children: PropTypes.object.isRequired,
};
