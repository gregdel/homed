import type React from "react";
import type { ReactNode } from "react";
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

  if (component.has_graph) {
    extras.push(
      <Link key="graphIcon" to={`/components/${id}/graph`}>
        <Icon name="chartLine" size={0.8} />
      </Link>,
    );
  }

  const roomName = "roomName" in device ? device.roomName : device.room;
  const title =
    friendlyName !== ""
      ? friendlyName
      : `${String(roomName)} - ${component.type} - ${device.name}`;

  return (
    <Card title={title} extra={extras}>
      {children}
    </Card>
  );
};
