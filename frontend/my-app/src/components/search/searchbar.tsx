import React, {useState} from 'react';

interface SearchBarProps {
  onSubmit: (e: any) => void;
  searchInput: string;
  setSearchInput: React.Dispatch<React.SetStateAction<string>>;
}

const SearchBar: React.FC<SearchBarProps> = ({onSubmit, searchInput, setSearchInput}) => {
  
  return (
    <form className="grow p-4 shadow-md rounded-lg" onSubmit={onSubmit}>
      <p className="text-lg font-bold text-purple-600 pb-3">Search for a Package</p>
      <div className="flex flex-row items-center ">
        <div className="flex flex-row items-center border rounded-lg border-purple-600 mr-2 grow shadow-sm">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth="1.5"
            stroke="currentColor"
            className="w-6 h-6 text-gray-400 mx-2"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
            />
          </svg>
          <input
            type="text"
            name="search"
            id="search"
            className="block w-full focus:outline-none px-1 py-1.5"
            placeholder="Search Package"
            value={searchInput}
            onChange={e => setSearchInput(e.target.value)}
          />
        </div>
        <button
          className="bg-purple-500 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded shadow-sm"
          onSubmit={onSubmit}
        >
          Search
        </button>
      </div>
    </form>
  );
};

export default SearchBar;
