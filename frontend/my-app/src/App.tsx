import { useState } from 'react'
import HomePage from './pages/home-page'
import ErrorPage from './pages/error-page';
import LoginPage from './pages/login-page';

import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import React from 'react';

const router = createBrowserRouter([
  {
    path: "/",
    element: <HomePage/>,
    errorElement: <ErrorPage/>
  },
  {
    path: "/login",
    element: <LoginPage/>
  }
]);

function App() {
  return (
    <RouterProvider router={router} />
  )
}

export default App
