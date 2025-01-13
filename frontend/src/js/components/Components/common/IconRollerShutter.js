import React, { useMemo } from "react";
import PropTypes from "prop-types";

export const IconRollerShutter = ({ percentOpen = 0 }) => {
  // Ensure percentOpen is between 0 and 100
  const normalizedPercentOpen = Math.max(0, Math.min(100, percentOpen));

  // SVG constants
  const SVG_CONFIG = {
    width: 600,
    height: 600,
    padding: 50,
    frameThickness: 60,
  };

  // Shutter element constants
  const SHUTTER_CONFIG = {
    startY: 175,
    elementHeight: 50,
    spacing: 12,
    width: 300,
    maxElements: 7, // Pre-calculated based on visible height
  };

  // Calculate shutter elements with fixed amount for smooth animation
  const shutterElements = useMemo(() => {
    return Array.from({ length: SHUTTER_CONFIG.maxElements }, (_, index) => {
      const baseY =
        SHUTTER_CONFIG.startY +
        index * (SHUTTER_CONFIG.elementHeight + SHUTTER_CONFIG.spacing);

      // Calculate how "visible" each element should be based on percentOpen
      const elementThreshold = (index + 1) * (100 / SHUTTER_CONFIG.maxElements);
      const isVisible = normalizedPercentOpen < 100 - elementThreshold;

      // Calculate partial visibility for the transitioning element
      const partialVisibility = isVisible ? 1 : 0;

      return {
        id: index,
        y: baseY,
        height: SHUTTER_CONFIG.elementHeight,
        scale: partialVisibility,
      };
    });
  }, [normalizedPercentOpen]);

  return (
    <svg
      className="roller-shutter-icon"
      style={{
        width: "12rem",
        height: "12rem",
      }}
      viewBox={`0 0 ${SVG_CONFIG.width} ${SVG_CONFIG.height}`}
      role="img"
      aria-label={`Roller shutter ${normalizedPercentOpen}% open`}
    >
      {/* Frame */}
      <g className="frame">
        {/* Top frame */}
        <rect
          className="frame-top"
          height={120}
          width={500}
          y={42}
          x={SVG_CONFIG.padding}
          fill="#000000"
        />
        {/* Left frame */}
        <rect
          className="frame-left"
          height={400}
          width={SVG_CONFIG.frameThickness}
          y={150}
          x={75}
          fill="#000000"
        />
        {/* Right frame */}
        <rect
          className="frame-right"
          height={400}
          width={SVG_CONFIG.frameThickness}
          y={150}
          x={465}
          fill="#000000"
        />
      </g>

      {/* Shutter elements with smooth animation */}
      <g className="shutter-elements">
        {shutterElements.map(({ id, height, y, scale }) => (
          <rect
            key={id}
            className="shutter-element"
            height={height}
            width={SHUTTER_CONFIG.width}
            y={y}
            x={150}
            fill="#000000"
            style={{
              transform: `scaleY(${scale})`,
              transformOrigin: `150px ${y}px`,
              transition: "transform 400ms cubic-bezier(0.4, 0, 0.2, 1)",
            }}
          />
        ))}
      </g>
    </svg>
  );
};

IconRollerShutter.propTypes = {
  percentOpen: PropTypes.number,
};

export default IconRollerShutter;
