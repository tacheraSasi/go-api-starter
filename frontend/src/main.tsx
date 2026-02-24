/* eslint-disable import/extensions */
/* eslint-disable @typescript-eslint/no-unsafe-call */
import React from 'react';

import ReactDOM from 'react-dom/client';

import App from './App.tsx';
import './styles/app.scss';

const rootElement = document.getElementById('root')!;
ReactDOM.createRoot(rootElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
