import React, {useState, useEffect} from 'react';
import {useNavigate} from 'react-router-dom';

import UploadModal from '../components/modals/upload';
import SearchBar from '../components/search/searchbar';
import {Package} from '../imports';
import TableElement from '../components/table/tableElement';
import {SERVICE} from '../imports';
import {getJWT} from '../utils/token';
import DeleteModal from '../components/modals/delete';
import UpdateModal from '../components/modals/update';

const HomePage = () => {
  const [isUploadModalOpen, setisUploadModalOpen] = useState(false);
  const [isDeleteModalOpen, setisDeleteModalOpen] = useState(false);
 
  const [packages, setPackages] = useState<Package[]>([]);
  const [packageState, setPackageState] = useState(false);
  const [loadingPackage, setLoadingPackages] = useState(false);
  const [packageSearch, setPackageSearch] = useState('');
  const navigate = useNavigate();

  const handleOpenUploadModal = () => {
    setisUploadModalOpen(true);
  };

  const handleCloseUploadModal = () => {
    setisUploadModalOpen(false);
  };

  const handleOpenDeleteModal = () => {
    setisDeleteModalOpen(true);
  };

  const handleCloseDeleteModal = () => {
    setisDeleteModalOpen(false);
  };

  

  const handleSubmitSearchBar = (e: any) => {
    e.preventDefault();
    console.log(packageSearch);
  };

  const packageChange = () => {
    setPackageState(!packageState);
  };

  useEffect(() => {
    setLoadingPackages(true);
    const timeOutId = setTimeout(() => {
      const token = getJWT();
      fetch(`${SERVICE}/packages?name=${packageSearch}`, {
        method: 'POST',
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
        .then(result => {
          console.log('Success:', result);
          setLoadingPackages(false);
          if (result.message) {
            setPackages([]);
          } else {
            setPackages(result);
          }
        })
        .catch(error => {
          setLoadingPackages(false);
          console.error('Error:', error);
        });
    }, 500);
    return () => clearTimeout(timeOutId);
  }, [packageState, packageSearch]);
  const packageList = packages ? (
    packages.map(p => <TableElement displayPackage={p} onChange={packageChange} key={p.ID} />)
  ) : (
    <></>
  );

  return (
    <div>
      <div className="flex flex-ro ml-4 mt-4 items-center">
        <p className="text-2xl font-bold text-purple-600 pr-8">Package Repository</p>
        <button
          className="bg-gray-400 text-gray-50 font-bold p-2 rounded-lg shadow-sm hover:bg-gray-700"
          onClick={handleOpenDeleteModal}
        >
          Reset Repo
        </button>
      </div>
      <div className="flex flex-row space-x-4 m-4">
        <SearchBar onSubmit={handleSubmitSearchBar} searchInput={packageSearch} setSearchInput={setPackageSearch} />
        <button
          className="bg-purple-500 hover:bg-purple-700 text-white font-bold py-2 px-4 mt-3 rounded-lg shadow-sm"
          onClick={handleOpenUploadModal}
        >
          Upload Package
        </button>
        {isUploadModalOpen && <UploadModal onClose={handleCloseUploadModal} onChange={packageChange} />}
        {isDeleteModalOpen && <DeleteModal onClose={handleCloseDeleteModal} onChange={packageChange} />}
      </div>

      <div className="rounded-xl shadow-lg m-4">
        {loadingPackage ? (
          <p className="text-xl p-4 text-gray-500">Getting Packages...</p>
        ) : (
          <table className="min-w-full divide-y divide-gray-200">
            <thead>
              <tr>
                <th className="px-6 py-3 text-left text-sm font-medium text-purple-600 uppercase tracking-wider">
                  Package Name
                </th>
                <th className="px-6 py-3 text-left text-sm font-medium text-purple-600 uppercase tracking-wider">
                  Package URL
                </th>
                <th className="px-6 py-3 text-left text-sm font-medium text-purple-600 uppercase tracking-wider">
                  Package Netscore
                </th>
                <th className="px-6 py-3 text-left text-sm font-medium text-purple-600 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-300">{packageList}</tbody>
          </table>
        )}
      </div>
    </div>
  );
};

export default HomePage;
