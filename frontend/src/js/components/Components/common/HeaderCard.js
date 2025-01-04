import React from "react";
import PropTypes from "prop-types";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
dayjs.extend(relativeTime);

import Icon from "@mdi/react";
import { mdiChartLine } from "@mdi/js";

import { Card } from "antd";

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
      {device.online === true ? dayjs(updatedAt).fromNow() : "offline"}
    </div>
  );

  if (graphURL && graphURL != "") {
    extras.push(
      <div key="graphIcon" style={{ marginRight: "-1em", marginLeft: "0.3em" }}>
        <Link to={`/components/${id}/graph`}>
          <Icon path={mdiChartLine} style={{ color: "#000000d9" }} size={0.8} />
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
