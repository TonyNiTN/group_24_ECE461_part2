import React from "react";
import "../index.css"

const LoginPage = () => {
    return(
        <div className="bg-amber-300 h-screen w-screen flex flex-col justify-center content-center">
            <p className="text-5xl font-bold text-purple-800">Package Manager</p>

            <div className="flex flex-col">
                <button className="rounded-3xl bg-slate-50 px-10 py-2 text-xl text-purple-800">Login</button>
                <button className="rounded-3xl bg-purple-500 px-10 py-2 text-xl text-slate-50">Sign Up</button>
            </div>
        </div>
    )
}

export default LoginPage