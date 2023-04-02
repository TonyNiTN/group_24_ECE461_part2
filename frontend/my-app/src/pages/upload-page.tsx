import React, { useState } from "react";


const UploadPage = () => {

    const [selectedFile, setSelectedFile] = useState('')
    const [isFileSelected, setIsFileSelected] = useState(false)


    const changeHandler = (event: any) => {
        event.preventDefault()
        setSelectedFile(event.target.files[0]);
		setIsFileSelected(true);
    }

    const submitHandler = () => {
        const formData = new FormData();

		formData.append('file', selectedFile);

		fetch(
			'/upload',
			{
				method: 'POST',
				body: formData,
			}
		)
			.then((response) => response.json())
			.then((result) => {
				console.log('Success:', result);
			})
			.catch((error) => {
				console.error('Error:', error);
			});
    }

    return (
       <div>
            <input type="file" onChange={changeHandler} className="text-md block border rounded-lg bg-purple-400"/>

                <button className="border rounded-2xl bg-purple-500" onClick={submitHandler}>
                    <div className="flex flex-row gap-x-4 items-center justify-center px-8 py-6">
                        <p className="text-xl text-gray-50">Upload Package</p>
                    </div>
                </button>
			
        </div>
    )
}

export default UploadPage;