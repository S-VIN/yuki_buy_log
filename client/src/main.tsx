import React from "react";
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from "react-router-dom";
import {ConfigProvider} from "antd";

import './index.css';

import App from './App.tsx';

const root = ReactDOM.createRoot(document.getElementById('root')!);

try {
  root.render(
    <React.StrictMode>
      <BrowserRouter>

      <ConfigProvider
        theme={{
            components: {
              Card: {
                 bodyPadding: 8,
              },
              Modal: {
                  bodyPadding: 0,
              }
            },
      }}
      >

            <App/>
        </ConfigProvider>
      </BrowserRouter>
    </React.StrictMode>
  );
} catch (e) {
  console.log('unsupport env')
}


