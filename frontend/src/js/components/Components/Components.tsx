import React, { useMemo } from "react";
import { useComponents } from "./../ComponentsContext";
import { Component } from "./Component";

interface ComponentsProps {
  typesFilter?: string[];
  noCard?: boolean;
}

interface FilteredComponent {
  id: string;
  roomName?: string;
  friendlyName: string;
  type: string;
  hide: boolean;
}

export const Components: React.FC<ComponentsProps> = ({
  typesFilter = [],
  noCard = false,
}) => {
  const { components, loading, error } = useComponents();

  const filteredAndSortedComponents = useMemo<FilteredComponent[]>(() => {
    return Object.values(components)
      .map((component) => ({
        id: component.values.id,
        roomName: component.values.device?.room,
        friendlyName: component.values.friendly_name,
        type: component.type,
        hide: component.values.hide,
      }))
      .filter((component) => {
        if (!component) {
          return false;
        }

        if (typesFilter.length === 0) {
          return true;
        }

        if (component.hide) {
          return false;
        }

        return typesFilter.includes(component.type);
      })
      .sort((a, b) => {
        // Sort by room name if both have it
        if (a.roomName && b.roomName) {
          const roomCompare = a.roomName.localeCompare(b.roomName);
          if (roomCompare !== 0) {
            return roomCompare;
          }
        }

        // Sort by friendly name if both have it
        if (a.friendlyName && b.friendlyName) {
          return a.friendlyName.localeCompare(b.friendlyName);
        }

        // Fall back to sorting by ID
        return a.id.localeCompare(b.id);
      });
  }, [components, typesFilter]);

  if (loading) {
    return (
      <div className="grid grid-cols-1 grid-cols-sm-2 grid-cols-lg-3">
        <div>Loading components...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="grid grid-cols-1 grid-cols-sm-2 grid-cols-lg-3">
        <div>Error loading components: {error}</div>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 grid-cols-sm-2 grid-cols-lg-3">
      {filteredAndSortedComponents.map(({ id, type }) => (
        <Component key={id} id={id} type={type} noCard={noCard} />
      ))}
    </div>
  );
};

export default Components;
