import type React from "react";
import type { CSSProperties, ReactNode } from "react";
import { relativeTime } from "../../../utils/relativeTime";
import { useComponents } from "../../ComponentsContext";
import { Link } from "../../Navigation";
import { Card } from "../../ui/Card";
import { Icon } from "../../ui/Icon";

interface HeaderCardProps {
  id: string;
  children: ReactNode;
}

export const HeaderCard: React.FC<HeaderCardProps> = ({ id, children }) => {
  const { getComponentById } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const {
    updated_at: updatedAt,
    graph_url: graphURL,
    friendly_name: friendlyName,
    device,
  } = component.values;

  const extras: ReactNode[] = [];
  extras.push(
    <div key="status">
      {"online" in device && device.online === true
        ? relativeTime(updatedAt)
        : "offline"}
    </div>,
  );

  if (graphURL && graphURL !== "") {
    const graphIconStyle: CSSProperties = {
      marginRight: "-1em",
      marginLeft: "0.3em",
    };
    const iconStyle: CSSProperties = { color: "#000000d9" };

    extras.push(
      <div key="graphIcon" style={graphIconStyle}>
        <Link to={`/components/${id}/graph`}>
          <Icon name="chartLine" style={iconStyle} size={0.8} />
        </Link>
      </div>,
    );
  }

  const roomName = "roomName" in device ? device.roomName : device.room;
  const title =
    friendlyName !== ""
      ? friendlyName
      : `${String(roomName)} - ${component.type} - ${device.name}`;

  const extraContainerStyle: CSSProperties = { display: "flex" };

  return (
    <Card title={title} extra={<div style={extraContainerStyle}>{extras}</div>}>
      {children}
    </Card>
  );
};
