import type React from "react";
import { useState } from "react";
import type { ScriptValues } from "../../types";
import { useComponents } from "../ComponentsContext";
import { IconToggle } from "./common/IconToggle";

interface ScriptProps {
  id: string;
}

const formatDuration = (durationMs: number): string => {
  if (durationMs <= 0) {
    return "";
  }

  if (durationMs < 1000) {
    return `${durationMs}ms`;
  }

  return `${(durationMs / 1000).toFixed(1)}s`;
};

const statusClasses: Record<ScriptValues["status"], string> = {
  idle: "text-secondary",
  running: "text-primary",
  success: "text-success",
  error: "text-error",
};

export const Script: React.FC<ScriptProps> = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();
  const [pending, setPending] = useState(false);

  const component = getComponentById(id);
  if (component === undefined) return null;

  const values = component.values as ScriptValues;
  const running = values.status === "running";
  const disabled = pending || running;
  const duration = formatDuration(values.last_duration_ms);
  const status = (
    <div className="flex flex-wrap items-center justify-center gap-xs text-secondary">
      <span
        className={`font-semibold text-capitalize ${statusClasses[values.status]}`}
      >
        {running ? "Running" : values.status}
      </span>
      {duration !== "" && <span>{duration}</span>}
      {values.last_exit_code !== null && (
        <span>exit {values.last_exit_code}</span>
      )}
      {values.last_error !== "" && (
        <span className="break-anywhere">{values.last_error}</span>
      )}
    </div>
  );

  const run = async () => {
    if (disabled) {
      return;
    }

    setPending(true);
    try {
      await updateComponent(id, {});
    } finally {
      setPending(false);
    }
  };

  return (
    <IconToggle
      activeLabel={status}
      ariaLabel={`Run ${values.friendly_name || id}`}
      colorOff="var(--color-primary)"
      colorOn="var(--color-primary)"
      disabled={disabled}
      iconOff="playCircleOutline"
      iconOn="autorenew"
      label={status}
      on={running}
      rotate={running}
      toggle={() => {
        void run();
      }}
    />
  );
};
