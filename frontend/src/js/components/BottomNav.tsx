import type React from "react";
import { useRef, useState } from "react";
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
];

export const BottomNav: React.FC = () => {
  const [moreOpen, setMoreOpen] = useState(false);
  const dialogRef = useRef<HTMLDialogElement>(null);
  const { currentPath, navigate } = useNav();

  const handleNavClick = (path: string) => {
    navigate(path);
    setMoreOpen(false);
    dialogRef.current?.close();
  };

  const handleMoreClick = () => {
    if (moreOpen) {
      dialogRef.current?.close();
      setMoreOpen(false);
    } else {
      dialogRef.current?.showModal();
      setMoreOpen(true);
    }
  };

  const handleDialogClose = () => {
    setMoreOpen(false);
  };

  const handleBackdropClick = (e: React.MouseEvent<HTMLDialogElement>) => {
    // Close if clicking the backdrop (the dialog element itself, not its children)
    if (e.target === dialogRef.current) {
      dialogRef.current?.close();
      setMoreOpen(false);
    }
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
          >
            <Icon name={item.icon} size={1} />
            <span>{item.label}</span>
          </button>
        ))}
        <button
          className={`bottom-nav-item ${isMoreActive || moreOpen ? "active" : ""}`}
          onClick={handleMoreClick}
        >
          <Icon name="dotsHorizontal" size={1} />
          <span>More</span>
        </button>
      </nav>

      <dialog
        ref={dialogRef}
        className="more-sheet"
        onClose={handleDialogClose}
        onClick={handleBackdropClick}
      >
        <div className="more-sheet-content">
          <div className="more-sheet-grid">
            {moreItems.map((item) => (
              <button
                key={item.path}
                className={`more-sheet-item ${currentPath === item.path ? "active" : ""}`}
                onClick={() => handleNavClick(item.path)}
              >
                <Icon name={item.icon} size={1.5} />
                <span>{item.label}</span>
              </button>
            ))}
          </div>
        </div>
      </dialog>
    </>
  );
};
