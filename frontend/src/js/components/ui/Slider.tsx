import type React from "react";

interface SliderProps {
  min?: number;
  max?: number;
  step?: number;
  value?: number;
  onChange?: (value: number) => void;
  onChangeComplete?: (value: number) => void;
}

export const Slider: React.FC<SliderProps> = ({
  min = 0,
  max = 100,
  step = 1,
  value,
  onChange,
  onChangeComplete,
}) => {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = Number.parseFloat(e.target.value);
    if (onChange) {
      onChange(newValue);
    }
  };

  const handleMouseUp = (
    e: React.MouseEvent<HTMLInputElement> | React.TouchEvent<HTMLInputElement>,
  ) => {
    if (onChangeComplete) {
      const target = e.target as HTMLInputElement;
      const newValue = Number.parseFloat(target.value);
      onChangeComplete(newValue);
    }
  };

  return (
    <input
      type="range"
      className="slider"
      min={min}
      max={max}
      step={step}
      value={value}
      onChange={handleChange}
      onMouseUp={handleMouseUp}
      onTouchEnd={handleMouseUp}
    />
  );
};

export default Slider;
