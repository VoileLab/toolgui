import React from 'react';
import ReactDOM from 'react-dom/client';
// After WSApp: it pulls in the library's stylesheets, and these override them.
import { WSApp } from './wsapp';
import './index.css';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <WSApp />
  </React.StrictMode>
);

