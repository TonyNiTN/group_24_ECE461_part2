import React, {useState} from 'react';
import {useNavigate} from 'react-router';

import {SERVICE} from '../imports';
import { setJWT } from '../utils/token';

const LoginPage = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loginSuc, setLoginSuc] = useState(false);
  const [loginFail, setLoginFail] = useState(false);
  const navigate = useNavigate();

  const changeToSignUp = () => {
    navigate('/signup');
  };

  const handleUsernameChange = (event: any) => {
    setUsername(event.target.value);
  };

  const handlePasswordChange = (event: any) => {
    setPassword(event.target.value);
  };

  const handleSubmit = (event: any) => {
    event.preventDefault();
    // Add logic to handle login submission here
    if (username == '') {
      alert('Username can not be empty');
      return;
    }

    if (password == '') {
      alert('Password can not be empty');
      return;
    }

    const data = {email: username, password: password, returnSecureToken: true};

    fetch(SERVICE + '/login', {
      method: 'POST',
      body: JSON.stringify(data),
    })
      .then(response => {
        if (response.status === 200) {
          return response.json();
        } else {
          throw new Error('Login unsuccessful. Try Again');
        }
      })
      .then(response => {
        setJWT(response.token);
        setLoginSuc(true);
         const timeOutId = setTimeout(() => {
           setLoginSuc(false);
           navigate('/home');
         }, 1000);
         return () => clearTimeout(timeOutId);
      })
      .catch(err => {
        console.log(err);
        setLoginFail(true);
        const timeOutId = setTimeout(() => {
          setLoginFail(false);
        }, 1000);
        return () => clearTimeout(timeOutId);
      });
  };

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <h1 className="text-center text-4xl font-extrabold text-purple-600 mb-8">Login</h1>
        <div className="mt-6">
          <form onSubmit={handleSubmit} className="bg-white p-8 rounded-lg shadow-md">
            <div className="mb-4">
              <label htmlFor="username" className="block text-gray-700 text-sm font-bold mb-2">
                Username
              </label>
              <input
                type="text"
                id="username"
                placeholder="Username"
                name="username"
                value={username}
                onChange={handleUsernameChange}
                className="border rounded-md py-2 px-3 w-full text-gray-700 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500"
                required
              />
            </div>
            <div className="mb-6">
              <label htmlFor="password" className="block text-gray-700 text-sm font-bold mb-2">
                Password
              </label>
              <input
                type="password"
                id="password"
                name="password"
                placeholder="Password"
                value={password}
                onChange={handlePasswordChange}
                className="border rounded-md py-2 px-3 w-full text-gray-700 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500"
                required
              />
            </div>
            {loginSuc ? <p className="text-green-500 pb-4">Login success!</p> : null}

            {loginFail ? <p className="text-red-500 pb-4">Login fail. Please try again.</p> : null}
            <button
              className="block w-full py-3 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-purple-600 hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500"
              type="submit"
              onClick={handleSubmit}
            >
              Login
            </button>
            <div className="flex justify-center mt-4">
              <p className="text-gray-600">Don't have an account?</p>
              <a className="text-purple-600 ml-2 font-bold hover:underline" onClick={changeToSignUp}>
                Sign up
              </a>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
