import type React from "react";
import type { EsphomeLightColorMode, EsphomeLightValues } from "../../types";
import { useComponents } from "../ComponentsContext";
import { Icon } from "../ui/Icon";
import { Slider } from "../ui/Slider";

const MAX_CHANNEL = 255;
const MIN_OPACITY = 0.2;
const DEFAULT_COLD_WHITE = 128;
const DEFAULT_WARM_WHITE = 127;
const WARM_WHITE_COLOR = { red: 255, green: 213, blue: 154 };
const COLD_WHITE_COLOR = { red: 221, green: 235, blue: 255 };

const COLOR_PRESETS = [
  { name: "Focus Blue", hex: "#6BA8FF" },
  { name: "Calm Teal", hex: "#48C8C2" },
  { name: "Sage", hex: "#7ECF9A" },
  { name: "Lavender", hex: "#B9A7FF" },
  { name: "Warm Amber", hex: "#FFC46B" },
] as const;

const toPercent = (value: number): number =>
  Math.round((clamp(value, 0, MAX_CHANNEL) / MAX_CHANNEL) * 100);

const toChannel = (value: number): number =>
  Math.round((clamp(value, 0, 100) / 100) * MAX_CHANNEL);

const clamp = (value: number, min: number, max: number): number =>
  Math.min(max, Math.max(min, value));

const channelToHex = (value: number): string =>
  clamp(value, 0, MAX_CHANNEL).toString(16).padStart(2, "0");

const rgbToHex = (red: number, green: number, blue: number): string =>
  `#${channelToHex(red)}${channelToHex(green)}${channelToHex(blue)}`;

const hexToRGB = (
  hex: string,
): { red: number; green: number; blue: number } => {
  const normalized = hex.replace("#", "");
  return {
    red: Number.parseInt(normalized.slice(0, 2), 16),
    green: Number.parseInt(normalized.slice(2, 4), 16),
    blue: Number.parseInt(normalized.slice(4, 6), 16),
  };
};

const DEFAULT_RGB_COLOR = hexToRGB(COLOR_PRESETS[0].hex);

const whiteMixToPercent = (coldWhite: number, warmWhite: number): number => {
  const total = coldWhite + warmWhite;
  if (total === 0) {
    return Math.round((DEFAULT_COLD_WHITE / MAX_CHANNEL) * 100);
  }

  return Math.round((coldWhite / total) * 100);
};

const percentToWhiteMix = (
  value: number,
): { coldWhite: number; warmWhite: number } => {
  const coldWhite = toChannel(value);
  return {
    coldWhite,
    warmWhite: MAX_CHANNEL - coldWhite,
  };
};

const mixChannel = (start: number, end: number, ratio: number): number =>
  Math.round(start + (end - start) * ratio);

const cwwwToHex = (coldWhite: number, warmWhite: number): string => {
  const total = coldWhite + warmWhite;
  const coldRatio =
    total === 0 ? DEFAULT_COLD_WHITE / MAX_CHANNEL : coldWhite / total;

  return rgbToHex(
    mixChannel(WARM_WHITE_COLOR.red, COLD_WHITE_COLOR.red, coldRatio),
    mixChannel(WARM_WHITE_COLOR.green, COLD_WHITE_COLOR.green, coldRatio),
    mixChannel(WARM_WHITE_COLOR.blue, COLD_WHITE_COLOR.blue, coldRatio),
  );
};

const lightDisplayColor = ({
  blue,
  coldWhite,
  green,
  isOn,
  mode,
  red,
  warmWhite,
}: {
  blue: number;
  coldWhite: number;
  green: number;
  isOn: boolean;
  mode: EsphomeLightColorMode;
  red: number;
  warmWhite: number;
}): string => {
  if (!isOn) {
    return "var(--color-icon-inactive)";
  }

  if (mode === "rgb") {
    return rgbToHex(red, green, blue);
  }

  return cwwwToHex(coldWhite, warmWhite);
};

interface EsphomeLightProps {
  id: string;
}

export const EsphomeLight: React.FC<EsphomeLightProps> = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();
  const data = getComponentById(id);
  if (!data?.values) {
    return null;
  }

  const values = data.values as EsphomeLightValues;
  const online = values.device.online;
  const isOn = online && values.on;
  const mode: EsphomeLightColorMode =
    values.color_mode === "rgb" ? "rgb" : "cwww";
  const brightness = toPercent(values.brightness);
  const red = values.red ?? MAX_CHANNEL;
  const green = values.green ?? MAX_CHANNEL;
  const blue = values.blue ?? MAX_CHANNEL;
  const coldWhite = values.cold_white ?? DEFAULT_COLD_WHITE;
  const warmWhite = values.warm_white ?? DEFAULT_WARM_WHITE;
  const whiteMix = whiteMixToPercent(coldWhite, warmWhite);
  const colorHex = rgbToHex(red, green, blue);
  const iconColor = lightDisplayColor({
    blue,
    coldWhite,
    green,
    isOn,
    mode,
    red,
    warmWhite,
  });
  const opacity = isOn ? Math.max(brightness / 100, MIN_OPACITY) : 1;
  const controlsEnabled = online && isOn;
  const statusLabel = !online
    ? "Offline"
    : isOn
      ? `On - ${brightness}%`
      : "Off";

  const sendLightCommand = ({
    nextBrightness = brightness,
    nextMode = mode,
    nextOn = isOn,
    nextRed = red,
    nextGreen = green,
    nextBlue = blue,
    nextColdWhite = coldWhite,
    nextWarmWhite = warmWhite,
  }: {
    nextBrightness?: number;
    nextMode?: EsphomeLightColorMode;
    nextOn?: boolean;
    nextRed?: number;
    nextGreen?: number;
    nextBlue?: number;
    nextColdWhite?: number;
    nextWarmWhite?: number;
  }) => {
    if (!online) {
      return;
    }

    const baseCommand = {
      state: nextOn ? "ON" : "OFF",
      color_mode: nextMode,
      brightness: toChannel(nextBrightness),
    };
    const rgbChannels =
      nextRed === 0 && nextGreen === 0 && nextBlue === 0
        ? DEFAULT_RGB_COLOR
        : { red: nextRed, green: nextGreen, blue: nextBlue };
    const whiteChannels =
      nextColdWhite === 0 && nextWarmWhite === 0
        ? { coldWhite: DEFAULT_COLD_WHITE, warmWhite: DEFAULT_WARM_WHITE }
        : { coldWhite: nextColdWhite, warmWhite: nextWarmWhite };

    const command =
      nextMode === "rgb"
        ? {
            ...baseCommand,
            color: {
              r: clamp(rgbChannels.red, 0, MAX_CHANNEL),
              g: clamp(rgbChannels.green, 0, MAX_CHANNEL),
              b: clamp(rgbChannels.blue, 0, MAX_CHANNEL),
            },
          }
        : {
            ...baseCommand,
            color: {
              c: clamp(whiteChannels.coldWhite, 0, MAX_CHANNEL),
              w: clamp(whiteChannels.warmWhite, 0, MAX_CHANNEL),
            },
          };

    void updateComponent(id, command);
  };

  return (
    <div
      className="select-none"
      style={{
        display: "flex",
        flexDirection: "column",
        gap: "var(--spacing-md)",
        justifyContent: "center",
        minHeight: "10rem",
      }}
    >
      <div
        style={{
          alignItems: "center",
          display: "flex",
          justifyContent: "space-evenly",
        }}
      >
        <button
          type="button"
          onClick={() => sendLightCommand({ nextOn: !isOn })}
          disabled={!online}
          aria-label={isOn ? "Turn light off" : "Turn light on"}
          style={{
            appearance: "none",
            background: "transparent",
            border: 0,
            color: iconColor,
            cursor: online ? "pointer" : "not-allowed",
            lineHeight: 0,
            opacity,
            padding: 0,
            transition: "color 0.3s ease-out, opacity 0.3s ease-out",
          }}
        >
          <Icon name={isOn ? "lightbulbOn" : "lightbulbOnOutline"} size={5} />
        </button>

        <div className="button-group button-group-sm button-group-vertical">
          <button
            type="button"
            className="button-group-item"
            aria-pressed={mode === "cwww"}
            disabled={!controlsEnabled}
            onClick={() => sendLightCommand({ nextMode: "cwww", nextOn: true })}
          >
            White
          </button>
          <button
            type="button"
            className="button-group-item"
            aria-pressed={mode === "rgb"}
            disabled={!controlsEnabled}
            onClick={() => sendLightCommand({ nextMode: "rgb", nextOn: true })}
          >
            Color
          </button>
        </div>
      </div>

      <div
        style={{
          color: online
            ? "var(--color-text-secondary)"
            : "var(--color-text-disabled)",
          fontWeight: 600,
          textAlign: "center",
        }}
      >
        {statusLabel}
      </div>

      <Slider
        min={0}
        max={100}
        step={1}
        value={brightness}
        disabled={!controlsEnabled}
        onChange={(nextBrightness) =>
          sendLightCommand({ nextBrightness, nextOn: true })
        }
        formatTooltip={(value) => `${value}%`}
      />

      {controlsEnabled && mode === "rgb" ? (
        <div>
          <div
            style={{
              display: "grid",
              gap: "var(--spacing-xs)",
              gridTemplateColumns: "repeat(5, 2.25rem)",
              justifyContent: "center",
            }}
          >
            {COLOR_PRESETS.map((preset) => {
              const rgb = hexToRGB(preset.hex);
              const active =
                preset.hex.toLowerCase() === colorHex.toLowerCase();

              return (
                <button
                  key={preset.hex}
                  type="button"
                  aria-label={preset.name}
                  aria-pressed={active}
                  disabled={!controlsEnabled}
                  onClick={() =>
                    sendLightCommand({
                      nextRed: rgb.red,
                      nextGreen: rgb.green,
                      nextBlue: rgb.blue,
                      nextMode: "rgb",
                      nextOn: true,
                    })
                  }
                  style={{
                    appearance: "none",
                    background: preset.hex,
                    border: active
                      ? "3px solid var(--color-text)"
                      : "1px solid var(--color-border)",
                    borderRadius: "var(--radius)",
                    boxShadow: active ? "var(--shadow)" : "none",
                    cursor: controlsEnabled ? "pointer" : "not-allowed",
                    minHeight: "2.25rem",
                    minWidth: "2.25rem",
                    outlineOffset: "2px",
                    padding: 0,
                  }}
                />
              );
            })}
          </div>
        </div>
      ) : null}

      {controlsEnabled && mode === "cwww" ? (
        <div>
          <div
            style={{
              display: "flex",
              justifyContent: "space-between",
              color: "var(--color-text-secondary)",
              fontWeight: 600,
              marginBottom: "var(--spacing-xs)",
            }}
          >
            <span>Warm</span>
            <span>White</span>
            <span>Cool</span>
          </div>
          <Slider
            min={0}
            max={100}
            step={1}
            value={whiteMix}
            disabled={!controlsEnabled}
            onChange={(nextWhiteMix) => {
              const nextMix = percentToWhiteMix(nextWhiteMix);
              sendLightCommand({
                nextColdWhite: nextMix.coldWhite,
                nextWarmWhite: nextMix.warmWhite,
                nextMode: "cwww",
                nextOn: true,
              });
            }}
            formatTooltip={(value) => `${value}% cool`}
          />
        </div>
      ) : null}
    </div>
  );
};
