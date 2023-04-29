const KEY = 'JWT';
const getJWT = () => {
  /**
   * get JWT from localStorage
   */
  return localStorage.getItem(KEY);
};

const setJWT = (token: string) => {
  /**
   * store JWT to localStorage
   */
  localStorage.setItem(KEY, token);
};

export {getJWT, setJWT};
