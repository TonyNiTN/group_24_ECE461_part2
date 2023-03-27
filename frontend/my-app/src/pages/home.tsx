import React from "react";

import Table from "../components/table/table";
import SearchBar from "../components/search/searchbar";
import UploadButton from "../components/buttons/uploadButton";
import OutlineButton from "../components/buttons/rateButton";

const Home = () => {
    return(
        <div className="">
            <SearchBar/>
            <UploadButton/>
            <OutlineButton/>
            <Table/>
        </div>
    )
}

export default Home