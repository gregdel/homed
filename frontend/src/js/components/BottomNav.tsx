import type React from "react";
import { useEffect, useRef, useState } from "react";
import { useNav } from "./Navigation";
import { Icon, type IconName } from "./ui/Icon";

interface NavItem {
  icon: IconName;
  label: string;
  path: string;
}

const primaryItems: NavItem[] = [
  { icon: "homeThermometer", label: "Temperature", path: "/temperature" },
  { icon: "lightbulb", label: "Lights", path: "/lights" },
  { icon: "windowShutter", label: "Shutters", path: "/shutters" },
  { icon: "power", label: "Switches", path: "/switches" },
];

const moreItems: NavItem[] = [
  { icon: "fan", label: "Fans", path: "/fans" },
  { icon: "ruler", label: "Sensors", path: "/sensors" },
  { icon: "lightningBolt", label: "Power", path: "/power" },
  { icon: "radiator", label: "TRVs", path: "/trv" },
  { icon: "thermometer", label: "Climate", path: "/climate_sensors" },
  { icon: "automation", label: "Automations", path: "/automations" },
];

export const BottomNav: React.FC = () => {
  const [moreOpen, setMoreOpen] = useState(false);
  const moreButtonRef = useRef<HTMLButtonElement>(null);
  const morePanelRef = useRef<HTMLElement>(null);
  const { currentPath, navigate } = useNav();

  useEffect(() => {
    if (!moreOpen) return;

    const handlePointerDown = (e: PointerEvent) => {
      const target = e.target as Node;
      if (
        morePanelRef.current?.contains(target) ||
        moreButtonRef.current?.contains(target)
      ) {
        return;
      }
      setMoreOpen(false);
    };

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        setMoreOpen(false);
        moreButtonRef.current?.focus();
      }
    };

    document.addEventListener("pointerdown", handlePointerDown);
    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.removeEventListener("pointerdown", handlePointerDown);
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [moreOpen]);

  useEffect(() => {
    if (currentPath !== "") {
      setMoreOpen(false);
    }
  }, [currentPath]);

  const handleNavClick = (path: string) => {
    navigate(path);
    setMoreOpen(false);
  };

  const handleMoreClick = () => {
    setMoreOpen((open) => !open);
  };

  const isMoreActive = moreItems.some((item) => currentPath === item.path);

  return (
    <>
      <nav className="bottom-nav">
        {primaryItems.map((item) => (
          <button
            key={item.path}
            className={`bottom-nav-item ${currentPath === item.path ? "active" : ""}`}
            onClick={() => handleNavClick(item.path)}
            type="button"
          >
            <Icon name={item.icon} size={1} />
            <span>{item.label}</span>
          </button>
        ))}
        <button
          ref={moreButtonRef}
          aria-controls="mobile-more-nav"
          aria-expanded={moreOpen}
          className={`bottom-nav-item ${isMoreActive || moreOpen ? "active" : ""}`}
          onClick={handleMoreClick}
          type="button"
        >
          <Icon name="dotsHorizontal" size={1} />
          <span>More</span>
        </button>
      </nav>

      {moreOpen && (
        <nav
          ref={morePanelRef}
          aria-label="More navigation"
          className="more-sheet"
          id="mobile-more-nav"
        >
          <div className="more-sheet-content">
            <div className="more-sheet-grid">
              {moreItems.map((item) => (
                <button
                  key={item.path}
                  aria-current={currentPath === item.path ? "page" : undefined}
                  className={`more-sheet-item ${currentPath === item.path ? "active" : ""}`}
                  onClick={() => handleNavClick(item.path)}
                  type="button"
                >
                  <Icon name={item.icon} size={1.5} />
                  <span>{item.label}</span>
                </button>
              ))}
            </div>
          </div>
        </nav>
      )}
    </>
  );
};
