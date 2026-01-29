import React from "react";
import { useNav } from "./../Navigation";

import { Icon } from "../ui/Icon";
import { apiDelete } from "../../utils/api";

import { AddOverride } from "./AddOverride";

interface OverrideData {
  id: string;
  start: string;
  stop?: string;
  value: number;
  on: boolean;
}

interface OverridesProps {
  refresh: () => void;
  overrides?: OverrideData[];
}

export const Overrides: React.FC<OverridesProps> = ({
  refresh,
  overrides = [],
}) => {
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

interface OverrideProps extends OverrideData {
  refresh: () => void;
}

const Override: React.FC<OverrideProps> = ({
  id,
  start,
  stop,
  value,
  on,
  refresh,
}) => {
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

  const formatDate = (date: string): string => {
    const d = new Date(date);
    return d.toLocaleString();
  };

  return (
    <div className="list-item">
      <div>
        {on && <Icon name="leaf" size={0.5} />}
        <strong>{value}°C</strong> from <strong>{formatDate(start)}</strong>
        {stop && (
          <>
            {" "}
            to <strong>{formatDate(stop)}</strong>
          </>
        )}
      </div>
      <div
        className="cursor-pointer"
        style={{ alignSelf: "flex-end" }}
        onClick={() => {
          void handleDelete();
        }}
      >
        <Icon name="trashCanOutline" size={1} />
      </div>
    </div>
  );
};
