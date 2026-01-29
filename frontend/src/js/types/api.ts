// API response types

export interface APIResponse<T> {
  status: "success" | "error";
  data: T;
}

export interface APIError {
  status: "error";
  message: string;
}

// Component command types (for updates)
export type SwitchCommand = "ON" | "OFF";

export interface TemperatureCommand {
  mode?: string;
  manual_target?: number;
  manual_until?: string;
  manual_duration?: string; // e.g., "1h"
  on?: boolean;
}

export type RollerShutterCommand = "open" | "close" | "stop";

export type ComponentCommand = string | Record<string, unknown>;
