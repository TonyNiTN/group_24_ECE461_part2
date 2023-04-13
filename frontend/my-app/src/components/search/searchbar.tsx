import React, {useState} from "react"

interface SearchBarProps {
    onSubmit: () => void
}

const SearchBar: React.FC<SearchBarProps> = ({onSubmit}) => {
    return(
        <div >
            <input type="text" name="search bar" placeholder="Search Package Repository"
            className="bg-slate-800 rounded-md border-white border text-md text-white p-2 w-full"/>
        </div>
    )
}

export default SearchBar