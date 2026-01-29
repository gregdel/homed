import React, { useState } from "react";
import { useNav } from "./../Navigation";

import { Modal } from "../ui/Modal";
import { apiPost } from "../../utils/api";

interface DefaultValueProps {
  defaultValue: string | number;
  refresh: () => void;
}

export const DefaultValue: React.FC<DefaultValueProps> = ({
  defaultValue,
  refresh,
}) => {
  const { params } = useNav();
  const [value, setValue] = useState<string | number>(defaultValue);
  const [open, setOpen] = useState<boolean>(false);

  const setDefault = async () => {
    try {
      await apiPost(`/components/${params.componentId}/schedule/default`, {
        value: parseInt(value.toString()),
      });
    } catch (error) {
      console.error("Error posting data:", error);
    } finally {
      refresh();
    }
  };

  const handleOk = () => {
    setOpen(false);
    void setDefault();
  };

  return (
    <div className="cursor-pointer">
      <div
        className="flex flex-col"
        onClick={() => {
          setOpen(true);
        }}
      >
        <div style={{ fontSize: "2em" }}>Default:</div>
        <div style={{ fontSize: "3em" }}>{value} °C</div>
      </div>
      <Modal
        title="Set default value"
        open={open}
        onOk={handleOk}
        onCancel={() => setOpen(false)}
      >
        <div className="form-group">
          <label className="form-label">Default value</label>
          <div className="form-control">
            <input
              type="number"
              step=".5"
              className="input"
              defaultValue={defaultValue.toString()}
              onChange={(e) => setValue(e.target.value)}
            />
          </div>
        </div>
      </Modal>
    </div>
  );
};
