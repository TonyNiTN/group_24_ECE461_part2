import React from 'react';
import {useNavigate} from 'react-router-dom';

const InitialPage = () => {
  const navigate = useNavigate();

  const changeToLogin = () => {
    navigate('login');
  };

  const changeToSignUp = () => {
    navigate('signup');
  };

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <h1 className="text-center text-4xl font-extrabold text-purple-600 mb-8">Package Manager</h1>
        <div className="mt-6">
          <div className="flex flex-col items-center justify-center space-y-3">
            <button
              className="grow w-1/2 py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-purple-600 hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500"
              onClick={changeToLogin}
            >
              Login
            </button>
            <button
              className="w-1/2 py-2 px-4 border border-transparent text-sm font-medium rounded-md text-purple-600 bg-white hover:text-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500"
              onClick={changeToSignUp}
            >
              Sign Up
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default InitialPage;
