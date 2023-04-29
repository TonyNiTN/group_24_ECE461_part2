import React from 'react';
import {Link} from 'react-router-dom';

const ErrorPage: React.FC = () => {
  return (
    <div className="mt-8">
      <h1 className="text-2xl font-bold">404</h1>
      <p>
        Page not found.{' '}
        <Link to="/" className="text-red-400 font-bold">
          Return Home
        </Link>
      </p>
    </div>
  );
};

export default ErrorPage;
