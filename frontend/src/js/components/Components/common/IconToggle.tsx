import type React from "react";
import type { ReactNode } from "react";
import type { CSSProperties } from "react";
import type { IconName } from "../../ui/Icon";
import { Icon } from "../../ui/Icon";

interface IconToggleProps {
  iconOn: IconName;
  iconOff: IconName;
  toggle: () => void;
  on: boolean;
  rotate?: boolean;
  disabled?: boolean;
  label?: ReactNode;
  activeLabel?: ReactNode;
  colorOn?: string;
  colorOff?: string;
  ariaLabel?: string;
}

export const IconToggle: React.FC<IconToggleProps> = ({
  iconOn,
  iconOff,
  toggle,
  on,
  rotate = false,
  disabled = false,
  label = "Off",
  activeLabel = "On",
  colorOn = "var(--color-light-on)",
  colorOff = "var(--color-icon-inactive)",
  ariaLabel,
}) => {
  const containerStyle: CSSProperties = {
    display: "flex",
    justifyContent: "center",
    flexFlow: "column nowrap",
    height: "30vh",
  };

  const iconStyle: CSSProperties = {
    color: disabled ? "var(--color-icon-inactive)" : on ? colorOn : colorOff,
    transition: "color 0.3s ease-out 0s",
  };

  const buttonStyle: CSSProperties = {
    alignSelf: "center",
    padding: 0,
    border: 0,
    background: "transparent",
    color: "inherit",
    cursor: disabled ? "not-allowed" : "pointer",
  };

  const textStyle: CSSProperties = {
    alignSelf: "center",
  };
  const fallbackAriaLabel =
    typeof (on ? activeLabel : label) === "string"
      ? String(on ? activeLabel : label)
      : undefined;
  const buttonAriaLabel = ariaLabel ?? fallbackAriaLabel;

  return (
    <div className="select-none" style={containerStyle}>
      <button
        {...(buttonAriaLabel !== undefined
          ? { "aria-label": buttonAriaLabel }
          : {})}
        disabled={disabled}
        onClick={toggle}
        style={buttonStyle}
        type="button"
      >
        <Icon
          name={on ? iconOn : iconOff}
          size={6}
          className={rotate ? "rotating" : "rotating paused"}
          style={iconStyle}
        />
      </button>
      <div style={textStyle}>{on ? activeLabel : label}</div>
    </div>
  );
};
