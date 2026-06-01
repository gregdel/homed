import type React from "react";
import { useComponents } from "../ComponentsContext";
import { Icon } from "../ui/Icon";
import { Slider } from "../ui/Slider";

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
    ON: "var(--color-light-on)",
    OFF: "var(--color-icon-inactive)",
  },
  STATES: {
    ON: "ON",
    OFF: "OFF",
  },
  ANIMATION: {
    DURATION: "0.3s",
    EASING: "cubic-bezier(0.4, 0, 0.2, 1)",
  },
} as const;

// Utility Functions
const calculateBrightness = (value: number): number =>
  Math.floor((value / LIGHT_CONSTANTS.MAX_BRIGHTNESS) * 100);

const calculateBrightnessCommand = (value: number): number =>
  Math.floor((value / 100) * LIGHT_CONSTANTS.MAX_BRIGHTNESS);

const calculateOpacity = (isOn: boolean, brightness: number): number => {
  if (!isOn) return 1;
  const opacity = brightness / 100;
  return opacity < LIGHT_CONSTANTS.MIN_OPACITY
    ? LIGHT_CONSTANTS.MIN_OPACITY
    : opacity;
};

// Sub-components
interface LightIconProps {
  isOn: boolean;
  opacity: number;
  onClick: () => void;
}

const LightIcon: React.FC<LightIconProps> = ({ isOn, opacity, onClick }) => (
  <>
    <Icon
      name={isOn ? "lightbulbOn" : "lightbulbOnOutline"}
      size={6}
      onClick={onClick}
      style={{
        cursor: "pointer",
        alignSelf: "center",
        color: isOn ? LIGHT_CONSTANTS.COLOR.ON : LIGHT_CONSTANTS.COLOR.OFF,
        opacity,
        transition: `color ${LIGHT_CONSTANTS.ANIMATION.DURATION} ease-out`,
      }}
    />
    <div style={{ alignSelf: "center" }}>{isOn ? "On" : "Off"}</div>
  </>
);

interface BrightnessSliderProps {
  brightness: number;
  onChange: (value: number) => void;
  isVisible: boolean;
}

const BrightnessSlider: React.FC<BrightnessSliderProps> = ({
  brightness,
  onChange,
  isVisible,
}) => (
  <div
    style={{
      opacity: isVisible ? 1 : 0,
      visibility: isVisible ? "visible" : "hidden",
      transition: `opacity ${LIGHT_CONSTANTS.ANIMATION.DURATION} ${LIGHT_CONSTANTS.ANIMATION.EASING}`,
      marginTop: "1rem",
    }}
  >
    <Slider
      min={LIGHT_CONSTANTS.SLIDER.MIN}
      max={LIGHT_CONSTANTS.SLIDER.MAX}
      step={LIGHT_CONSTANTS.SLIDER.STEP}
      value={brightness}
      onChange={onChange}
    />
  </div>
);

interface LightContainerProps {
  children: React.ReactNode;
}

const LightContainer: React.FC<LightContainerProps> = ({ children }) => (
  <div
    className="select-none"
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

// Main Component
interface EsphomeLightProps {
  id: string;
}

interface EsphomeLightValues {
  device?: {
    online: boolean;
  };
  on: boolean;
  brightness: number;
  [key: string]: unknown;
}

export const EsphomeLight: React.FC<EsphomeLightProps> = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();

  const data = getComponentById(id);
  if (!data?.values) {
    return null;
  }

  const values = data.values as unknown as EsphomeLightValues;
  if (!values.device) {
    return null;
  }
  const isOn = values.device.online && values.on;
  const brightness = calculateBrightness(values.brightness);
  const opacity = calculateOpacity(isOn, brightness);

  const updateLight = (brightnessValue: number, isLightOn: boolean) => {
    const command = {
      color_mode: "cwww",
      state: isLightOn ? LIGHT_CONSTANTS.STATES.ON : LIGHT_CONSTANTS.STATES.OFF,
      brightness: calculateBrightnessCommand(brightnessValue),
      color: {
        c: LIGHT_CONSTANTS.MAX_BRIGHTNESS,
        w: LIGHT_CONSTANTS.MAX_BRIGHTNESS,
      },
    };
    void updateComponent(id, command);
  };

  const handleToggle = () => {
    updateLight(brightness, !isOn);
  };

  const handleBrightnessChange = (value: number) => {
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
