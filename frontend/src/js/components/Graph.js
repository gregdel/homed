import React from "react";
import { useSelector } from "react-redux";
import { useParams } from "react-router-dom";

export const Graph = () => {
  const { componentId: id } = useParams();

  const data = useSelector((state) => state.components.components.get(id));
  if (!data || !data.values || !data.values.graph_url) {
    return null;
  }

  const url = data.values.graph_url;

  return (
    <div style={{ margin: "-1em", height: "100vh" }}>
      <iframe style={{ width: "100%", height: "100%" }} src={url} />
    </div>
  );
};
