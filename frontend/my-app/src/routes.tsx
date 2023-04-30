import React, { useEffect } from "react";
import { Routes, Route, useLocation, useNavigate } from "react-router-dom";

import HomePage from './pages/home-page';
import ErrorPage from './pages/error-page';
import InitialPage from './pages/initial-page';
import LoginPage from './pages/login-page';
import SignUpPage from './pages/signup-page';
import { getJWT } from "./utils/token";
import { SERVICE } from "./imports";

const NO_TOKEN_PATHS = ["/", "/login", "/signup"];

const Router: React.FC = () => {
  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    if (!NO_TOKEN_PATHS.includes(location.pathname)) {
      if (getJWT() === undefined || getJWT() === null) navigate("/");
    }
  }, [location.pathname]);

  return (
    <Routes location={location} key={location.pathname}>
      <Route path="/" element={<InitialPage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/signup" element={<SignUpPage />} />
      <Route path="/home" element={<HomePage />} />
      <Route path="/*" element={<ErrorPage />} />
    </Routes>
  );
};

export default Router;
