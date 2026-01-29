import React, { useState } from "react";
import PropTypes from "prop-types";
import { useNav } from "../Navigation";
import { useNotifications } from "../NotificationsContext";

import { Icon } from "../ui/Icon";
import { Modal } from "../ui/Modal";
import { Switch } from "../ui/Switch";
import { apiPost, apiPut } from "../../utils/api";

import { daysMap } from "./DailySchedule";

const DaysCheckboxes = ({ checked, setChecked }) => {
  const handleCheck = (i) => {
    let values = [...checked];
    values[i] = !values[i];
    setChecked(values);
  };

  return (
    <div className="flex flex-col">
      <hr className="divider" />
      {[1, 2, 3, 4, 5, 6, 0].map((day) => (
        <div
          key={`days-checkbox-${day}`}
          className="cursor-pointer flex justify-between"
          style={{
            marginLeft: "3em",
            marginRight: "3em",
          }}
          onClick={() => handleCheck(day)}
        >
          <h3 className="text-light">{daysMap[day]}</h3>
          <div className="flex items-center">
            <div>
              <Icon
                name={checked[day] ? "checkboxBlank" : "checkboxBlankOutline"}
                size={2}
              />
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};
DaysCheckboxes.propTypes = {
  checked: PropTypes.array.isRequired,
  setChecked: PropTypes.func.isRequired,
};

export const TimeSlotModal = ({
  start: defaultStart,
  stop: defaultStop,
  target: defaultTarget = 16,
  on: defaultOn = false,
  edit = false,
  refresh,
  open,
  setOpen,
  id,
  day,
}) => {
  const { params } = useNav();
  const { addNotificationError } = useNotifications();

  const [start, setStart] = useState(defaultStart);
  const [stop, setStop] = useState(defaultStop);
  const [target, setTarget] = useState(defaultTarget);
  const [on, setOn] = useState(defaultOn);
  const [checked, setChecked] = useState([
    false,
    false,
    false,
    false,
    false,
    false,
    false,
  ]);

  const add = async (day, data) => {
    try {
      await apiPost(
        `/components/${params.componentId}/schedule/daily/${day}`,
        data
      );
    } catch (error) {
      addNotificationError(error.message);
    } finally {
      refresh();
    }
  };

  const update = async (data) => {
    try {
      await apiPut(
        `/components/${params.componentId}/schedule/daily/${day}/${id}`,
        data
      );
    } catch (error) {
      addNotificationError(error.message);
    } finally {
      refresh();
    }
  };

  const handleOk = () => {
    setOpen(false);
    const data = {
      start: start,
      stop: stop ? stop : null,
      value: new Number(target),
      on: on,
    };

    if (edit) {
      update(data);
    } else {
      if (day === undefined) {
        checked.map((v, day) => {
          if (v === false) {
            return;
          }
          add(day, data);
        });
      } else {
        add(day, data);
      }
    }
  };

  const handleCancel = () => {
    setOpen(false);
  };

  const title = edit ? "Edit timeslot" : "Add in schedule";

  return (
    <Modal title={title} open={open} onOk={handleOk} onCancel={handleCancel}>
      <div className="form-group">
        <label className="form-label">From</label>
        <div className="form-control">
          <input
            type="time"
            className="input"
            defaultValue={defaultStart}
            onChange={(e) => setStart(e.target.value)}
          />
        </div>
      </div>

      <div className="form-group">
        <label className="form-label">To</label>
        <div className="form-control">
          <input
            type="time"
            className="input"
            defaultValue={defaultStop}
            onChange={(e) => setStop(e.target.value)}
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
            defaultValue={defaultTarget}
            onChange={(e) => setTarget(e.target.value)}
          />
        </div>
      </div>

      <div className="form-group">
        <label className="form-label">Opportunistic</label>
        <div className="form-control">
          <Switch checked={on} onChange={() => setOn(!on)} />
        </div>
      </div>

      {day === undefined && (
        <DaysCheckboxes checked={checked} setChecked={setChecked} />
      )}
    </Modal>
  );
};
TimeSlotModal.propTypes = {
  open: PropTypes.bool.isRequired,
  setOpen: PropTypes.func.isRequired,
  day: PropTypes.number,
  id: PropTypes.string,
  start: PropTypes.any,
  stop: PropTypes.any,
  target: PropTypes.number,
  on: PropTypes.bool,
  edit: PropTypes.bool,
  refresh: PropTypes.func.isRequired,
};
