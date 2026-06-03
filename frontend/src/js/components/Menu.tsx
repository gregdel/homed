import type React from "react";
import { useEffect, useState } from "react";
import { BottomNav } from "./BottomNav";
import { useComponents } from "./ComponentsContext";
import { useNav } from "./Navigation";
import { getAvailableMenuItems } from "./navigationMenu";
import { Icon } from "./ui/Icon";

const MOBILE_BREAKPOINT = 992;

const Sidebar: React.FC = () => {
  const { currentPath, navigate } = useNav();
  const { components } = useComponents();
  const menuItems = getAvailableMenuItems(components);

  if (menuItems.length === 0) {
    return null;
  }

  return (
    <nav className="sidebar">
      <ul className="nav">
        {menuItems.map((item) => (
          <li
            key={item.key}
            className={`nav-item ${currentPath === item.path ? "active" : ""}`}
          >
            <button onClick={() => navigate(item.path)}>
              <Icon name={item.icon} size={1} />
              <span>{item.label}</span>
            </button>
          </li>
        ))}
      </ul>
    </nav>
  );
};

export const AppMenu: React.FC = () => {
  const [isMobile, setIsMobile] = useState(
    window.innerWidth < MOBILE_BREAKPOINT,
  );

  useEffect(() => {
    const handleResize = () => {
      setIsMobile(window.innerWidth < MOBILE_BREAKPOINT);
    };
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  if (isMobile) {
    return <BottomNav />;
  }

  return <Sidebar />;
};
