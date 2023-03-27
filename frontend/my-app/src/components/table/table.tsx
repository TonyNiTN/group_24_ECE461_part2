import TableElement from "./tableElement"

const entry = {
    packageName: "ReactJS",
    rating: "5 stars",
    size: "5GB",
    numOfDep: "20 dependencies, 40 dependents",
    score: "4/5"
}

const Table = () => {
    return (
        <table className="table-auto">
            <thead>
                <tr>
                    <th>Package Name</th>
                    <th>Package Size</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>ReactJS</td>
                    <td>5 GB</td>
                </tr>
            </tbody>
        </table>
    )
}

export default Table