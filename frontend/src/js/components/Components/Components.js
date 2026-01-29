import React, { useMemo } from "react";
import PropTypes from "prop-types";
import { useComponents } from "./../ComponentsContext";
import { Component } from "./Component";

export const Components = ({ typesFilter = [], noCard = false }) => {
  const { components, loading, error } = useComponents();

  const filteredAndSortedComponents = useMemo(() => {
    return Object.values(components)
      .map((component) => ({
        id: component.values.id,
        roomName: component.values.device?.room || undefined,
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
        // Sort by friendly name if both have it
        if (a.friendlyName && b.friendlyName) {
          return a.friendlyName.localeCompare(b.friendlyName);
        }

        // Sort by room name if both have it
        if (a.roomName && b.roomName) {
          return a.roomName.localeCompare(b.roomName);
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

Components.propTypes = {
  typesFilter: PropTypes.array,
  noCard: PropTypes.bool,
};

export default Components;
