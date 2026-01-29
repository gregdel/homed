// Navigation types

export interface NavigationParams {
  componentId?: string;
}

export interface NavigationContext {
  currentPath: string;
  params: NavigationParams;
  navigate: (path: string) => void;
}
