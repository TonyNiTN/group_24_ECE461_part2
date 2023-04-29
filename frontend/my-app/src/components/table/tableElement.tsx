import React, {useState} from 'react';
import {Package} from '../../imports';
import {SERVICE} from '../../imports';
import {getJWT} from '../../utils/token';
import { useNavigate } from 'react-router-dom';

interface TableElementProps {
  onChange: () => void;
  displayPackage: Package;
}

const TableElement: React.FC<TableElementProps> = ({onChange, displayPackage}) => {
  const [showScore, setShowScore] = useState(false);
  const navigate = useNavigate()
  const packageId = displayPackage.ID;

  const handleDownloadPackage = () => {
    const token = getJWT();

    fetch(`${SERVICE}/packages/${packageId}/download`, {
      method: 'GET',
      headers: {
        Authorization: 'Bearer ' + token,
        'Content-Type': 'application/octet-stream', // or the appropriate MIME type for your file
      },
    })
      .then(response => {
        if(response.status === 401 || response.status === 403) {
          navigate("/");
        } else if (response.status !== 200) {
          throw new Error('Can not download pacakage. Status code: ' + response.status);
        } 
        return response.blob()
      })
      .then(blob => {
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', `${displayPackage.Name}.zip`);
        document.body.appendChild(link);
        link.click();
        link.remove();
      })
      .catch(error => {
        console.error('Error:', error);
      });
  };

  const handleDeletePackage = () => {
    onChange();
    const token = getJWT();
    fetch(`${SERVICE}/packages/${packageId}/delete`, {
      method: 'DELETE',
      headers: {
        Authorization: 'Bearer ' + token,
      },
      body: null,
    })
      .then(response => {
        if (response.status === 401 || response.status === 403) {
          navigate('/');
        }
        console.log(response);
        return response.json();
      })
      .then(result => {
        console.log('Success:', result);
      })
      .catch(error => {
        console.error('Error:', error);
      });
  };

  const handleScorePackage = () => {
    onChange();
    const token = getJWT();
    fetch(`${SERVICE}/packages/${packageId}/score`, {
      method: 'GET',
      headers: {
        Authorization: 'Bearer ' + token,
      },
      body: null,
    })
      .then(response => {
        if (response.status === 401 || response.status === 403) {
          navigate('/');
        }
        return response.json()
      })
      .then(result => {
        console.log('Success:', result);
      })
      .catch(error => {
        console.error('Error:', error);
      });
  };

  const packageScoreBreakdown = (
    <div className="fixed z-50 inset-0 overflow-y-auto">
      <div className="flex items-center justify-center min-h-screen">
        <div className="relative bg-white w-full max-w-md mx-auto rounded-lg shadow-lg p-6">
          <button className="absolute top-0 right-0 mt-2 mr-2" onClick={() => setShowScore(false)}>
            <svg className="h-6 w-6 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
          <h2 className="text-2xl font-bold mb-4 text-gray-800">Package Score Breakdown</h2>
          <p className="text-gray-700">License Score: {displayPackage.LicenseScore}</p>
          <p className="text-gray-700">Correctness Score: {displayPackage.CorrectnessScore}</p>
          <p className="text-gray-700">Bus Factor Score: {displayPackage.BusFactorScore}</p>
          <p className="text-gray-700">Responsiveness Score: {displayPackage.ResponsivenessScore}</p>
          <p className="text-gray-700">Ramp Up Score: {displayPackage.RampUpScore}</p>
          <p className="text-gray-700">Version Score: {displayPackage.VersionScore}</p>
          <p className="text-gray-700">Review Score: {displayPackage.ReviewScore}</p>
        </div>
      </div>
    </div>
  );

  return (
    <tr className="bg-white hover:bg-gray-50">
      <td className="px-6 py-4 whitespace-nowrap text-gray-700">{displayPackage.Name}</td>
      <td className="px-6 py-4 whitespace-nowrap text-gray-700">{displayPackage.URL}</td>
      <td className="px-6 py-4 whitespace-nowrap text-gray-700">
        {displayPackage.NetScore}
        <button className="text-sm text-gray-400 pl-2" onClick={() => setShowScore(true)}>
          Show Breakdown
        </button>
        {showScore && packageScoreBreakdown}
      </td>
      <td className="px-6 py-4 whitespace-nowrap">
        <div className="flex items-center space-x-4">
          <button
            onClick={handleDownloadPackage}
            className="bg-purple-500 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded"
          >
            Download
          </button>
          <button
            onClick={handleScorePackage}
            className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
          >
            Score
          </button>
          <button
            onClick={handleDeletePackage}
            className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded"
          >
            Delete
          </button>
        </div>
      </td>
    </tr>
  );
};

export default TableElement;
