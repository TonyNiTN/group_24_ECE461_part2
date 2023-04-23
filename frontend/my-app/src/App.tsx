import { useState } from 'react'
import HomePage from './pages/home-page'
import ErrorPage from './pages/error-page';
import InitialPage from './pages/initial-page';
import LoginPage from './pages/login-page';
import SignUpPage from './pages/signup-page';

import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import React from 'react';
import { error } from 'console';
import { element } from 'prop-types';

const router = createBrowserRouter([
  {
    path: '/',
    element: <InitialPage />,
    errorElement: <ErrorPage />,
  },
  {
    path: 'login',
    element: <LoginPage />,
    errorElement: <ErrorPage />,
  },
  {
    path: 'signup',
    element: <SignUpPage />,
    errorElement: <ErrorPage />,
  },
  {
    path: 'home',
    element: <HomePage />,
    errorElement: <ErrorPage />,
  },
]);

function App() {
  return (
    <RouterProvider router={router} />
  )
}

export default App
