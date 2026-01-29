import React from "react";
import { useNav } from "./../Navigation";

import PropTypes from "prop-types";

import { Icon } from "../ui/Icon";
import { apiDelete } from "../../utils/api";

import { AddOverride } from "./AddOverride";

export const Overrides = ({ refresh, overrides = [] }) => {
  return (
    <>
      <div className="flex justify-between items-baseline">
        <h3>Overrides</h3>
        <AddOverride refresh={refresh} />
      </div>
      {overrides.length !== 0 && (
        <div className="list">
          {overrides.map((v) => (
            <Override key={v.id} refresh={refresh} {...v} />
          ))}
        </div>
      )}
      {overrides.length === 0 && <div>No override defined</div>}
    </>
  );
};
Overrides.propTypes = {
  overrides: PropTypes.array,
  refresh: PropTypes.func.isRequired,
};

const Override = ({ id, start, stop, value, on, refresh }) => {
  const { params } = useNav();

  const handleDelete = async () => {
    try {
      await apiDelete(
        `/components/${params.componentId}/schedule/overrides/${id}`
      );
    } catch (error) {
      console.error("Error deleting override:", error);
    } finally {
      refresh();
    }
  };

  const formatDate = (date) => {
    const d = new Date(date);
    return d.toLocaleString();
  };

  return (
    <div className="list-item">
      <div>
        {on && <Icon name="leaf" size={0.5} />}
        <strong>{value}°C</strong> from <strong>{formatDate(start)}</strong> to{" "}
        <strong>{formatDate(stop)}</strong>
      </div>
      <div
        className="cursor-pointer"
        style={{ alignSelf: "flex-end" }}
        onClick={handleDelete}
      >
        <Icon name="trashCanOutline" size={1} />
      </div>
    </div>
  );
};
Override.propTypes = {
  id: PropTypes.string.isRequired,
  start: PropTypes.string.isRequired,
  stop: PropTypes.string.isRequired,
  value: PropTypes.number.isRequired,
  on: PropTypes.bool.isRequired,
  refresh: PropTypes.func.isRequired,
};
