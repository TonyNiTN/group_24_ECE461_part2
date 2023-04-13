import React, {useState} from "react";

import Table from "../components/table/table";
import SearchBar from "../components/search/searchbar";
import UploadButton from "../components/buttons/uploadButton";
import OutlineButton from "../components/buttons/rateButton";
import UploadModal from "../components/modals/upload";

const HomePage = () => {
    const [isUploadModalOpen, setisUploadModalOpen] = useState(false);

    const handleOpenModal = () => {
        setisUploadModalOpen(true);
    };

    const handleCloseModal = () => {
        setisUploadModalOpen(false);
    };

    return(
            <div>
                <button onClick={handleOpenModal}>Press me</button>
                {isUploadModalOpen && <UploadModal onClose={handleCloseModal}/>}
            </div>
          
    
    )
}

export default HomePage