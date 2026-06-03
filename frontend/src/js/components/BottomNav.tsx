import type React from "react";
import { useEffect, useRef, useState } from "react";
import { useComponents } from "./ComponentsContext";
import { useNav } from "./Navigation";
import { getAvailableMenuItems } from "./navigationMenu";
import { Icon } from "./ui/Icon";

const MOBILE_PRIMARY_LIMIT = 4;

export const BottomNav: React.FC = () => {
  const [moreOpen, setMoreOpen] = useState(false);
  const moreButtonRef = useRef<HTMLButtonElement>(null);
  const morePanelRef = useRef<HTMLElement>(null);
  const { currentPath, navigate } = useNav();
  const { components } = useComponents();
  const menuItems = getAvailableMenuItems(components);
  const primaryItems = menuItems.slice(0, MOBILE_PRIMARY_LIMIT);
  const moreItems = menuItems.slice(MOBILE_PRIMARY_LIMIT);

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

  if (menuItems.length === 0) {
    return null;
  }

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
            <span>{item.mobileLabel}</span>
          </button>
        ))}
        {moreItems.length > 0 && (
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
        )}
      </nav>

      {moreOpen && moreItems.length > 0 && (
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
                  <span>{item.mobileLabel}</span>
                </button>
              ))}
            </div>
          </div>
        </nav>
      )}
    </>
  );
};
