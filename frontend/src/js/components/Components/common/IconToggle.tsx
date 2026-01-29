import type React from "react";
import type { CSSProperties } from "react";
import type { IconName } from "../../ui/Icon";
import { Icon } from "../../ui/Icon";

interface IconToggleProps {
  iconOn: IconName;
  iconOff: IconName;
  toggle: () => void;
  on: boolean;
  rotate?: boolean;
}

export const IconToggle: React.FC<IconToggleProps> = ({
  iconOn,
  iconOff,
  toggle,
  on,
  rotate = false,
}) => {
  const containerStyle: CSSProperties = {
    display: "flex",
    justifyContent: "center",
    flexFlow: "column nowrap",
    height: "30vh",
    WebkitTapHighlightColor: "transparent",
    userSelect: "none",
    outline: "none",
  };

  const iconStyle: CSSProperties = {
    cursor: "pointer",
    color: on ? "var(--color-light-on)" : "var(--color-icon-inactive)",
    transition: "color 0.3s ease-out 0s",
    alignSelf: "center",
  };

  const textStyle: CSSProperties = {
    alignSelf: "center",
  };

  return (
    <div style={containerStyle}>
      <Icon
        name={on ? iconOn : iconOff}
        size={6}
        onClick={toggle}
        className={rotate ? "rotating" : "rotating paused"}
        style={iconStyle}
      />
      <div style={textStyle}>{on ? "On" : "Off"}</div>
    </div>
  );
};
