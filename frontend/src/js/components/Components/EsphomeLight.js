import React from "react";
import PropTypes from "prop-types";
import { useComponents } from "../ComponentsContext";
import { Slider } from "antd";
import Icon from "@mdi/react";
import { mdiLightbulbOnOutline, mdiLightbulbOn } from "@mdi/js";

// Constants
const LIGHT_CONSTANTS = {
  MAX_BRIGHTNESS: 255,
  MIN_OPACITY: 0.2,
  SLIDER: {
    MIN: 0,
    MAX: 100,
    STEP: 10,
  },
  COLOR: {
    ON: "#ffec3d",
    OFF: "#00000040",
  },
  STATES: {
    ON: "ON",
    OFF: "OFF",
  },
  ANIMATION: {
    DURATION: "0.3s",
    EASING: "cubic-bezier(0.4, 0, 0.2, 1)",
  },
};

// Utility Functions
const calculateBrightness = (value) =>
  Math.floor((value / LIGHT_CONSTANTS.MAX_BRIGHTNESS) * 100);

const calculateBrightnessCommand = (value) =>
  Math.floor((value / 100) * LIGHT_CONSTANTS.MAX_BRIGHTNESS);

const calculateOpacity = (isOn, brightness) => {
  if (!isOn) return 1;
  const opacity = brightness / 100;
  return opacity < LIGHT_CONSTANTS.MIN_OPACITY
    ? LIGHT_CONSTANTS.MIN_OPACITY
    : opacity;
};

// Sub-components
const LightIcon = ({ isOn, opacity, onClick }) => (
  <>
    <Icon
      path={isOn ? mdiLightbulbOn : mdiLightbulbOnOutline}
      onClick={onClick}
      style={{
        cursor: "pointer",
        color: isOn ? LIGHT_CONSTANTS.COLOR.ON : LIGHT_CONSTANTS.COLOR.OFF,
        opacity,
        userSelect: "none",
        outline: "none",
        WebkitTapHighlightColor: "transparent",
        transition: `
          color ${LIGHT_CONSTANTS.ANIMATION.DURATION} ${LIGHT_CONSTANTS.ANIMATION.EASING},
          opacity ${LIGHT_CONSTANTS.ANIMATION.DURATION} ${LIGHT_CONSTANTS.ANIMATION.EASING},
          transform ${LIGHT_CONSTANTS.ANIMATION.DURATION} ${LIGHT_CONSTANTS.ANIMATION.EASING}
        `,
        transform: `scale(${isOn ? 1.1 : 1})`,
      }}
    />
    <div
      style={{
        alignSelf: "center",
        opacity: isOn ? 1 : 0.7,
        transform: `translateY(${isOn ? 0 : "-2px"})`,
        transition: `
          opacity ${LIGHT_CONSTANTS.ANIMATION.DURATION} ${LIGHT_CONSTANTS.ANIMATION.EASING},
          transform ${LIGHT_CONSTANTS.ANIMATION.DURATION} ${LIGHT_CONSTANTS.ANIMATION.EASING}
        `,
      }}
    >
      {isOn ? "On" : "Off"}
    </div>
  </>
);

LightIcon.propTypes = {
  isOn: PropTypes.bool.isRequired,
  opacity: PropTypes.number.isRequired,
  onClick: PropTypes.func.isRequired,
};

const BrightnessSlider = ({ brightness, onChange, isVisible }) => (
  <div
    style={{
      display: "grid",
      gridTemplateRows: isVisible ? "1fr" : "0fr",
      opacity: isVisible ? 1 : 0,
      transition: `
        grid-template-rows ${LIGHT_CONSTANTS.ANIMATION.DURATION} ${LIGHT_CONSTANTS.ANIMATION.EASING},
        opacity ${LIGHT_CONSTANTS.ANIMATION.DURATION} ${LIGHT_CONSTANTS.ANIMATION.EASING},
      `,
    }}
  >
    <div style={{ overflow: "hidden", minHeight: 0 }}>
      <Slider
        min={LIGHT_CONSTANTS.SLIDER.MIN}
        max={LIGHT_CONSTANTS.SLIDER.MAX}
        step={LIGHT_CONSTANTS.SLIDER.STEP}
        defaultValue={brightness}
        value={brightness}
        onChange={onChange}
      />
    </div>
  </div>
);

BrightnessSlider.propTypes = {
  brightness: PropTypes.number.isRequired,
  onChange: PropTypes.func.isRequired,
  isVisible: PropTypes.bool.isRequired,
};

const LightContainer = ({ children }) => (
  <div
    style={{
      display: "flex",
      justifyContent: "center",
      flexFlow: "column nowrap",
      alignItems: "stretch",
      height: "30vh",
    }}
  >
    {children}
  </div>
);

LightContainer.propTypes = {
  children: PropTypes.node.isRequired,
};

// Main Component
export const EsphomeLight = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();

  const data = getComponentById(id);
  if (!data?.values?.device) {
    return null;
  }

  const isOn = data.values.device.online && data.values.on;
  const brightness = calculateBrightness(data.values.brightness);
  const opacity = calculateOpacity(isOn, brightness);

  const updateLight = (brightnessValue, isLightOn) => {
    const command = {
      color_mode: "cwww",
      state: isLightOn ? LIGHT_CONSTANTS.STATES.ON : LIGHT_CONSTANTS.STATES.OFF,
      brightness: calculateBrightnessCommand(brightnessValue),
      color: {
        c: LIGHT_CONSTANTS.MAX_BRIGHTNESS,
        w: LIGHT_CONSTANTS.MAX_BRIGHTNESS,
      },
    };
    updateComponent(id, command);
  };

  const handleToggle = () => {
    updateLight(brightness, !isOn);
  };

  const handleBrightnessChange = (value) => {
    updateLight(value, isOn);
  };

  return (
    <LightContainer>
      <LightIcon isOn={isOn} opacity={opacity} onClick={handleToggle} />
      <BrightnessSlider
        brightness={brightness}
        onChange={handleBrightnessChange}
        isVisible={isOn}
      />
    </LightContainer>
  );
};

EsphomeLight.propTypes = {
  id: PropTypes.string.isRequired,
};
