import React, {useState} from 'react';
import {useNavigate} from 'react-router-dom';

import {getJWT} from '../../utils/token';
import { SERVICE } from '../../imports';

interface ModalProps {
  onClose: () => void;
  onChange: () => void;
}

const DeleteModal: React.FC<ModalProps> = ({onClose, onChange}) => {
  const token = getJWT();
  const navigate = useNavigate();
  const [deleteConfirm, setDeleteConfrim] = useState('');

  const submitHandler = () => {
    onChange();
    fetch(`${SERVICE}/packages/delete`, {
      method: 'DELETE',
      headers: {
        Authorization: 'Bearer ' + token,
      },
    })
      .then(response => {
        if (response.status === 401 || response.status === 403) {
          navigate('/');
        }
        return response.json();
      })
      .then(result => {})
      .catch(error => {});
  };

  const validate = () => {
    return deleteConfirm === "Delete"
  }

  return (
    <div className="fixed z-50 inset-0 overflow-y-auto">
      <div className="flex items-center justify-center min-h-screen">
        <div className="relative bg-white w-full max-w-md mx-auto rounded-lg shadow-lg p-6">
          <button className="absolute top-0 right-0 mt-2 mr-2" onClick={onClose}>
            <svg className="h-6 w-6 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
          <h2 className="text-2xl font-bold mb-4">Confirm Delete</h2>
          <p className="text-sm text-gray-600 pb-2">
            You will delete all packages in the repo. Type "Delete" to Delete
          </p>
          <form onSubmit={submitHandler} className="">
            <div className="flex flex-col space-y-4">
              <input
                type="text"
                name="name"
                className="block w-full py-2 px-3 border border-purple-500 bg-white rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm"
                placeholder="Delete"
                required
                value={deleteConfirm}
                onChange={e => setDeleteConfrim(e.target.value)}
              />
            </div>
            <button
              className="bg-red-500 disabled:opacity-50 text-white font-bold py-2 px-4 mt-3 rounded shadow-sm"
              type="submit"
              onSubmit={submitHandler}
              disabled={!validate()}
            >
              Delete
            </button>
          </form>
        </div>
      </div>
    </div>
  );
};

export default DeleteModal;
