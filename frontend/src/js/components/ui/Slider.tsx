import type React from "react";
import { useState } from "react";

interface SliderProps {
  min?: number;
  max?: number;
  step?: number;
  value?: number;
  onChange?: (value: number) => void;
  onChangeComplete?: (value: number) => void;
  formatTooltip?: (value: number) => string;
}

export const Slider: React.FC<SliderProps> = ({
  min = 0,
  max = 100,
  step = 1,
  value = min,
  onChange,
  onChangeComplete,
  formatTooltip,
}) => {
  const [active, setActive] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = Number.parseFloat(e.target.value);
    if (onChange) {
      onChange(newValue);
    }
  };

  const handleMouseUp = (
    e: React.MouseEvent<HTMLInputElement> | React.TouchEvent<HTMLInputElement>,
  ) => {
    setActive(false);
    if (onChangeComplete) {
      const target = e.target as HTMLInputElement;
      const newValue = Number.parseFloat(target.value);
      onChangeComplete(newValue);
    }
  };

  const handleStart = () => {
    setActive(true);
  };

  // Calculate percentage for filled track
  const percentage = ((value - min) / (max - min)) * 100;
  const tooltipText = formatTooltip ? formatTooltip(value) : String(value);

  return (
    <div className="slider-container">
      {active && (
        <div className="slider-tooltip" style={{ left: `${percentage}%` }}>
          {tooltipText}
        </div>
      )}
      <input
        type="range"
        className="slider"
        min={min}
        max={max}
        step={step}
        value={value}
        onChange={handleChange}
        onMouseDown={handleStart}
        onTouchStart={handleStart}
        onMouseUp={handleMouseUp}
        onTouchEnd={handleMouseUp}
        style={{ "--slider-percent": `${percentage}%` } as React.CSSProperties}
      />
    </div>
  );
};

export default Slider;
