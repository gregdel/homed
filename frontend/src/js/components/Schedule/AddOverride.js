import React, { useState } from "react";
import { useNav } from "./../Navigation";
import PropTypes from "prop-types";

import { Modal } from "../ui/Modal";
import { Switch } from "../ui/Switch";
import { Icon } from "../ui/Icon";
import { apiPost } from "../../utils/api";

export const AddOverride = ({ refresh }) => {
  const { params } = useNav();

  const [show, setShow] = useState(false);
  const [from, setFrom] = useState();
  const [to, setTo] = useState();
  const [target, setTarget] = useState(16);
  const [on, setOn] = useState(false);

  const showModal = () => {
    setShow(true);
  };

  const add = async (data) => {
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
      value: new Number(target),
      on: on,
    };
    add(data);
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
                setTarget(e.target.value);
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
AddOverride.propTypes = {
  refresh: PropTypes.func.isRequired,
};
