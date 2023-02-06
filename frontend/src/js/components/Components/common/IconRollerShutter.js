import React from "react";

import PropTypes from "prop-types";

export const IconRollerShutter = ({ percentOpen }) => {
  const percentClose = 100 - percentOpen;
  const totalHeight = 375;

  const elementStartY = 175;
  const elementHeight = 50;
  const elementSpacing = 13;

  const maxElements = totalHeight / (elementHeight + elementSpacing);
  const factor = 100 / maxElements;

  const elementsCount = percentClose / factor;
  const elementsRemain = percentClose % factor;

  var elements = [];
  for (var i = 0; i < elementsCount; i++) {
    elements.push([i, elementHeight]);
  }

  if (elements.length !== 0 && elementsRemain !== 0) {
    elements[elements.length - 1][1] =
      (elementHeight / factor) * elementsRemain;
  }

  return (
    <svg
      style={{
        width: "12rem",
        height: "12rem",
      }}
      viewBox="0 0 600 600"
      role="presentation"
    >
      <g>
        <g id="svg_background">
          <rect
            id="svg_top"
            height="120"
            width="500"
            y="42"
            x="50"
            stroke="#000"
            fill="#000000"
          />
          <rect
            id="svg_left"
            height="400"
            width="60"
            y="150"
            x="75"
            stroke="#000"
            fill="#000000"
          />
          <rect
            id="svg_right"
            height="400"
            width="60"
            y="150"
            x="465"
            stroke="#000"
            fill="#000000"
          />
        </g>

        {elements.map(([i, height]) => {
          return (
            <rect
              key={i}
              height={height}
              width="300"
              style={{
                transition: "height 200ms linear",
              }}
              y={elementStartY + i * elementHeight + i * elementSpacing}
              x="150"
              stroke="#000"
              fill="#000000"
            />
          );
        })}
      </g>
    </svg>
  );
};

IconRollerShutter.propTypes = {
  percentOpen: PropTypes.number.isRequired,
};
