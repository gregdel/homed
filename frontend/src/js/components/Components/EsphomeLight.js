import React from "react";
import { useComponents } from "../ComponentsContext";

import PropTypes from "prop-types";

import { mdiLightbulbOnOutline, mdiLightbulbOn } from "@mdi/js";
import { Slider } from "antd";

import Icon from "@mdi/react";

export const EsphomeLight = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();

  const data = getComponentById(id);
  if (!data || !data.values || !data.values.device) {
    return null;
  }

  const on = data.values.device.online && data.values.on;
  const brightness = Math.floor((data.values.brightness / 255) * 100);

  const opacity = on ? (brightness / 100).toFixed(2) : 1;

  const update = (b, o) => {
    const cmd = {
      color_mode: "cwww",
      state: o ? "ON" : "OFF",
      brightness: Math.floor((b / 100) * 255),
      color: {
        c: 255,
        w: 255,
      },
    };
    updateComponent(id, cmd);
  };

  const toggle = () => {
    update(brightness, !on);
  };

  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        flexFlow: "column nowrap",
        alignItems: "stretch",
        height: "30vh",
      }}
    >
      <Icon
        path={on ? mdiLightbulbOn : mdiLightbulbOnOutline}
        onClick={toggle}
        style={{
          cursor: "pointer",
          color: on ? "#ffec3d" : "#00000040",
          opacity: opacity < 0.2 ? 0.2 : opacity,
          transition: "color 0.3s ease-out 0s",
        }}
      />
      <div
        style={{
          alignSelf: "center",
        }}
      >
        {on ? "On" : "Off"}
      </div>
      {on && (
        <Slider
          min={0}
          max={100}
          step={10}
          defaultValue={brightness}
          value={brightness}
          onChange={(v) => {
            update(v, on);
          }}
        />
      )}
    </div>
  );
};

EsphomeLight.propTypes = {
  id: PropTypes.string.isRequired,
};
