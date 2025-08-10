import type { ComponentType, JSX } from 'react';

import AddCheckScreen from './Screens/AddCheckScreen.tsx';
import ViewChecksScreen from "./Screens/ViewChecksScreen.tsx";

interface Route {
  path: string;
  Component: ComponentType;
  title?: string;
  icon?: JSX.Element;
}

export const routes: Route[] = [
  { path: '/', Component: AddCheckScreen, title: 'AddCheckScreen' },
  { path: '/add-check-screen', Component: AddCheckScreen, title: 'AddCheckScreen' },
  { path: '/view-checks-screen', Component: ViewChecksScreen, title: 'ViewChecksScreen' },
];