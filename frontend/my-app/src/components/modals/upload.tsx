import React, {useState} from 'react';
import { useNavigate } from 'react-router-dom';

import { getJWT } from '../../utils/token';
import { SERVICE } from '../../imports';

interface ModalProps {
  onClose: () => void;
  onChange: () => void;
}

const UploadModal: React.FC<ModalProps> = ({onClose, onChange}) => {
  const [selectedFile, setSelectedFile] = useState<File>(new File([''], ''));
  const [isFileSelected, setIsFileSelected] = useState(false);
  const [packageName, setPackageName] = useState('');
  const [packageURL, setPackageURL] = useState('');
  const [fileUpload, setFileUpload] = useState(false);
  const [statusMessage, setStatusMessage] = useState('');
  const navigate = useNavigate();

  const changeHandler = (event: any) => {
    event.preventDefault();
    if (event.target.files[0]) {
      setSelectedFile(event.target.files[0]);
      setIsFileSelected(true);
      return;
    }
    setIsFileSelected(false);
  };

  const submitHandler = () => {
    onChange();
    const token = getJWT();
    const formData = new FormData();
    
    formData.append('name', packageName);
    formData.append('url', packageURL);
    formData.append('file', selectedFile);

    fetch(`${SERVICE}/package`, {
      method: 'POST',
      headers: {
        Authorization: 'Bearer ' + token,
      },
      body: formData,
    })
      .then(response => {
        if (response.status === 401 || response.status === 403) {
          navigate('/');
        }
        return response.json();
      })
      .then(result => {
        console.log('Success:', result);
        setStatusMessage('File uploaded successfully!');
        setFileUpload(true);
        setTimeout(() => {
          setFileUpload(false);
        }, 1000);
      })
      .catch(error => {
        console.error('Error:', error);
        setStatusMessage('Error uploading file please try again.');
        setFileUpload(true);
        setTimeout(() => {
          setFileUpload(false);
        }, 1000);
      });
  };

  return (
    <div className="fixed z-50 inset-0 overflow-y-auto">
      <div className="flex items-center justify-center min-h-screen">
        <div className="relative bg-white w-full max-w-md mx-auto rounded-lg shadow-lg p-6">
          <button className="absolute top-0 right-0 mt-2 mr-2" onClick={onClose}>
            <svg className="h-6 w-6 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
          <h2 className="text-2xl font-bold mb-4">Upload Package</h2>
          {/* <form onSubmit={submitHandler} className=""> */}
            <div className="flex flex-col space-y-4">
              <label className="block w-full border border-purple-500 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm cursor-pointer py-3 px-4">
                <span className="block font-medium text-purple-500">Select a file</span>
                <input
                  type="file"
                  className="hidden"
                  onChange={changeHandler}
                  name="file"
                  accept="application/zip"
                  required
                />
              </label>

              <input
                type="text"
                name="name"
                className="block w-full py-2 px-3 border border-purple-500 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="Package Name"
                required
                value={packageName}
                onChange={e => setPackageName(e.target.value)}
              />

              <input
                type="text"
                name="url"
                className="block w-full py-2 px-3 border border-purple-500 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="Package URL"
                required
                value={packageURL}
                onChange={e => setPackageURL(e.target.value)}
              />

              <p className="text-sm text-gray-400">
                {isFileSelected ? `Selected file: ${selectedFile.name}` : 'No file selected'}
              </p>
              <p className="text-sm text-gray-400">{isFileSelected ? `File type: ${selectedFile.type}` : ''}</p>
              <p className="text-sm text-gray-400">{isFileSelected ? `File size: ${selectedFile.size} bytes` : ''}</p>
            </div>
            {fileUpload ? <></> : <p className="text-sm text-purple-700 py-3">{statusMessage}</p>}
            <button
              className="bg-purple-500 hover:bg-purple-700 text-white font-bold py-2 px-4 mt-3 rounded shadow-sm"
              // type='submit'
              onClick={submitHandler}
              // onSubmit={submitHandler}
            >
              Upload
            </button>
          {/* </form> */}
        </div>
      </div>
    </div>
  );
};

export default UploadModal;
