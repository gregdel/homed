import { useEffect } from "react";

export const resumeRefreshEvent = "homed:resume";

export const useResumeRefresh = (callback: () => void): void => {
  useEffect(() => {
    window.addEventListener(resumeRefreshEvent, callback);
    return () => {
      window.removeEventListener(resumeRefreshEvent, callback);
    };
  }, [callback]);
};
