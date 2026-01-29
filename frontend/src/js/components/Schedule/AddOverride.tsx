import React, { useState } from "react";
import { useNav } from "./../Navigation";

import { Modal } from "../ui/Modal";
import { Switch } from "../ui/Switch";
import { Icon } from "../ui/Icon";
import { apiPost } from "../../utils/api";

interface AddOverrideProps {
  refresh: () => void;
}

export const AddOverride: React.FC<AddOverrideProps> = ({ refresh }) => {
  const { params } = useNav();

  const [show, setShow] = useState<boolean>(false);
  const [from, setFrom] = useState<string | undefined>();
  const [to, setTo] = useState<string | undefined>();
  const [target, setTarget] = useState<number>(16);
  const [on, setOn] = useState<boolean>(false);

  const showModal = () => {
    setShow(true);
  };

  const add = async (data: {
    start: string | undefined;
    stop: string | null;
    value: number;
    on: boolean;
  }) => {
    try {
      await apiPost(
        `/components/${params.componentId}/schedule/overrides`,
        data
      );
    } catch (error) {
      console.error("Error posting data:", error);
    } finally {
      refresh();
    }
  };

  const handleOk = () => {
    setShow(false);
    const data = {
      start: from,
      stop: to ? to : null,
      value: Number(target),
      on: on,
    };
    void add(data);
  };

  const handleCancel = () => {
    setShow(false);
  };

  return (
    <div style={{ marginTop: "0.5em" }}>
      <button className="btn btn-primary" onClick={showModal}>
        <Icon name="calendarPlus" size={1} />
      </button>
      <Modal
        title="Add schedule override"
        open={show}
        onOk={handleOk}
        onCancel={handleCancel}
      >
        <div className="form-group">
          <label className="form-label">From</label>
          <div className="form-control">
            <input
              type="datetime-local"
              className="input"
              onChange={(e) => {
                const date = new Date(e.target.value);
                setFrom(date.toISOString());
              }}
            />
          </div>
        </div>

        <div className="form-group">
          <label className="form-label">To</label>
          <div className="form-control">
            <input
              type="datetime-local"
              className="input"
              onChange={(e) => {
                const date = new Date(e.target.value);
                setTo(date.toISOString());
              }}
            />
          </div>
        </div>

        <div className="form-group">
          <label className="form-label">Target</label>
          <div className="form-control">
            <input
              type="number"
              step=".5"
              className="input"
              defaultValue={target}
              onChange={(e) => {
                setTarget(Number(e.target.value));
              }}
            />
          </div>
        </div>

        <div className="form-group">
          <label className="form-label">Opportunistic</label>
          <div className="form-control">
            <Switch checked={on} onChange={() => setOn(!on)} />
          </div>
        </div>
      </Modal>
    </div>
  );
};
