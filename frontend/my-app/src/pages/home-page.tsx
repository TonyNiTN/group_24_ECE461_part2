import React, {useState, useEffect} from 'react';

import UploadModal from '../components/modals/upload';
import SearchBar from '../components/search/searchbar';
import {Package} from '../imports';
import TableElement from '../components/table/tableElement';
import {SERVICE} from '../imports';

const HomePage = () => {
  const [isUploadModalOpen, setisUploadModalOpen] = useState(false);
  const [packages, setPackages] = useState<Package[]>([]);
  const [packageState, setPackageState] = useState(false);
  const [loadingPackage, setLoadingPackages] = useState(false);
  const [packageSearch, setPackageSearch] = useState('');

  const handleOpenModal = () => {
    setisUploadModalOpen(true);
  };

  const handleCloseModal = () => {
    setisUploadModalOpen(false);
  };

  const handleSubmitSearchBar = (e: any) => {
    e.preventDefault();
    console.log(packageSearch);
  };

  const taskChange = () => {
    setPackageState(!packageState);
  };

  useEffect(() => {
    setLoadingPackages(true);
    const timeOutId = setTimeout(() => {
      // List Task method
      fetch(`${SERVICE}/packages/search?name=${packageSearch}`, {
        method: 'GET',
      })
        .then(response => response.json())
        .then(result => {
          console.log('Success:', result);
          setLoadingPackages(false);
          setPackages(result);
        })
        .catch(error => {
          setLoadingPackages(false);
          console.error('Error:', error);
        });
    }, 500);
  }, []);
  //packageState, packageSearch
  const packageList = packages ? (
    packages.map(p => <TableElement displayPackage={p} onChange={taskChange} key={p.ID} />)
  ) : (
    <></>
  );

  return (
    <div>
      <p className="text-2xl font-bold text-purple-600 ml-4 mt-4">Package Repository</p>
      <div className="flex flex-row space-x-4 m-4">
        <SearchBar onSubmit={handleSubmitSearchBar} searchInput={packageSearch} setSearchInput={setPackageSearch} />
        <button
          className="bg-purple-500 hover:bg-purple-700 text-white font-bold py-2 px-4 mt-3 rounded-lg shadow-sm"
          onClick={handleOpenModal}
        >
          Upload Package
        </button>
        {isUploadModalOpen && <UploadModal onClose={handleCloseModal} />}
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
