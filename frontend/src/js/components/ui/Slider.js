import React from "react";
import PropTypes from "prop-types";

export const Slider = ({
  min = 0,
  max = 100,
  step = 1,
  value,
  onChange,
  onChangeComplete,
}) => {
  const handleChange = (e) => {
    const newValue = parseFloat(e.target.value);
    if (onChange) {
      onChange(newValue);
    }
  };

  const handleMouseUp = (e) => {
    if (onChangeComplete) {
      const newValue = parseFloat(e.target.value);
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

Slider.propTypes = {
  min: PropTypes.number,
  max: PropTypes.number,
  step: PropTypes.number,
  value: PropTypes.number,
  onChange: PropTypes.func,
  onChangeComplete: PropTypes.func,
};

export default Slider;
