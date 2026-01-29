import type React from "react";
import { useComponents } from "./ComponentsContext";
import { useNav } from "./Navigation";

export const Graph: React.FC = () => {
  const { params } = useNav();
  const { getComponentById } = useComponents();

  const data = getComponentById(params.componentId || "");
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
